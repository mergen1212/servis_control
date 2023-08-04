package db

import (
	"fmt"
	"os"
	"testing"
	"time"
)

func TestPrepareDBMustCreateTables(t *testing.T) {
	db, _ := GetDB()
	db.PrepareDB()

	tables := []string{"user", "project"}
	for i := 0; i < len(tables); i++ {
		table := tables[i]
		row := db.oneRow("select name from sqlite_master where name = ?", table)
		var dbTableName string
		err := row.Scan(&dbTableName)
		if err != nil {
			t.Error("Error happened", err)
		}

		if dbTableName != table {
			t.Error(fmt.Sprintf("Table %s not found", table))
		}
	}
	os.Remove(databasePath)
}

func TestGetUser(t *testing.T) {
	db, _ := GetDB()
	db.PrepareDB()

	user := User{1, 1}
	db.exec("insert into user (id, telegram_id) values (?, ?)", user.ID, user.TelegramID)

	testUser, err := db.GetUser(1)
	if err != nil {
		t.Error("GetUser returned error", err)
	}

	if testUser.ID != user.ID {
		t.Error("test_user id is not equal", testUser.ID, user.ID)
	}
	os.Remove(databasePath)
}

func TestGetProject(t *testing.T) {
	db, _ := GetDB()
	db.PrepareDB()

	user := User{1, 1}
	db.exec("insert into user (id, telegram_id) values (?, ?)", user.ID, user.TelegramID)

	project := Project{1, "asd", user.ID, time.Now()}
	db.exec(
		"insert into project (id, hash, user_id) values (?, ?, ?)",
		project.ID, project.Hash, project.UserID)

	dbProject, err := db.GetProjectByHash(project.Hash)
	if err != nil {
		t.Error("GetProjectByHash returned error", err)
	}

	if dbProject.ID != project.ID {
		t.Error("dbProject id is not equal", dbProject.ID, project.ID)
	}
	os.Remove(databasePath)
}

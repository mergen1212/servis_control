package main

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"watchman/pkg/db"
)

func handleRequest(w *httptest.ResponseRecorder, r *http.Request) {
	router, _ := getRouter()
	router.ServeHTTP(w, r)
}

func TestUserNotFound(t *testing.T) {
	request, _ := http.NewRequest("GET", "/1/1", nil)
	recorder := httptest.NewRecorder()
	handleRequest(recorder, request)
	if recorder.Code != 404 {
		t.Error("User does not exist, must be 404")
	}
}

func TestUserBadUserId(t *testing.T) {
	request, _ := http.NewRequest("GET", "/foo/1", nil)
	recorder := httptest.NewRecorder()
	handleRequest(recorder, request)
	if recorder.Code != 400 {
		t.Error("User id is string, must be 400")
	}
}

func TestProjectNotFound(t *testing.T) {
	database, _ := db.GetDB()
	database.Exec("insert into user (id, telegram_id) values (1, 1)")
	request, _ := http.NewRequest("GET", "/1/bad_project_hash", nil)
	recorder := httptest.NewRecorder()
	handleRequest(recorder, request)
	if recorder.Code != 404 {
		t.Error("User does not exist, must be 404")
	}
}

func TestUserTriedToUseNotHisProject(t *testing.T) {
	user1 := 1
	user2 := 2
	database, _ := db.GetDB()
	database.Exec(fmt.Sprintf("insert into user (id, telegram_id) values (%d, 1)", user1))
	database.Exec(fmt.Sprintf("insert into user (id, telegram_id) values (%d, 1)", user2))

	projectHash := "foo"
	database.Exec(fmt.Sprintf("insert into project (id, hash, user_id) values (2, '%s', %d)", projectHash, user1))

	request, _ := http.NewRequest("GET", fmt.Sprintf("/%d/%s", user2, projectHash), nil)
	recorder := httptest.NewRecorder()
	handleRequest(recorder, request)
	if recorder.Code != 404 {
		t.Error("User 2 does not have that project, must be 404")
	}
}


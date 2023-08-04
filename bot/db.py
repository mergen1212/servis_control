import uuid
import aiosqlite

from contextlib import asynccontextmanager
from tp import Project, TelegramID, User, UserID

DB_NAME = "storage.db"


def make_project_hash() -> str:
    return uuid.uuid4().hex


@asynccontextmanager
async def open_db():
    async with aiosqlite.connect(f"../{DB_NAME}") as db:
        yield db


def prepare_user(user_id: int, tg_id: int) -> User:
    return {
        "id": UserID(user_id),
        "telegram_id": TelegramID(tg_id),
    }


async def get_or_create_user(telegram_id: TelegramID) -> User:
    async with open_db() as db:
        cursor = await db.execute(
            "select id, telegram_id from user where telegram_id = ? limit 1",
            [telegram_id],
        )
        user = await cursor.fetchone()
        if user:
            return prepare_user(*user)

        cursor = await db.execute(
            "insert into user (telegram_id) values (?)", [telegram_id]
        )
        user_id = cursor.lastrowid
        return {
            "id": UserID(user_id),
            "telegram_id": telegram_id,
        }


async def create_project(user_id: UserID) -> Project:
    async with open_db() as db:
        hash = make_project_hash()
        cursor = await db.execute(
            "insert into project (hash, user_id) values (?, ?)",
            [user_id, hash],
        )
        project_id = cursor.lastrowid
        return {
            "id": project_id,
            "hash": hash,
            "user_id": user_id,
            "updated": None,
        }

from typing import TypedDict, NewType


TelegramID = NewType("TelegramID", int)
UserID = NewType("UserID", int)


class User(TypedDict):
    id: UserID
    telegram_id: TelegramID


class Project(TypedDict):
    id: int
    hash: str
    user_id: UserID
    updated: int | None

from enum import Enum
from typing import Any


class Metadata:
    trace_id: str
    user_id: str

    def __init__(self, trace_id: str, user_id: str) -> None:
        self.trace_id = trace_id
        self.user_id = user_id


class EventTypes(Enum):
    ADD_USER = "AddUser"
    MODIFY_USER = "ModifyUser"
    USER_ADDED = "UserAdded"
    USER_MODIFIED = "UserModified"


class Message:
    data: Any
    id: str
    metadata: Metadata
    type: EventTypes

    def __init__(self, data: Any, id: str, metadata: Metadata, type: EventTypes) -> None:
        self.data = data
        self.id = id
        self.metadata = metadata
        self.type = type


class Categories(Enum):
    ESTIMATION = "estimation"
    IDENTITY = "identity"
    PROJECTS = "projects"


class Commands(Enum):
    ADD_USER = "AddUser"
    MODIFY_USER = "ModifyUser"


class AddUserCommandData:
    auth0_id: str
    email: str
    email_verified: bool
    first_name: str
    id: str
    last_name: str
    locale: str
    picture: str

    def __init__(self, auth0_id: str, email: str, email_verified: bool, first_name: str, id: str, last_name: str, locale: str, picture: str) -> None:
        self.auth0_id = auth0_id
        self.email = email
        self.email_verified = email_verified
        self.first_name = first_name
        self.id = id
        self.last_name = last_name
        self.locale = locale
        self.picture = picture


class AddUserCommandType(Enum):
    ADD_USER = "AddUser"


class AddUserCommand:
    data: AddUserCommandData
    id: str
    metadata: Metadata
    type: AddUserCommandType

    def __init__(self, data: AddUserCommandData, id: str, metadata: Metadata, type: AddUserCommandType) -> None:
        self.data = data
        self.id = id
        self.metadata = metadata
        self.type = type


class ModifyUserCommandData:
    first_name: str
    last_name: str
    locale: str
    picture: str

    def __init__(self, first_name: str, last_name: str, locale: str, picture: str) -> None:
        self.first_name = first_name
        self.last_name = last_name
        self.locale = locale
        self.picture = picture


class ModifyUserCommandType(Enum):
    MODIFY_USER = "ModifyUser"


class ModifyUserCommand:
    data: ModifyUserCommandData
    id: str
    metadata: Metadata
    type: ModifyUserCommandType

    def __init__(self, data: ModifyUserCommandData, id: str, metadata: Metadata, type: ModifyUserCommandType) -> None:
        self.data = data
        self.id = id
        self.metadata = metadata
        self.type = type


class Events(Enum):
    USER_ADDED = "UserAdded"
    USER_MODIFIED = "UserModified"


class UserAddedEventData:
    auth0_id: str
    email: str
    email_verified: bool
    first_name: str
    id: str
    last_name: str
    locale: str
    picture: str

    def __init__(self, auth0_id: str, email: str, email_verified: bool, first_name: str, id: str, last_name: str, locale: str, picture: str) -> None:
        self.auth0_id = auth0_id
        self.email = email
        self.email_verified = email_verified
        self.first_name = first_name
        self.id = id
        self.last_name = last_name
        self.locale = locale
        self.picture = picture


class UserAddedEventType(Enum):
    USER_ADDED = "UserAdded"


class UserAddedEvent:
    data: UserAddedEventData
    id: str
    metadata: Metadata
    type: UserAddedEventType

    def __init__(self, data: UserAddedEventData, id: str, metadata: Metadata, type: UserAddedEventType) -> None:
        self.data = data
        self.id = id
        self.metadata = metadata
        self.type = type


class UserModifiedEventData:
    first_name: str
    last_name: str
    locale: str
    picture: str

    def __init__(self, first_name: str, last_name: str, locale: str, picture: str) -> None:
        self.first_name = first_name
        self.last_name = last_name
        self.locale = locale
        self.picture = picture


class UserModifiedEventType(Enum):
    USER_MODIFIED = "UserModified"


class UserModifiedEvent:
    data: UserModifiedEventData
    id: str
    metadata: Metadata
    type: UserModifiedEventType

    def __init__(self, data: UserModifiedEventData, id: str, metadata: Metadata, type: UserModifiedEventType) -> None:
        self.data = data
        self.id = id
        self.metadata = metadata
        self.type = type

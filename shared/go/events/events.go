// This file was generated from JSON Schema using quicktype, do not modify it directly.
// To parse and unparse this JSON data, add this code to your project and do:
//
//    eventTypes, err := UnmarshalEventTypes(bytes)
//    bytes, err = eventTypes.Marshal()
//
//    message, err := UnmarshalMessage(bytes)
//    bytes, err = message.Marshal()
//
//    metadata, err := UnmarshalMetadata(bytes)
//    bytes, err = metadata.Marshal()
//
//    categories, err := UnmarshalCategories(bytes)
//    bytes, err = categories.Marshal()
//
//    commands, err := UnmarshalCommands(bytes)
//    bytes, err = commands.Marshal()
//
//    addUserCommand, err := UnmarshalAddUserCommand(bytes)
//    bytes, err = addUserCommand.Marshal()
//
//    modifyUserCommand, err := UnmarshalModifyUserCommand(bytes)
//    bytes, err = modifyUserCommand.Marshal()
//
//    events, err := UnmarshalEvents(bytes)
//    bytes, err = events.Marshal()
//
//    userAddedEvent, err := UnmarshalUserAddedEvent(bytes)
//    bytes, err = userAddedEvent.Marshal()
//
//    userModifiedEvent, err := UnmarshalUserModifiedEvent(bytes)
//    bytes, err = userModifiedEvent.Marshal()

package events

import "encoding/json"

func UnmarshalEventTypes(data []byte) (EventTypes, error) {
	var r EventTypes
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *EventTypes) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func UnmarshalMessage(data []byte) (Message, error) {
	var r Message
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *Message) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func UnmarshalMetadata(data []byte) (Metadata, error) {
	var r Metadata
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *Metadata) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func UnmarshalCategories(data []byte) (Categories, error) {
	var r Categories
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *Categories) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func UnmarshalCommands(data []byte) (Commands, error) {
	var r Commands
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *Commands) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func UnmarshalAddUserCommand(data []byte) (AddUserCommand, error) {
	var r AddUserCommand
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *AddUserCommand) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func UnmarshalModifyUserCommand(data []byte) (ModifyUserCommand, error) {
	var r ModifyUserCommand
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *ModifyUserCommand) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func UnmarshalEvents(data []byte) (Events, error) {
	var r Events
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *Events) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func UnmarshalUserAddedEvent(data []byte) (UserAddedEvent, error) {
	var r UserAddedEvent
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *UserAddedEvent) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func UnmarshalUserModifiedEvent(data []byte) (UserModifiedEvent, error) {
	var r UserModifiedEvent
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *UserModifiedEvent) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type Message struct {
	Data     interface{} `json:"data"`    
	ID       string      `json:"id"`      
	Metadata Metadata    `json:"metadata"`
	Type     EventTypes  `json:"type"`    
}

type Metadata struct {
	TraceID string `json:"traceId"`
	UserID  string `json:"userId"` 
}

type AddUserCommand struct {
	Data     AddUserCommandData `json:"data"`    
	ID       string             `json:"id"`      
	Metadata Metadata           `json:"metadata"`
	Type     AddUserCommandType `json:"type"`    
}

type AddUserCommandData struct {
	Auth0ID       string `json:"auth0Id"`      
	Email         string `json:"email"`        
	EmailVerified bool   `json:"emailVerified"`
	FirstName     string `json:"firstName"`    
	ID            string `json:"id"`           
	LastName      string `json:"lastName"`     
	Locale        string `json:"locale"`       
	Picture       string `json:"picture"`      
}

type ModifyUserCommand struct {
	Data     ModifyUserCommandData `json:"data"`    
	ID       string                `json:"id"`      
	Metadata Metadata              `json:"metadata"`
	Type     ModifyUserCommandType `json:"type"`    
}

type ModifyUserCommandData struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"` 
	Locale    string `json:"locale"`   
	Picture   string `json:"picture"`  
}

type UserAddedEvent struct {
	Data     UserAddedEventData `json:"data"`    
	ID       string             `json:"id"`      
	Metadata Metadata           `json:"metadata"`
	Type     UserAddedEventType `json:"type"`    
}

type UserAddedEventData struct {
	Auth0ID       string `json:"auth0Id"`      
	Email         string `json:"email"`        
	EmailVerified bool   `json:"emailVerified"`
	FirstName     string `json:"firstName"`    
	ID            string `json:"id"`           
	LastName      string `json:"lastName"`     
	Locale        string `json:"locale"`       
	Picture       string `json:"picture"`      
}

type UserModifiedEvent struct {
	Data     UserModifiedEventData `json:"data"`    
	ID       string                `json:"id"`      
	Metadata Metadata              `json:"metadata"`
	Type     UserModifiedEventType `json:"type"`    
}

type UserModifiedEventData struct {
	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"` 
	Locale    string `json:"locale"`   
	Picture   string `json:"picture"`  
}

type EventTypes string
const (
	EventTypesAddUser EventTypes = "AddUser"
	EventTypesModifyUser EventTypes = "ModifyUser"
	EventTypesUserAdded EventTypes = "UserAdded"
	EventTypesUserModified EventTypes = "UserModified"
)

type Categories string
const (
	Estimation Categories = "estimation"
	Identity Categories = "identity"
	Projects Categories = "projects"
)

type Commands string
const (
	CommandsAddUser Commands = "AddUser"
	CommandsModifyUser Commands = "ModifyUser"
)

type AddUserCommandType string
const (
	TypeAddUser AddUserCommandType = "AddUser"
)

type ModifyUserCommandType string
const (
	TypeModifyUser ModifyUserCommandType = "ModifyUser"
)

type Events string
const (
	EventsUserAdded Events = "UserAdded"
	EventsUserModified Events = "UserModified"
)

type UserAddedEventType string
const (
	TypeUserAdded UserAddedEventType = "UserAdded"
)

type UserModifiedEventType string
const (
	TypeUserModified UserModifiedEventType = "UserModified"
)

package msg

import "encoding/json"

func UnmarshalEventTypes(data []byte) (EventTypes, error) {
	var r EventTypes
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *EventTypes) Marshal() ([]byte, error) {
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

func UnmarshalMembershipCreatedEvent(data []byte) (MembershipCreatedEvent, error) {
	var r MembershipCreatedEvent
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *MembershipCreatedEvent) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func UnmarshalMembershipCreatedForProjectEvent(data []byte) (MembershipCreatedForProjectEvent, error) {
	var r MembershipCreatedForProjectEvent
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *MembershipCreatedForProjectEvent) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func UnmarshalMembershipUpdatedEvent(data []byte) (MembershipUpdatedEvent, error) {
	var r MembershipUpdatedEvent
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *MembershipUpdatedEvent) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func UnmarshalMembershipDeletedEvent(data []byte) (MembershipDeletedEvent, error) {
	var r MembershipDeletedEvent
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *MembershipDeletedEvent) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func UnmarshalProjectCreatedEvent(data []byte) (ProjectCreatedEvent, error) {
	var r ProjectCreatedEvent
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *ProjectCreatedEvent) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func UnmarshalProjectUpdatedEvent(data []byte) (ProjectUpdatedEvent, error) {
	var r ProjectUpdatedEvent
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *ProjectUpdatedEvent) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

func UnmarshalProjectDeletedEvent(data []byte) (ProjectDeletedEvent, error) {
	var r ProjectDeletedEvent
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *ProjectDeletedEvent) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type MembershipCreatedEvent struct {
	Data     MembershipCreatedEventData `json:"data"`
	Metadata Metadata                   `json:"metadata"`
	Type     MembershipCreatedEventType `json:"type"`
}

type MembershipCreatedEventData struct {
	TenantID     string `json:"tenantId"`
	CreatedAt    string `json:"createdAt"`
	MembershipID string `json:"membershipId"`
	Role         string `json:"role"`
	TeamID       string `json:"teamId"`
	UpdatedAt    string `json:"updatedAt"`
	UserID       string `json:"userId"`
}

type MembershipCreatedForProjectEvent struct {
	Data     MembershipCreatedForProjectEventData `json:"data"`
	Metadata Metadata                             `json:"metadata"`
	Type     MembershipCreatedForProjectEventType `json:"type"`
}

type MembershipCreatedForProjectEventData struct {
	TenantID     string `json:"tenantId"`
	CreatedAt    string `json:"createdAt"`
	MembershipID string `json:"membershipId"`
	ProjectID    string `json:"projectId"`
	Role         string `json:"role"`
	TeamID       string `json:"teamId"`
	UpdatedAt    string `json:"updatedAt"`
	UserID       string `json:"userId"`
}

type MembershipUpdatedEvent struct {
	Data     MembershipUpdatedEventData `json:"data"`
	Metadata Metadata                   `json:"metadata"`
	Type     MembershipUpdatedEventType `json:"type"`
}

type MembershipUpdatedEventData struct {
	MembershipID string `json:"membershipId"`
	Role         string `json:"role"`
	UpdatedAt    string `json:"updatedAt"`
}

type MembershipDeletedEvent struct {
	Data     MembershipDeletedEventData `json:"data"`
	Metadata Metadata                   `json:"metadata"`
	Type     MembershipDeletedEventType `json:"type"`
}

type MembershipDeletedEventData struct {
	MembershipID string `json:"membershipId"`
}

type ProjectCreatedEvent struct {
	Data     ProjectCreatedEventData `json:"data"`
	Metadata Metadata                `json:"metadata"`
	Type     ProjectCreatedEventType `json:"type"`
}

type ProjectCreatedEventData struct {
	TenantID    string   `json:"tenantId"`
	Active      bool     `json:"active"`
	ColumnOrder []string `json:"columnOrder"`
	CreatedAt   string   `json:"createdAt"`
	Description string   `json:"description"`
	Name        string   `json:"name"`
	Prefix      string   `json:"prefix"`
	ProjectID   string   `json:"projectId"`
	Public      bool     `json:"public"`
	TeamID      string   `json:"teamId"`
	UpdatedAt   string   `json:"updatedAt"`
	UserID      string   `json:"userId"`
}

type ProjectUpdatedEvent struct {
	Data     ProjectUpdatedEventData `json:"data"`
	Metadata Metadata                `json:"metadata"`
	Type     ProjectUpdatedEventType `json:"type"`
}

type ProjectUpdatedEventData struct {
	Active      *bool    `json:"active,omitempty"`
	ColumnOrder []string `json:"columnOrder,omitempty"`
	Description *string  `json:"description,omitempty"`
	Name        *string  `json:"name,omitempty"`
	ProjectID   string   `json:"projectId"`
	Public      *bool    `json:"public,omitempty"`
	TeamID      *string  `json:"teamId,omitempty"`
	UpdatedAt   string   `json:"updatedAt"`
}

type ProjectDeletedEvent struct {
	Data     ProjectDeletedEventData `json:"data"`
	Metadata Metadata                `json:"metadata"`
	Type     ProjectDeletedEventType `json:"type"`
}

type ProjectDeletedEventData struct {
	ProjectID string `json:"projectId"`
}

type EventTypes string

const (
	EventTypesMembershipCreated           EventTypes = "MembershipCreated"
	EventTypesMembershipCreatedForProject EventTypes = "MembershipCreatedForProject"
	EventTypesMembershipDeleted           EventTypes = "MembershipDeleted"
	EventTypesMembershipUpdated           EventTypes = "MembershipUpdated"
	EventTypesProjectCreated              EventTypes = "ProjectCreated"
	EventTypesProjectDeleted              EventTypes = "ProjectDeleted"
	EventTypesProjectUpdated              EventTypes = "ProjectUpdated"
)

type Categories string

const (
	Project Categories = "Project"
	Users   Categories = "Users"
)

type Events string

const (
	EventsMembershipCreated           Events = "MembershipCreated"
	EventsMembershipCreatedForProject Events = "MembershipCreatedForProject"
	EventsMembershipDeleted           Events = "MembershipDeleted"
	EventsMembershipUpdated           Events = "MembershipUpdated"
	EventsProjectCreated              Events = "ProjectCreated"
	EventsProjectDeleted              Events = "ProjectDeleted"
	EventsProjectUpdated              Events = "ProjectUpdated"
)

type MembershipCreatedEventType string

const (
	TypeMembershipCreated MembershipCreatedEventType = "MembershipCreated"
)

type MembershipCreatedForProjectEventType string

const (
	TypeMembershipCreatedForProject MembershipCreatedForProjectEventType = "MembershipCreatedForProject"
)

type MembershipUpdatedEventType string

const (
	TypeMembershipUpdated MembershipUpdatedEventType = "MembershipUpdated"
)

type MembershipDeletedEventType string

const (
	TypeMembershipDeleted MembershipDeletedEventType = "MembershipDeleted"
)

type ProjectCreatedEventType string

const (
	TypeProjectCreated ProjectCreatedEventType = "ProjectCreated"
)

type ProjectUpdatedEventType string

const (
	TypeProjectUpdated ProjectUpdatedEventType = "ProjectUpdated"
)

type ProjectDeletedEventType string

const (
	TypeProjectDeleted ProjectDeletedEventType = "ProjectDeleted"
)

package project

type ProjectMember struct {
	ID string 
	UserID string 
	ProjectID string
	Maintainer bool
	InviteSent bool
	InviteAccepted bool
}

package membership

// Enums
type Role int

const (
	Admin = iota
	Editor
	Commenter
	Viewer
)

func (r Role) String() string {
	return [...]string{"admin", "editor", "commenter", "viewer"}[r]
}

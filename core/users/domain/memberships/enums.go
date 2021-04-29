package memberships

// Enums
type Role int

const (
	Administrator = iota
	Editor
	Commenter
	Viewer
)

func (r Role) String() string {
	return [...]string{"administrator", "editor", "commenter", "viewer"}[r]
}

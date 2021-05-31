package memberships

// Role type for enumerated values
type Role int

// Roles
const (
	Administrator = iota
	Editor
	Commenter
	Viewer
)

func (r Role) String() string {
	return [...]string{"administrator", "editor", "commenter", "viewer"}[r]
}

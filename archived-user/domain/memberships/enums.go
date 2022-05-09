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

// String retrieves the corresponding string value for a role
func (r Role) String() string {
	return [...]string{"administrator", "editor", "commenter", "viewer"}[r]
}

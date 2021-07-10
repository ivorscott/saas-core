package testhelpers

import (
	"github.com/devpies/devpie-client-core/users/platform/database"
	mockDB "github.com/devpies/devpie-client-core/users/platform/database/mocks"
	"net/url"
)

// Repo initializes a mock repository
func Repo() *database.Repository {
	return &database.Repository{
		SqlxStorer:      &mockDB.SqlxStorer{},
		SquirrelBuilder: &mockDB.Squirreler{},
		URL:             url.URL{},
	}
}

// StringPointer converts a string to a string pointer
func StringPointer(s string) *string {
	return &s
}

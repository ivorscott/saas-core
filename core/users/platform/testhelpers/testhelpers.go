package testhelpers

import (
	"github.com/devpies/devpie-client-core/users/platform/database"
	mockDB "github.com/devpies/devpie-client-core/users/platform/database/mocks"
	"net/url"
)

func Repo() *database.Repository {
	return &database.Repository{
		SqlxStorer: &mockDB.SqlxStorer{},
		Squirreler: &mockDB.Squirreler{},
		URL:        url.URL{},
	}
}

func StringPointer(s string) *string {
	return &s
}

package testutils

import (
	"database/sql"
	"github.com/go-testfixtures/testfixtures/v3"
	_ "github.com/lib/pq" // required by testfixtures
)

var fixturesDir = "fixtures"

func loadFixtures(db *sql.DB) error {
	fixtures, err := testfixtures.New(
		testfixtures.Database(db),
		testfixtures.Dialect(dbDriver),
		testfixtures.Directory(resFile(fixturesDir)),
	)
	if err != nil {
		return err
	}

	return fixtures.Load()
}

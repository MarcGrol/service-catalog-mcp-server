package sloconstants

import "os"

const (
	// SLODatabaseFilenameKey offers a typestrong key for the catalog database filename
	SLODatabaseFilenameKey = "slo-databasefile"
)

// SLODatabaseFilename is the default filename for the SQLite database that describes the SLO.
func SLODatabaseFilename() string {
	homedir := os.Getenv("HOME")
	if homedir == "" {
		homedir = "/Users/marcgrol"
	}

	return homedir + "/src/service-catalog-mcp-server/data/slos.sqlite"
}

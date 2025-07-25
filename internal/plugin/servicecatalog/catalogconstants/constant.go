package catalogconstants

import "os"

const (
	// CatalogDatabaseFilenameKey offers a typestrong key for the catalog database filename
	CatalogDatabaseFilenameKey = "catalog-databasefile"
)

// CatalogDatabaseFilename returns the filename of the catalog database
// TODO: make this configurable
func CatalogDatabaseFilename() string {
	homedir := os.Getenv("HOME")
	if homedir == "" {
		homedir = "/Users/marcgrol"
	}

	return homedir + "/src/service-catalog-mcp-server/internal/plugin/servicecatalog/service-catalog.sqlite"
}

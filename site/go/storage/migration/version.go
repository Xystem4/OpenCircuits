package version

import (
	"database/sql"
	"errors"
	"sort"

	_ "github.com/mattn/go-sqlite3" // We don't really need code in this package, just symbols.
)

type versionStruct struct {
	id      int
	dateTag string
}

type sqliteVersionStorageInterface struct {
	db *sql.DB

	// loadEntryStmt    *sql.Stmt
	// storeEntryStmt   *sql.Stmt
	// createEntryStmt  *sql.Stmt
	// enumCircuitsStmt *sql.Stmt
	// deleteEntryStmt  *sql.Stmt
}

// TableInit initiaize version table. Called by main.
func TableInit(db *sql.DB) error {
	vs := parser(db)
	if !integrityCheck(vs) {
		return errors.New("Table corrupted")
	}

	return nil
}

// Pasering data from version table
func parser(db *sql.DB) []versionStruct {
	rows, err := db.Query("SELECT id, date-tag FROM version")
	defer rows.Close()
	if err != nil {
		return nil
	}

	var v versionStruct
	var vs []versionStruct
	for rows.Next() {
		err := rows.Scan(&v.id, &v.dateTag)
		if err != nil {
			return nil
		}
		vs = append(vs, v)
	}
	return vs
}

// Check table's integreity.
func integrityCheck(vs []versionStruct) bool {

	// sort version slice
	sort.Slice(vs, func(i, j int) bool {
		return vs[i].id < vs[j].id
	})

	// TODO when new versions are available, verify their record
	// format, date, etc.
	// loop over to check integrity...
	// for i, v := range vs {
	//		if ... {
	//			return false
	// 		}
	//
	// }

	return true
}

// Pre versioning handle. If format matches, return true.
func isPreVersioning(db *sql.DB) bool {
	rows, _ := db.Query("PRAGMA table_info(circuits)") // meta info
	defer rows.Close()

	var fieldNames []string

	for rows.Next() {

		// meta info table scheme
		// only using names

		var cid int64
		var name string
		var dtype string
		var notnull bool
		var dfltValue string
		var pk bool

		err := rows.Scan(&cid, &name, &dtype, &notnull, &dfltValue, &pk)
		if err != nil {
			return errors.New("Table corrupted")
		}
		fieldNames = append(fieldNames, name)
	}

	// This structure is the default for pre versioning versions.
	if fieldNames[0] == "id" &&
		fieldNames[0] == "name" &&
		fieldNames[2] == "designer" &&
		fieldNames[3] == "ownerId" &&
		fieldNames[4] == "version" &&
		fieldNames[5] == "thumbnail" {
		return true
	} else {
		return false
	}
}
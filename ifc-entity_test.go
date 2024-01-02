package datarepo_test

import (
	"database/sql"
	"fmt"
	DRS "github.com/fbaube/datarepo/sqlite"
	"io/ioutil"
	"os"
)

// ExampleRepoEntity tests (for SQLite) the interface:
//   - Handle() *sqlx.DB  // (noun) the handle to the DB
//   - Type() D.DB_type   // DB_SQLite ("sqlite", equiv.to "sqlite3")
//   - Path() string      // file/URL (or dir/URL, if uses multiple files)
//   - IsURL() bool       // false for a local SQLite file
//   - IsSingleFile() bool // true for SQLite
//
// .
func ExampleRepoEntity() {
	var F *os.File
	// S implements repo.RepoLifecycle
	var S *DRS.SqliteRepo
	var e error
	// func TempFile(dir, pattern string) (f *os.File, err error)
	F, e = ioutil.TempFile("", "*.db")
	if e != nil {
		panic(e)
	}
	// We can't print out the path, because it is
	// random but the output would have to be matched.
	// fmt.Println("Path:", F.Name())
	F.Close()

	// Now the test begins
	S, e = DRS.OpenRepoAtPath(F.Name())
	if e != nil {
		panic(e)
	}
	//   - Handle() *sqlx.DB  // (noun) the handle to the DB
	//   - Type() D.DB_type   // DB_SQLite ("sqlite", equiv.to "sqlite3")
	//   - Path() string      // file/URL (or dir/URL, if uses multiple files)
	//   - IsURL() bool       // false for a local SQLite file
	//   - IsSingleFile() bool // true for SQLite
	var H *sql.DB
	var T string
	var isU, isS bool
	H = S.Handle()
	fmt.Printf("Typeof Handle: %T\n", H)
	T = string(S.Type())
	fmt.Printf("Type of DB: %s\n", T)
	isU = S.IsURL()
	fmt.Printf("IsURL: %t\n", isU)
	isS = S.IsSingleFile()
	fmt.Printf("IsSingleFile: %t\n", isS)
	// Output:
	// Typeof Handle: *sqlx.DB
	// Type of DB: sqlite
	// IsURL: false
	// IsSingleFile: true
}

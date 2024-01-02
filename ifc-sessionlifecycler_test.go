package datarepo_test

import (
	"fmt"
	DRS "github.com/fbaube/datarepo/sqlite"
	"io/ioutil"
	"os"
)

// ExampleSessionLifecycle tests (for SQLite) the interface:
//   - Open() error
//   - IsOpen() bool
//   - Verify() error
//   - Flush() error
//   - Close() error
//
// .
func ExampleSessionLifecycle() {
	var F *os.File
	// S implements repo.SessionLifecycle
	var S *DRS.SqliteRepo
	var e error
	// func TempFile(dir, pattern string) (f *os.File, err error)
	F, e = ioutil.TempFile("", "*.db")
	if e != nil {
		panic(e)
	}
	// We would like to print out the path, but we
	// can't, because it is random but the output
	// has to be matched for the test to pass.
	// fmt.Println("Path:", F.Name())
	F.Close()

	// Now the test begins.
	S, e = DRS.OpenRepoAtPath(F.Name())
	if e != nil {
		panic(e)
	}
	fmt.Printf("IsOpen #1: %t\n", S.IsOpen())
	e = S.Verify()
	if e != nil {
		panic(e)
	}
	e = S.Close()
	if e != nil {
		panic(e)
	}
	fmt.Printf("IsOpen #2: %t\n", S.IsOpen())
	// Output:
	// IsOpen #1: true
	// IsOpen #2: false
}

package datarepo

import(
	"fmt"
	DRS "github.com/fbaube/datarepo/sqlite"
)

func init() {
  // var SimplR SimpleRepo
     var pSqltR *DRS.SqliteRepo
     pSqltR = new(DRS.SqliteRepo)
  // SimplR = pSqltR
     fmt.Printf("SqltR: %#v \n", *pSqltR) // , "SimplR", *SimplR)
     /*
     pSR, ok = env.SimpleRepo.(*SqliteRepo)
     if !ok {
     	println("PROBLEM! *SqliteRepo is not SimpleRepo.")
	}
	*/
}
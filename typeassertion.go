package datarepo

import(
	DRS "github.com/fbaube/datarepo/sqlite"
)

func init() {
     var SimplR SimpleRepo
     var pSqltR *DRS.SqliteRepo
     pSqltR = new(DRS.SqliteRepo)
     SimplR = pSqltR
     println("pSqltR", pSqltR, "SimplR", SimplR)
     /*
     pSR, ok = env.SimpleRepo.(*SqliteRepo)
     if !ok {
     	fmt.Printf("PROBLEM! *SqliteRepo is not SimpleRepo.")
	pSR = 
     	*/
	}
package sqlite

import(
	// "database/sql"
	"errors"
	"fmt"
	)

// Begin checks that there is no open Tx and then opens one.
// If one is already open, it returns an error "SQLITE_BUSY".
// (Altho this could easily be replaced with a short wait of
// some tens of milliseconds, and if it is still busy after
// a "long enough" wait, then throw a panic that says that
// we probably have deadlock caused by a failure to end a
// transaction.) 
//
// Note tho that in effect this is a queue for transactions,
// which cannot run concurrently. This is not really optimal, 
// but it could be improved upon in the future.
// .
func (p *SqliteRepo) Begin() error {
     if p.Tx != nil {
     	return errors.New("SqliteRepo.Begin: SQLITE_BUSY")
	}
     var e error 
     p.Tx, e = p.Handle().Begin()
     if e != nil {
     	p.Tx = nil 
     	return fmt.Errorf("SqliteRepo.Begin: Repo.Begin failed: %w", e)
	}
     return nil
}
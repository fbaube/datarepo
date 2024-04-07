package sqlite

import(
	// "database/sql"
	"errors"
	"fmt"
	)

// There is no SQLite pragma for BEGIN IMMEDIATE, 
// possibly because it is per-connection, not per-
// database. func BeginImmed() uses a workaround. 

/*
https://www.hwaci.com/sw/sqlite/lang_transaction.html
If the transaction is immediate, then RESERVED locks are acquired
on all databases as soon as the BEGIN command is executed, without
waiting for the database to be used. After a BEGIN IMMEDIATE, no
other database connection will be able to write to the database
or do a BEGIN IMMEDIATE or BEGIN EXCLUSIVE. Other processes can
continue to read from the database, however.
*/

// Begin checks that there is no open Tx and then opens one.
// If one is already open, it returns an error "SQLITE_BUSY".
// (This could easily be replaced with a short wait of some
// tens of milliseconds, but AFAIK it already does this, and
// if it is still busy after a "long enough" wait, then throw
// a panic that says that we probably have deadlock caused by
// a failure to end a transaction.) 
//
// Note tho that in effect this is a queue for transactions,
// which cannot run concurrently. This is not really optimal, 
// but it could be improved upon in the future. Perhaps an
// explicit queue for transactions. 
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

func (p *SqliteRepo) BeginImmed() error {
     if p.Tx != nil {
     	return errors.New("SqliteRepo.BeginImmed: SQLITE_BUSY")
	}
     var e error 
     p.Tx, e = p.Handle().Begin() 
     if e != nil {
     	p.Tx = nil 
     	return fmt.Errorf("SqliteRepo.BeginImmed: sql.DB.Begin failed: %w", e)
	}
     // https://github.com/mattn/go-sqlite3/issues/400#issuecomment-598953685
     _, e = p.Tx.Exec("ROLLBACK; BEGIN IMMEDIATE") 
     if e != nil {
     	p.Tx = nil 
     	return fmt.Errorf("SqliteRepo.BeginImmed: sql.DB.Exec failed: %w", e)
	}
     return nil
}


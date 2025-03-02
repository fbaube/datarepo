package datarepo

// QuerySpecRunner runs query specs!
//
// https://go.dev/wiki/SQLInterface#executing-queries
//  - ExecContext gets NO rows:
//    result, err := db.ExecContext(ctx,
//   "INSERT INTO users (name, age) VALUES ($1, $2)", "gopher", 27)
//  - type [sql.Result] interface {
//	LastInsertId() (int64, error) ; 
//	RowsAffected() (int64, error) }
//   - QueryRowContext gets 1 row
//   - QueryContext gets N rows
// .
type QuerySpecRunner [keyT comparable] interface {
     	// RunQuery0 is Exec() is get NO rows
	// RunQuerySpec0(*DRU.QuerySpec) (any, error) 
	// RunUniquery is QueryRowContext() is get ONE row (like: by_ID)
	// RunUniquerySpec(*UniquerySpec[keyT]) (any, error)
	RunUniquerySpec(*UniquerySpec) (any, error)
	// RunQuery is QueryContext is get N rows 
	RunQuerySpec(*QuerySpec) ([]any, error) 
}

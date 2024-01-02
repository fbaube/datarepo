package utils

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	S "strings"
)

// This file is interesting for a few reasons.
// 1) These funcs are all top-level funcs, NOT methods on structs.
//    This is the prefer'd way to do things with generics, says Go.
// 2) In the present implementation, each struct has its own package,
//    and in the package, SQL statements get prepared and specify
//    the DB table, which "should" correlate to the specific struct
//    defined in the package, which is returned to the caller. This
//    is quite independent tho of the effects of generics.
// 3) Note that there is NO use of SQL NULL. Not only does it
//    mess up pure relational logic (*gasp*), it also causes
//    crappy, complicated code. We require instead that missing
//    data is marked and handled at the app level. In debugging
//    statements we use phi "ϕ" to mark empty strings.

// dump writes to os.Stdout
func dump(op string, q string, args ...any) {
	fmt.Printf("db.%s args: %s \n   %s \n",
		op, StringifyAnyArgs(args...), q)
}

// Create executes the INSERT statement found in q.
// It returns the last inserted ID if any, as an int64.
// Thus it does not need (or use) generics.
// .
func Create(ctx context.Context, q string, args ...any) (int64, error) {
	dump("CREATE", q, args...)
	res, err := R.DB.ExecContext(ctx, q, args...)
	if err != nil {
		return 0, fmt.Errorf("db.Create() failed: %w \n\t (%s)", err, q)
	}
	return res.LastInsertId()
}

// GetRow returns a single row from the DB as type T.
// If no row is found, it returns the T zero value and ErrNotFound.
//
// Example: retval,e :=
// db.GetRow[User](ctx, q_getByID, id) // retval.(User) is true.
// .
func GetRow[T any, PT Row[T]](ctx context.Context, q string, args ...any) (T, error) {
	dump("GetROW", q, args...)
	row := R.DB.QueryRowContext(ctx, q, args...)
	var t T
	ptr := PT(&t)
	if err := row.Scan(ptr.PtrFields()...); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return t, ErrNotFound
		}
		return t, fmt.Errorf("db.GetRow(%T) "+
			"row.Scan error: %w \n\t (%s)", t, err, q)
	}
	return t, nil
}

// GetField returns a single column value from the DB as type T.
// If no column is found, it returns the T zero value with no error.
//
// Example:
//   - func GetField [T Field](ctx, q, args ...any) (T, error)
//   - return db.GetField[int64](ctx, q_getCount)
//
// .
func GetField[T Field](ctx context.Context, q string, args ...any) (T, error) {
	dump("GetFLD", q, args...)
	row := R.DB.QueryRowContext(ctx, q, args...)
	var t T
	if err := row.Scan(&t); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return t, ErrNotFound
		}
		return t, fmt.Errorf("db.GetField(%T) "+
			"row.Scan error: %w \n\t (%s)", t, err, q)
	}
	return t, nil
}

// GetRows returns rows from the DB as type []T.
// T must implement PtrFields() with a ptr receiver type.
func GetRows[T any, PT Row[T]](ctx context.Context, q string, args ...any) ([]T, error) {
	dump("GetRWS", q, args...)
	rows, err := R.DB.QueryContext(ctx, q, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		var t T
		return nil, fmt.Errorf("db.GetRows(%T) "+
			"QueryContext error: %w \n\t (%s)", t, err, q)
	}
	defer func() { _ = rows.Close() }()

	var result []T
	for rows.Next() {
		var t T
		ptr := PT(&t)
		if err := rows.Scan(ptr.PtrFields()...); err != nil {
			return nil, fmt.Errorf("db.GetRows(%T) "+
				"row.Scan error: %w \n\t (%s)", t, err, q)
		}
		result = append(result, t)
	}
	if err := rows.Err(); err != nil {
		var t T
		return nil, fmt.Errorf("db.GetRows(%T) "+
			"rows.Err(): %w \n\t (%s)", t, err, q)
	}
	return result, nil
}

// GetFields returns fields from the DB as type []T.
// T must satisfy the Field constraint: T must be
// a basic Go data type. As the main use case, an
// int64 primary index satisfies this constraint.
//
// Example:
//   - func GetFields[T Field](ctx, q, args ...any) ([]T, error)
//   - return db.GetFields[int64](ctx, q_getByIDsWithBio)
//
// .
func GetFields[T Field](ctx context.Context, q string, args ...any) ([]T, error) {
	dump("GetBSX", q, args...)
	rows, err := R.DB.QueryContext(ctx, q, args...)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		var t T
		return nil, fmt.Errorf("db.GetFields(%T) "+
			"QueryContext error; %w \n\t (%s)", t, err, q)
	}
	defer func() { _ = rows.Close() }()

	var result []T
	for rows.Next() {
		var t T
		if err := rows.Scan(&t); err != nil {
			return nil, fmt.Errorf("db.GetFields(%T) "+
				"row.Scan error: %w \n\t (%s)", t, err, q)
		}
		result = append(result, t)
	}
	if err := rows.Err(); err != nil {
		var t T
		return nil, fmt.Errorf("db.GetFields(%T) "+
			"rows.Err(): %w \n\t (%s)", t, err, q)
	}
	return result, nil
}

// InArgs returns ( placeholders and args )
// as formatted for a WHERE IN clause. For
// example, calling InArgs([]int{1,2,3})
// will return ("?,?,?", []any{1,2,3}).
func InArgs[T Field](inArgs []T) (string, []any) {
	// fmt.Printf("db.InArgs:  In: %+v \n", inArgs)
	outArgs := make([]any, len(inArgs))
	for i, a := range inArgs {
		outArgs[i] = a
	}
	s := S.Repeat("?,", len(outArgs)-1) + "?"
	// fmt.Printf("db.InArgs: Out: \"%s\":%+v \n", s, outArgs)
	return s, outArgs
}

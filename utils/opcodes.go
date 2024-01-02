package utils

/*
CURRENT API in DB (based on the article)

 - func Create                      (ctx, q, args ...any) (int64, error)
 - func GetRow   [T any, PT Row[T]] (ctx, q, args ...any)     (T, error)
 - func GetRows  [T any, PT Row[T]] (ctx, q, args ...any)   ([]T, error)
 - func GetField [T Field]          (ctx, q, args ...any)     (T, error)
 - func GetFields[T Field]          (ctx, q, args ...any)   ([]T, error)
 - func InArgs   [T Field]          (tt []T)             (string, []any)

NOTE: Drop "Row", even tho it refers to the constraint.

CURRENT API in "user/" (based on the article)

 - Create: func Create(ctx context.Context, u User) (User, error)
 - GetRow[User]: func GetByID(ctx context.Context, id int64) (User, error)
 - GetRows[User]: func GetByIDs(ctx context.Context, ids []int64) ([]User, error)
 - GetField[int64]: func Count(ctx context.Context) (int64, error)
 - GetFields[int64]: func GetIDsWithBio(ctx context.Context) ([]int64, error)

Our verbs (for e.g. a TypesRepo) are:

 - Get (retrieve,find,fetch,list) (& count)
 - Add (create,insert)
 - Mod (update)
 - Del (remove)

Basic grammar:

 - Add,Mod,Del,Get
 - One of:
   - -,All=Mult,Everything[T]
   - Where[s][T]
   - ByID[s][T]

*/

type DbOp string

const (
	// Odd ops
	OpCreateTable DbOp = "CreateTable"
	OpCount            = "Count" // () (int, error)
	// Ops on single items:
	//  - 5 functions (2 generic) - lacking ModByID
	//  - All operate by using the ID, except of course Add(T)
	OpAdd     = "Add"     // (T) (int, error)
	OpDel     = "Del"     // (T) error
	OpMod     = "Mod"     // (T) error
	OpGetByID = "GetByID" // [T](int) (T, error)
	OpDelByID = "DelByID" // [T](int) error
	// Ops on multiple items where the count is KNOWN:
	//  - 6 functions (3 generic)
	OpGetByIDs = "GetByIDs" // [T]([]int) ([]T, []error)
	OpDelByIDs = "DelByIDs" // [T]([]int) []error
	OpModByIDs = "ModByIDs" // [T]([]int, actions ...string) []error
	OpAddAll   = "AddAll"   // ([]T) ([]int, []error)
	OpDelAll   = "DelAll"   // ([]T) []error
	OpModAll   = "ModAll"   // ([]T) []error
	// Ops on multiple items (WHERE, Everything)
	// where the count is NOT known:
	//  - 6 functions (all generic) - lacking ModEverything
	OpGetIDsWhere   = "GetIDsWhere"   // [T](cond string) ([]int64, error)
	OpGetWhere      = "GetWhere"      // [T](cond string) ([]T, error)
	OpDelWhere      = "DelWhere"      // [T](cond string) ([]int64, error)
	OpModWhere      = "ModWhere"      // [T](cond s, action s) ([]int64, e)
	OpGetEverything = "GetEverything" // [T]() ([]T, error)
	OpDelEverything = "DelEverything" // [T]() (error)
)

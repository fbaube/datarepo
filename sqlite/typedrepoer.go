package sqlite

import(
	"database/sql"
	DRU "github.com/fbaube/datarepo/utils"
	RM "github.com/fbaube/datarepo/rowmodels"
	)

type Zork struct {
	I int
	S string
}

// TypedRepoer does
//   - Insert (add,Create,save,store)
//   - Select (get,Read,  find,list)
//   - Update (mod,Update,modify)
//   - Delete (del,Delete,remove)
// . 
type TypedRepoer[T any] interface {
     	RM.RowModeler // Must satisfy this interface
	*T // Must be a pointer
	
	Insert(T)  (int, error)
	Update(T)       (error)
	Delete(T)       (error)
	DeleteByID(int) (error)
	SelectByID(int) (T, error)

	SelectByQuery(DRU.QuerySpec)    (T,   error)
	SelectsByQuery(DRU.QuerySpec) ([]T, []error)

	Inserts([]T) ([]int, []error)
	Updates([]T)        ([]error)
	Deletes       ([]T) ([]error)
	DeleteByIDs([]int)  ([]error)
	SelectByIDs([]int)  ([]T, []error)
	
	SelectAll() ([]T, error)
	DeleteAll()      (error)
}

/*
type TypedRepoer[T any] interface {
	Add(T) (int, error)
	Mod(T) (T, error)
	Del(T) error
	GetByID(int) (T, error)
	GetByIDs([]int) ([]T, error)
	// ModByID(int) (T, error) // WTH? 
	DelByID(int) error
	GetAll() ([]T, []error)
	// AddAll([]T) ([]int, []error)
	ModAll([]T) ([]T, []error)
	DelAll([]T) []error
	DelAllByID([]int) []error
}
*/

func (p *Zork) GetByID(id int) Zork {
	Z := new(Zork)
	Z.I = id
	return *p
}

// file:///Users/fbaube/src/M/utils/repo/Abt_GenericsRepo/Generic-Repos-in-Go.html

type TypedRepo[T any] struct {
	db *sql.DB
}

func (r *TypedRepo[T]) Create(t T) { // error {
	// return r.db.Create(&t).Error
	// FAIL?! return T(0)
}

func (r *TypedRepo[T]) Get(id uint) (*T, error) {
	var t T
	// err := r.db.Where("id = ?", id).First(&t).Error
	return &t, nil // err
}

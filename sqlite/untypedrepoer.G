package sqlite

type UntypedRepoer interface {
	GetByID() (any, error)
	Add(any) (int, error)
	ModByID(any) (any, error)
}

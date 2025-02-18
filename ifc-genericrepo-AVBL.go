package datarepo

type GenericRepo[T any] interface {
	Create(T) T
	GetAll() []T
	GetByID(uint) (T, error)
	UpdateByID(uint, T) (T, error)
	DeleteByID(uint) (bool, error)
	GetSome(string) []T
}

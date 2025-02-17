package common

type Document interface {
	Data() map[string]interface{}
	Id() string
	Delete() error
	Set(any) error
}

type Collection interface {
	Query
	Doc(string) (Document, error)
	Add(any) (Document, error)
}

type Query interface {
	Where(string, string, interface{}) Query
	Limit(int) Query
	Select(...string) Query
	OrderBy(string, bool) Query
	Get() ([]Document, error)
}

type DAO interface {
	Collection(string) Collection
}

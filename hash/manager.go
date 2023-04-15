package hash

type Manager interface {
	Hash(data string) (string, error)
	Validate(data string, hash string) error
}

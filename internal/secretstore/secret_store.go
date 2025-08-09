package secretstore

type SecretStore interface {
	Store(key string, secret []byte) error
	Get(key string) ([]byte, error)
	Remove(key string) error
}

package storage

// Adapter is an interface for pluggable components that store log messages.
type Adapter interface {
	Start()
	Write(string, string) error
	Read(string, int) ([]string, error)
	Destroy(string) error
	Reopen() error
	Stop()
}

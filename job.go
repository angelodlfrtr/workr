package workr

// Job represent a job
type Job interface {
	Name() string
	New() Job
	Run(*Workr) error
	Bytes() ([]byte, error)
	Load(data []byte) error
}

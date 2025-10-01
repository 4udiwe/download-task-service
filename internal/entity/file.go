package entity

type FileStatus string

const (
	FilePending FileStatus = "pending"
	FileRunning FileStatus = "in_progress"
	FileDone    FileStatus = "done"
	FileFailed  FileStatus = "failed"
)

type File struct {
	URL    string     `json:"url"`
	Status FileStatus `json:"status"`
	Path   string     `json:"path,omitempty"`
	Error  string     `json:"error,omitempty"`
}

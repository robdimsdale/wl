package wl

// FolderRevision contains information about the revision of folder.
type FolderRevision struct {
	ID         uint   `json:"id" yaml:"id"`
	TypeString string `json:"type" yaml:"type"`
	Revision   uint   `json:"revision" yaml:"revision"`
}

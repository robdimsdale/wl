package wundergo

// FolderRevision contains information about the revision of folder.
type FolderRevision struct {
	ID         uint   `json:"id"`
	TypeString string `json:"type"`
	Revision   uint   `json:"revision"`
}

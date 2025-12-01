package dto

import "github.com/google/uuid"

// FileResponse adalah struktur yang HANYA berisi field yang ingin ditampilkan.
type FileResponse struct {
	ID		uuid.UUID `json:"id"`
	Filename	string `json:"filename"`
	Filetype	string `json:"filetype"`
	Filepath	string `json:"filepath"`
}
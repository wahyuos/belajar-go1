package service

import (
	"bejalar-dasar/internal/model"
	"bejalar-dasar/internal/repository"
	"errors"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"gorm.io/gorm"
)

const maxFileSize = 500 * 1024 // 500 KB
const uploadDir = "./uploads"   // Direktori penyimpanan file fisik

type FileService interface {
	// PrepareFile: Hanya melakukan validasi dan membuat model file (tanpa menyimpan ke DB)
	PrepareFile(fileHeader *multipart.FileHeader) (*model.File, error)
	// SaveMetadata: Tugas tunggal untuk menyimpan model yang sudah lengkap ke DB
	SaveMetadata(fileRecord *model.File) error 
	GetFileRecord(id string) (*model.File, error)
}

type fileService struct {
	repo repository.FileRepository
}

func NewFileService(repo repository.FileRepository) FileService {
	return &fileService{repo}
}

// PrepareFile: Logika Validasi dan Penamaan File
func (s *fileService) PrepareFile(fileHeader *multipart.FileHeader) (*model.File, error) {
	// 1. Validasi Ukuran Maksimal (500 KB)
	if fileHeader.Size > maxFileSize {
		return nil, errors.New("file size exceeds 500 KB limit")
	}

	// 2. Validasi Tipe File (Gambar atau PDF)
	fileExtension := strings.ToLower(filepath.Ext(fileHeader.Filename))
	filetype := ""

	switch fileExtension {
	case ".jpg", ".jpeg", ".png", ".gif":
		filetype = "image"
	case ".pdf":
		filetype = "pdf"
	default:
		return nil, errors.New("file type not supported. only accepts images or PDFs")
	}

	// 3. Penamaan File Unik dan Path
	fileName := strings.ReplaceAll(time.Now().Format("20060102150405.000"), ".", "") + fileExtension
	filePath := filepath.Join(uploadDir, fileName)

	// 4. Buat dan Kembalikan Model (Belum disimpan ke DB)
	fileRecord := &model.File{
		Filename: fileHeader.Filename,
		Filetype: filetype,
		Filepath: filePath,
	}

	return fileRecord, nil
}

// SaveMetadata: Tugas tunggal untuk menyimpan metadata file
func (s *fileService) SaveMetadata(fileRecord *model.File) error {
    // Panggil Repository untuk menyimpan ke DB
	return s.repo.Create(fileRecord)
}

func (s *fileService) GetFileRecord(id string) (*model.File, error) {
	// Panggil Repository
	file, err := s.repo.GetByID(id) 
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, errors.New("file not found")
		}
		return nil, errors.New("failed to retrieve file data: " + err.Error())
	}
	return file, nil
}
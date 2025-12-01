package handler

import (
	"bejalar-dasar/internal/service"
	"bejalar-dasar/pkg/response"
	"fmt"
	"net/http"
	"os"
	"path/filepath"

	"github.com/gin-gonic/gin"
)

// Handler
type FileHandler struct {
	FileService service.FileService
}

// UploadFileHandler: Endpoint POST untuk upload file
func (h *FileHandler) UploadFileHandler(c *gin.Context) {
	// 1. Validasi Inputan HTTP (Handler Responsibility)
	fileHeader, err := c.FormFile("file") 
	if err != nil {
		// response.Error(c, http.StatusBadRequest, "Validation error", gin.H{"error": "field 'file' tidak ditemukan atau gagal diambil"})
		response.Error(c, http.StatusBadRequest, "Validation error", err.Error())
		return
	}

	// 2. Panggil Service untuk Validasi Bisnis dan Persiapan Model
	fileRecord, err := h.FileService.PrepareFile(fileHeader)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Validation error", err.Error()) // Error validasi bisnis
		return
	}
    
	// 3. Simpan File Fisik (Aksi I/O File)
	filePath := fileRecord.Filepath 
	
    // Pastikan direktori dibuat sebelum menyimpan
    if err := os.MkdirAll(filepath.Dir(filePath), os.ModePerm); err != nil {
        // c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal membuat direktori upload"})
		response.Error(c, http.StatusInternalServerError, "Failed to create upload directory", err.Error())
        return
    }

	if err := c.SaveUploadedFile(fileHeader, filePath); err != nil {
		// c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal menyimpan file fisik"})
		response.Error(c, http.StatusInternalServerError, "Failed to upload file", err.Error())
		return
	}

	// 4. Panggil Service untuk Menyimpan Metadata ke DB
	if err := h.FileService.SaveMetadata(fileRecord); err != nil {
		// ROLLBACK FILE FISIK JIKA DB GAGAL
		os.Remove(filePath) 
		
		// 5. Response Error
		// c.JSON(http.StatusInternalServerError, gin.H{"error": "gagal menyimpan metadata ke DB. File fisik dihapus. Detail: " + err.Error()})
		response.Error(c, http.StatusInternalServerError, "Failed to save file", err.Error())
		return
	}

	// 6. Response Sukses (Handler Responsibility)
	response.Success(c, http.StatusOK, "File uploaded successfully", nil)
}

func (h *FileHandler) DownloadFileHandler(c *gin.Context) {
	// 1. Ambil File ID dari URL parameter
	id := c.Param("id")
    
	// 2. Panggil Service untuk mendapatkan metadata
	fileRecord, err := h.FileService.GetFileRecord(id)
	if err != nil {
		// Error sudah dihandle oleh service, termasuk RecordNotFound
		// c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		response.Error(c, http.StatusNotFound, "File not found", err.Error())
		return
	}

	// 3. Tentukan Path File Fisik
	// Catatan: fileRecord.Filepath adalah path di dalam container (misalnya /app/uploads/unique_name.jpg)
	filePath := fileRecord.Filepath
	
    // 4. Set Content Type (Penting!)
    // Ini memberi tahu browser cara menampilkan file (display/download)
	contentType := ""
	switch fileRecord.Filetype {
	case "image":
		// Gin secara otomatis akan mencoba menebak tipe MIME berdasarkan ekstensi
		contentType = "image/" + filepath.Ext(fileRecord.Filename)[1:] // contoh: image/jpeg
	case "pdf":
		contentType = "application/pdf"
	default:
		contentType = "application/octet-stream" // Default untuk download paksa
	}
    
    // 5. Streaming File
    // Gin menyediakan helper yang aman untuk melayani file
    // Content-Disposition: 'inline' akan mencoba menampilkan file di browser (seperti gambar dan PDF).
    // Jika ingin download paksa, gunakan 'attachment; filename='
    
    c.Header("Content-Type", contentType)
    c.Header("Content-Disposition", fmt.Sprintf("inline; filename=%s", fileRecord.Filename))
    
    // Perintah Gin untuk menyajikan file yang tersimpan di disk
    c.File(filePath)
    
    
}
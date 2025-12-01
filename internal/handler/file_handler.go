package handler

import (
	"bejalar-dasar/internal/service"
	"bejalar-dasar/pkg/response"
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
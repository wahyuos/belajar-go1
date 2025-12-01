package main

import (
	"bejalar-dasar/internal/handler"
	"bejalar-dasar/internal/model"
	"bejalar-dasar/internal/repository"
	"bejalar-dasar/internal/service"
	"bejalar-dasar/pkg/config"
	"bejalar-dasar/pkg/database"
	"bejalar-dasar/pkg/middleware"
	"bejalar-dasar/pkg/response"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	"golang.org/x/crypto/bcrypt"
)

func main() {
	// 1. Load Env
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system env")
	}

	// 2. Connect DB
	db := database.GetDatabaseConnection()

	// buat user login
	// Hash Password
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte("password123"), bcrypt.DefaultCost)

	// Buat user
	user := model.User{
		Base: model.Base{
			ID: uuid.MustParse("123e4567-e89b-12d3-a456-426614174000"),
		},
		Username: "admin",
		Password: string(hashedPassword),
	}

	// Simpan ke database (akan otomatis di-hash)
	db.FirstOrCreate(&user, model.User{Username: "admin"})

	// 3. Auto Migration
	err := db.AutoMigrate(
		&model.User{},
		&model.Category{}, 
		&model.Article{},
		&model.File{},
	)
	if err != nil {
		log.Fatal("Migration failed:", err)
	}

	// 4. Setup Layer (Dependency Injection)
	// -- User Module --
    userRepo := repository.NewUserRepository(db)
	userService := service.NewUserService(userRepo)
    authHandler := handler.NewAuthHandler(userService)

	// -- Category Module --
	catRepo := repository.NewCategoryRepository(db)
	catService := service.NewCategoryService(catRepo)
	catHandler := handler.NewCategoryHandler(catService)

	// -- Article Module --
	articleRepo := repository.NewArticleRepository(db)
	// Perhatikan: ArticleService butuh catRepo juga untuk validasi
	articleService := service.NewArticleService(articleRepo, catRepo) 
	articleHandler := handler.NewArticleHandler(articleService)

	// -- Category Module --
	fileRepo := repository.NewFileRepository(db)
	fileService := service.NewFileService(fileRepo)
	fileHandler := handler.FileHandler{FileService: fileService}

	// Note: Anda perlu mengimplementasikan repo/service/handler untuk Article 
	// dengan pola yang sama persis seperti Category di atas.

	// 5. Setup Router
	// Muat Konfigurasi
    cfg, err := config.GinMode()
    if err != nil {
        log.Fatalf("Failed to load GinMode configuration: %v", err)
    }

    // Terapkan Gin Mode berdasarkan Konfigurasi (.env)
    // gin.SetMode(gin.ReleaseMode) atau gin.SetMode(gin.DebugMode)
    if cfg.GinMode == gin.ReleaseMode {
        gin.SetMode(gin.ReleaseMode)
        log.Println("RELEASE mode.")
    } else {
        // Default ke Debug jika tidak diatur atau nilainya bukan 'release'
        gin.SetMode(gin.DebugMode)
        log.Println("DEBUG mode (default).")
    }
	
	// jika dev
	r := gin.Default()

	// index
	r.GET("/", func (ctx *gin.Context)  {
		response.Success(ctx, http.StatusOK, "Welcome", nil)
	})

	api := r.Group("/api/v1")
	{
		// health
		api.GET("/health", func (ctx *gin.Context)  {
			response.Success(ctx, http.StatusOK, "Everything ok", nil)
		})

		// RUTE AUTHENTIKASI (TIDAK PERLU MIDDLEWARE)
        api.POST("/auth/login", authHandler.Login)

		// Rute yang MEMERLUKAN otentikasi
        // Buat Grup baru dan terapkan AuthMiddleware
        protected := api.Group("/", middleware.AuthMiddleware())
        {
			files := protected.Group("/files")
			{
				// Endpoint untuk upload file
				files.POST("/upload", fileHandler.UploadFileHandler)
				files.GET("/:id", fileHandler.DownloadFileHandler)
			}

			categories := protected.Group("/categories")
			{
				categories.GET("/", catHandler.GetAll)
				categories.POST("/", catHandler.Create)
				categories.PUT("/:id", catHandler.Update)
				categories.DELETE("/:id", catHandler.Delete)
			}
			
			articles := protected.Group("/articles")
			{
				articles.GET("/", articleHandler.GetAll)
				articles.GET("/:id", articleHandler.GetByID)
				articles.POST("/", articleHandler.Create)
				articles.PUT("/:id", articleHandler.Update)
				articles.DELETE("/:id", articleHandler.Delete)
			}
		}
	}

	// 6. Run Server
	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}
	r.Run(":" + port)
}
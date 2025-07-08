package main

import (
	"barberia-api-class/internal/domain"
	"barberia-api-class/internal/repository"
	"barberia-api-class/internal/service"
	"log"
	"time"

	httptrans "barberia-api-class/internal/transport/http"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

func initConfig() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AutomaticEnv()

	// Valores por defecto
	viper.SetDefault("PORT", "8080")
	viper.SetDefault("DB_URL", "postgres://barberia_user:barberia_pass@127.0.0.1:5432/barberia-db?sslmode=disable")
	viper.SetDefault("GIN_MODE", "debug")

	if err := viper.ReadInConfig(); err != nil {
		log.Printf("No se encontro config.yaml, usando variables de entorno y valores por defecto: %v", err)
	}
}

func main() {
	initConfig()

	dbURL := viper.GetString("DB_URL")
	port := viper.GetString("PORT")
	ginMode := viper.GetString("GIN_MODE")

	// Configurar modo de Gin
	gin.SetMode(ginMode)

	log.Println("Iniciando Barberia API")
	log.Printf("Prueto: %s", port)
	log.Printf("Modo Gin: %s", ginMode)

	// Conexion a PostgresSQL con GORM v2
	db, err := gorm.Open(postgres.Open(dbURL), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})
	if err != nil {
		log.Fatalf("Fallo conexion a la base de datos: %v", err)
	}

	// Ejecutar migraciones
	if err := db.AutoMigrate(&domain.Appointment{}, &domain.Product{}); err != nil {
		log.Fatalf("Error en migraciones: %v", err)
	}

	// Configurar pool de conexiones
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Error obteniendo instancia de la base de datos: %v", err)
	}
	sqlDB.SetMaxIdleConns(10)
	sqlDB.SetMaxOpenConns(100)
	sqlDB.SetConnMaxLifetime(time.Hour)

	// Inyeccion de dependencias
	log.Println("Configurando Dependencias...")
	apptRepo := repository.NewGormAppointmentRepo(db)
	prodRepo := repository.NewGormProductRepo(db)
	apptSvc := service.NewAppointmentService(apptRepo, prodRepo)
	prodSvc := service.NewProductService(prodRepo)

	// Aranque de Gin
	log.Println("Configurando rutas...")
	router := httptrans.NewRouter(apptSvc, prodSvc)

	log.Printf("Servidor de Barberia escuchando en puerto :%s", port)
	log.Printf("Health check disponible en: http://localhost:%s/health", port)
	log.Printf("API endpoints en: http://localhots:%s/api/v1", port)

	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Error Iniciando Servidor: %v", err)
	}
}

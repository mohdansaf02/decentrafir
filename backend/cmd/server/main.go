package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/decntrafir/backend/internal/config"
	"github.com/decntrafir/backend/internal/controllers"
	"github.com/decntrafir/backend/internal/database"
	"github.com/decntrafir/backend/internal/routes"
	"github.com/decntrafir/backend/internal/services"
	"github.com/decntrafir/backend/pkg/blockchain"
	"github.com/decntrafir/backend/pkg/ipfs"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

// @title Blockchain FIR Management API
// @version 1.0
// @description Tamper-proof FIR management with blockchain and IPFS
// @host localhost:8080
// @BasePath /
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatal(err)
	}
	gin.SetMode(cfg.GinMode)

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	db, err := database.Connect(ctx, cfg.MongoURI, cfg.MongoDB)
	if err != nil {
		log.Fatalf("mongodb: %v", err)
	}
	defer db.Disconnect(context.Background())

	chain, err := blockchain.NewClient(cfg.BlockchainRPC, cfg.ContractAddress, cfg.DeployerPrivateKey)
	if err != nil {
		log.Printf("blockchain client warning: %v (running in mock mode)", err)
		chain = &blockchain.Client{}
	}

	ipfsClient := ipfs.NewPinataClient(cfg.PinataJWT, cfg.PinataAPIKey, cfg.PinataSecret)
	auditSvc := services.NewAuditService(db)
	notificationSvc := services.NewNotificationService(db, cfg)
	authSvc := services.NewAuthService(db, cfg, auditSvc)
	firSvc := services.NewFIRService(db, chain, auditSvc, notificationSvc)
	evidenceSvc := services.NewEvidenceService(db, ipfsClient, auditSvc, cfg.MaxUploadMB)
	userSvc := services.NewUserService(db, auditSvc)

	handlers := routes.Handlers{
		Auth:     controllers.NewAuthController(authSvc),
		FIR:      controllers.NewFIRController(firSvc),
		Evidence: controllers.NewEvidenceController(evidenceSvc),
		Admin:    controllers.NewAdminController(userSvc, auditSvc, notificationSvc),
	}

	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())
	r.Use(cors.New(cors.Config{
		AllowOrigins:     cfg.CORSOrigins,
		AllowMethods:     []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	routes.Setup(r, cfg, handlers)

	srv := &http.Server{Addr: ":" + cfg.Port, Handler: r}
	go func() {
		log.Printf("server listening on :%s", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatal(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer shutdownCancel()
	_ = srv.Shutdown(shutdownCtx)
}

package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/go-chi/chi"
	"github.com/go-chi/chi/middleware"

	auth "github.com/DrAnonymousNet/foodshare/Auth"
	core "github.com/DrAnonymousNet/foodshare/Core"
	foodshare "github.com/DrAnonymousNet/foodshare/FoodShareApp"
	notifications "github.com/DrAnonymousNet/foodshare/Notifications"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/swaggo/http-swagger/example/go-chi/docs" // docs is generated by Swag CLI, you have to import it.
	httpSwagger "github.com/swaggo/http-swagger/v2"
	"gorm.io/gorm"
)

func GetModels() []interface{} {
	Models := []interface{}{
		auth.User{},
		auth.JwtToken{},
		notifications.Notification{},
		foodshare.DonationRequest{},
		foodshare.Donation{},
	}
	return Models
}

var db *gorm.DB // Assume you have a GORM database connection

// @title           Swagger Example API
// @version         1.0
// @description     This is a sample server celler server.
// @termsOfService  http://swagger.io/terms/

// @contact.name   API Support
// @contact.url    http://www.swagger.io/support
// @contact.email  support@swagger.io

// @license.name  Apache 2.0
// @license.url   http://www.apache.org/licenses/LICENSE-2.0.html

// @host      localhost:8000
// @BasePath  /api/v1

// @securityDefinitions.basic  BasicAuth

// @externalDocs.description  OpenAPI
// @externalDocs.url          https://swagger.io/resources/open-api/
func main() {
	godotenv.Load(".env")

	portString := os.Getenv("PORTSTRING")

	db, err := core.SetupDatabase()
	if err != nil {
		log.Fatal(err)
	}

	// Create ENUM types in the database
	db.Debug().Exec(`
		DO $$ BEGIN
			CREATE TYPE enum_donatable_obj_type AS ENUM (
				'FoodStuff', 'Cloths', 'MedicalSupplies',
				'SchoolSupplies', 'PersonalCareSupplies', 'BooksAndToys'
			);
		EXCEPTION
			WHEN duplicate_object THEN null;
		END $$;
		
		DO $$ BEGIN
			CREATE TYPE enum_request_status_type AS ENUM (
				'PartiallyFulfilled', 'FullyFulfilled'
			);
		EXCEPTION
			WHEN duplicate_object THEN null;
		END $$;
		
		DO $$ BEGIN
			CREATE TYPE enum_request_from_type AS ENUM (
				'WareHouse', 'Community'
			);
		EXCEPTION
			WHEN duplicate_object THEN null;
		END $$;
		
		DO $$ BEGIN
			CREATE TYPE enum_donation_status_type AS ENUM (
				'Pending', 'PickedUp'
			);
		EXCEPTION
			WHEN duplicate_object THEN null;
		END $$;
		
		DO $$ BEGIN
			CREATE TYPE enum_gender AS ENUM (
				'Male', 'Female'
			);
		EXCEPTION
			WHEN duplicate_object THEN null;
		END $$;`)

	db.AutoMigrate(GetModels()...)

	router := chi.NewRouter()
	router.Use(cors.Handler(
		cors.Options{
			AllowedOrigins:   []string{"https://*", "http://*"},
			AllowedMethods:   []string{"GET", "POST", "DELETE", "PUT", "PATCH", "OPTIONS"},
			AllowedHeaders:   []string{"*"},
			ExposedHeaders:   []string{"link"},
			AllowCredentials: false,
			MaxAge:           300,
		}),
	)

	router.Use(middleware.RequestID)
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {

		w.Header().Set("Content-Type", "text/plain")

		w.Write([]byte("Hello World!"))

	})
	router.Get("/swagger/*", httpSwagger.Handler(
		httpSwagger.URL("http://127.0.0.1:8000/swagger/doc.json"), //The url pointing to API definition
	))
	router.Post("/api/v1/auth/login/", auth.GenerateJWTTokenHandler)

	router.Post("/api/v1/auth/", auth.CreateUser)

	//router.Mount("/api/v1/auth", auth.AuthRoutes())
	d := foodshare.SetViewSet()
	router.With(auth.UserContext).Post("/api/v1/foodshare/donations/", d.CreateDonation)
	router.With(auth.UserContext).Get("/api/v1/foodshare/donations/", d.ListDonations)
	router.With(auth.UserContext).Get("/api/v1/foodshare/donations/{uid}/", d.GetDonation)
	router.With(auth.UserContext).Patch("/api/v1/foodshare/donations/{uid}/", d.UpdateDonation)

	srv := &http.Server{
		Handler: router,
		Addr:    ":" + portString,
	}
	log.Printf("Starting server at port %v, %v", portString, srv.Addr)

	// Channel to signal server shutdown
	shutdownChan := make(chan struct{})

	// Run the server in a separate goroutine
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("Error starting the server: %s\n", err)
			close(shutdownChan) // Signal shutdown if an error occurs during server startup
		}
	}()
	// Capture OS signals for graceful shutdown
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	// Wait for a signal or server error
	select {
	case <-sigChan:
		fmt.Println("Received shutdown signal. Shutting down server gracefully...")
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := srv.Shutdown(ctx); err != nil {
			fmt.Printf("Error during server shutdown: %s\n", err)
		}
	case <-shutdownChan:
		fmt.Println("Server startup failed. Shutting down...")
	}

	fmt.Println("Server shutdown complete.")
}

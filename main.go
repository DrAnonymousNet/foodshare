package main

import (

	"log"
	"net/http"
	"os"


	"github.com/go-chi/chi"

	"github.com/joho/godotenv"
	auth "github.com/DrAnonymousNet/foodshare/Auth"
	core "github.com/DrAnonymousNet/foodshare/Core"
	"gorm.io/gorm"
)

var db *gorm.DB // Assume you have a GORM database connection


func main(){
	godotenv.Load(".env")

	portString := os.Getenv("portString")

	db, err := core.SetupDatabase()
	if err != nil{
		log.Fatal(err)
	}
	db.AutoMigrate(getModels()...)

	router := chi.NewRouter()

	router.Mount("/v1", auth.AuthRoutes())

	srv := &http.Server{
		Handler: router,
		Addr : ":"+portString,
	}
	srv.ListenAndServe()

}
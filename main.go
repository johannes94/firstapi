package main

import (
	"log"
	"net/http"

	"firstapi/api"
	"firstapi/db"

	"github.com/gorilla/mux"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func main() {
	database := db.InitDB()
	defer database.Close()

	router := mux.NewRouter()
	apiRouter := router.PathPrefix("/api").Subrouter()

	productRepo := db.NewGormProductRepository(database)
	productsController := api.NewProductsController(productRepo)

	apiRouter.HandleFunc("/products", productsController.GetProductList).Methods("GET")
	apiRouter.HandleFunc("/products", productsController.CreateProduct).Methods("POST")
	apiRouter.HandleFunc("/products/{id}", productsController.GetProduct).Methods("GET")
	apiRouter.HandleFunc("/products/{id}", productsController.UpdateProduct).Methods("PUT")
	apiRouter.HandleFunc("/products/{id}", productsController.DeleteProduct).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8080", router))

	// review pluralsight chapter for channels by nigel
}

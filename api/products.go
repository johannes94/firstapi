package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"firstapi/db"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

type productsController struct {
	db *gorm.DB
}

func NewProductsController(db *gorm.DB) *productsController {
	return &productsController{db}
}

func (pc productsController) CreateProduct(w http.ResponseWriter, r *http.Request) {

	var newProduct db.Product
	if err := json.NewDecoder(r.Body).Decode(&newProduct); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Request Body could not be parsed as Product")
		fmt.Println(err)
		return
	}

	if newProduct.Name == nil || newProduct.Price == nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Name or Price not specified in Request")
		return
	}

	if newProduct.Id != 0 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Product id will be set automatically, remove it from the Request Body")
		return
	}

	if err := pc.db.Create(&newProduct).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Internal error with storing the Product")
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Created")
}

func (pc productsController) GetProductList(w http.ResponseWriter, r *http.Request) {
	var products []db.Product

	if err := pc.db.Find(&products).Error; err != nil {
		http.NotFound(w, r)
		fmt.Println(err)
		return
	}

	w.Header().Add("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(products); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error": "Failed to deserialize database response"}`)
		fmt.Println(err)
		return
	}
}

func (pc productsController) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	productId := mux.Vars(r)["id"]
	var product db.Product

	var newProduct db.Product
	if err := json.NewDecoder(r.Body).Decode(&newProduct); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Request Body could not be parsed as Product")
		fmt.Println(err)
		return
	}

	if err := pc.db.First(&product, productId).Error; err != nil {
		if err.Error() == "record not found" {
			http.NotFound(w, r)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Println(err)
		}
		return
	}

	product.Name = newProduct.Name
	product.Price = newProduct.Price

	if err := pc.db.Save(product).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Internal error with storing the Product")
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Updated")
}

func (pc productsController) GetProduct(w http.ResponseWriter, r *http.Request) {

	productId := mux.Vars(r)["id"]
	var product db.Product

	if err := pc.db.First(&product, productId).Error; err != nil {
		if err.Error() == "record not found" {
			http.NotFound(w, r)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Println(err)
		}
		return
	}

	w.Header().Add("Content-Type", "application/json")

	if err := json.NewEncoder(w).Encode(product); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, `{"error": "Failed to deserialize database response"}`)
		fmt.Println(err)
		return
	}

}

func (pc productsController) DeleteProduct(w http.ResponseWriter, r *http.Request) {

	productId := mux.Vars(r)["id"]
	var product db.Product

	if err := pc.db.First(&product, productId).Error; err != nil {
		if err.Error() == "record not found" {
			http.NotFound(w, r)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Println(err)
		}
		return
	}

	if err := pc.db.Delete(&product).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Internal error with deleting the Product")
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Deleted Product")
}

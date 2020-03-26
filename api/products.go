package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"firstapi/db"

	"github.com/gorilla/mux"
)

//ProductsController defines HTTP Handlers for Products
type ProductsController struct {
	productRepository db.ProductRepository
}

//NewProductsController creates a new ProductsController with the DB set
func NewProductsController(pr db.ProductRepository) ProductsController {
	return ProductsController{pr}
}

func (pc ProductsController) CreateProduct(w http.ResponseWriter, r *http.Request) {

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

	if newProduct.ID != 0 {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Product id will be set automatically, remove it from the Request Body")
		return
	}

	if err := pc.productRepository.Create(&newProduct); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Internal error with storing the Product")
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Created")
}

func (pc ProductsController) GetProductList(w http.ResponseWriter, r *http.Request) {
	products, err := pc.productRepository.GetAll()

	if err != nil {
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

func (pc ProductsController) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	productID := mux.Vars(r)["id"]

	var newProduct db.Product
	if err := json.NewDecoder(r.Body).Decode(&newProduct); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		fmt.Fprintf(w, "Request Body could not be parsed as Product")
		fmt.Println(err)
		return
	}

	product, err := pc.productRepository.GetByID(productID)
	if err != nil {
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

	if err := pc.productRepository.Update(&product); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Internal error with storing the Product")
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Updated")
}

func (pc ProductsController) GetProduct(w http.ResponseWriter, r *http.Request) {

	productID := mux.Vars(r)["id"]

	product, err := pc.productRepository.GetByID(productID)

	if err != nil {
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

func (pc ProductsController) DeleteProduct(w http.ResponseWriter, r *http.Request) {

	productID := mux.Vars(r)["id"]

	product, err := pc.productRepository.GetByID(productID)
	if err != nil {
		if err.Error() == "record not found" {
			http.NotFound(w, r)
		} else {
			w.WriteHeader(http.StatusInternalServerError)
			fmt.Println(err)
		}
		return
	}

	if err := pc.productRepository.Delete(&product).Error; err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		fmt.Fprintf(w, "Internal error with deleting the Product")
		return
	}

	w.WriteHeader(http.StatusOK)
	fmt.Fprintf(w, "Deleted Product")
}

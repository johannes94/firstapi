package api

import (
	"encoding/json"
	"firstapi/db"
	"net/http"
	"net/http/httptest"
	"testing"

	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

type mockProductRepository struct {
	ExpectedForGetAll []db.Product
}

func (mp mockProductRepository) GetAll() ([]db.Product, error) {
	return mp.ExpectedForGetAll, nil
}
func (mp mockProductRepository) GetByID(id string) (db.Product, error) {
	return db.Product{}, nil
}
func (mp mockProductRepository) Create(product *db.Product) error {
	return nil
}
func (mp mockProductRepository) Update(product *db.Product) error {
	return nil
}
func (mp mockProductRepository) Delete(product *db.Product) error {
	return nil
}

func TestGetAllProducts(t *testing.T) {
	rr := httptest.NewRecorder()

	productRepo := new(mockProductRepository)
	productRepo.ExpectedForGetAll = []db.Product{
		db.Product{Name: "test", ID: 1, Price: 5},
		db.Product{Name: "test2", ID: 2, Price: 6},
	}

	productsController := ProductsController{productRepo}

	handler := http.HandlerFunc(productsController.GetProductList)
	req, err := http.NewRequest("GET", "/api/procucts", nil)
	if err != nil {
		t.Fatal(err)
	}

	handler.ServeHTTP(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	var products []db.Product

	if err := json.Unmarshal(rr.Body.Bytes(), &products); err != nil {
		t.Errorf("failed to json decode return value got:\n%v", rr.Body.String())
	}

	if len(products) != len(productRepo.ExpectedForGetAll) {
		t.Errorf("slice size is not equal")
		t.Logf("Got: %v", products)
		t.Logf("Want: %v", productRepo.ExpectedForGetAll)
	}

	products[0].Name = "super neuer name"
	for index := range products {
		expected := productRepo.ExpectedForGetAll[index]
		actual := products[index]
		if actual != expected {
			t.Errorf("Product not equal")
			t.Logf("Got: %v", actual)
			t.Logf("Want: %v", expected)
		}
	}

}

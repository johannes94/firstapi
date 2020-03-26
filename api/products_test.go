package api

import (
	"firstapi/db"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	_ "github.com/jinzhu/gorm/dialects/sqlite"
)

func TestGetAllProducts(t *testing.T) {
	rr := httptest.NewRecorder()

	database := db.InitDB()
	defer func() {
		database.Close()
		os.Remove("test.db")
	}()

	productRepo := db.NewGormProductRepository(database)
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

	t.Log(rr.Body)
	// Check the response body is what we expect.
	// expected := `{"alive": true}`
	// if rr.Body.String() != expected {
	//     t.Errorf("handler returned unexpected body: got %v want %v",
	//         rr.Body.String(), expected)
	// }
}

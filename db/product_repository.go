package db

import (
	"fmt"

	"github.com/jinzhu/gorm"
)

type ProductRepository interface {
	GetAll() ([]Product, error)
	GetByID(id string) (Product, error)
	Create(product *Product) error
	Update(product *Product) error
	Delete(product *Product) error
}

type GormProductRepository struct {
	db *gorm.DB
}

func NewGormProductRepository(db *gorm.DB) GormProductRepository {
	return GormProductRepository{db}
}

func (repo GormProductRepository) GetAll() ([]Product, error) {
	var products []Product

	err := db.Find(&products).Error
	return products, err
}

func (repo GormProductRepository) GetByID(id string) (Product, error) {

	var product Product

	err := db.First(&product, id).Error
	return product, err
}

func (repo GormProductRepository) Create(product *Product) error {
	res := db.Create(product)
	err := res.Error

	fmt.Println(res.Value)
	return err
}

func (repo GormProductRepository) Update(product *Product) error {
	return db.Save(product).Error
}

func (repo GormProductRepository) Delete(product *Product) error {
	return db.Delete(product).Error
}

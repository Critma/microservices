package data

import (
	"encoding/json"
	"fmt"
	"io"
)

type Product struct {
	ID          int    `json:"id"`
	Name        string `json:"name" binding:"required,max=120"`
	Price       int    `json:"price" binding:"required,gt=0"`
	Description string `json:"description"`
	CreatedAt   string `json:"-"`
	UpdatedAt   string `json:"-"`
}

var products []*Product = []*Product{
	{ID: 0, Name: "cap1", Description: "burl", CreatedAt: "2025-01-01", UpdatedAt: "2025-01-01"},
	{ID: 1, Name: "cap2", Description: "curl", CreatedAt: "2025-01-01", UpdatedAt: "2025-01-01"},
}

type Products []*Product

func NewProductsList() Products {
	return products
}

func AddProduct(p *Product) {
	p.ID = getNewID()
	products = append(products, p)
}

func UpdateProduct(id int, p *Product) error {
	_, index, err := findProduct(id)
	if err != nil {
		return err
	}
	p.ID = id
	products[index] = p
	return nil
}

var ErrProductNotFound error = fmt.Errorf("product not found")

func findProduct(id int) (*Product, int, error) {
	for i, v := range products {
		if v.ID == id {
			return v, i, nil
		}
	}

	return nil, -1, ErrProductNotFound
}

func getNewID() int {
	prod := products[len(products)-1]
	return prod.ID + 1
}

func (pl *Products) ToJson(w io.Writer) error {
	encoder := json.NewEncoder(w)
	return encoder.Encode(pl)
}

func (p *Product) FromJson(r io.Reader) error {
	decoder := json.NewDecoder(r)
	return decoder.Decode(p)
}

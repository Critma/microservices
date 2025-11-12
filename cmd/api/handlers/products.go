package handlers

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/critma/prodapi/cmd/api/data"
	"github.com/gin-gonic/gin"
)

const (
	PRODUCT_KEY = "product"
)

type Products struct {
	l *log.Logger
}

func NewProducts(l *log.Logger) *Products {
	return &Products{l}
}

// GetProducts godoc
//
//	@Summary		Get a list of products
//	@Description	Get a list of all products
//	@Tags			products
//	@Accept			json
//	@Produce		json
//	@Success		200	{array}		data.Product
//	@Failure		500	{object}	HTTPError
//	@Router			/products/ [get]
func (p *Products) GetProducts(c *gin.Context) {
	pl := data.NewProductsList()
	err := pl.ToJson(c.Writer)

	if err != nil {
		http.Error(c.Writer, "Unable to marshal json", 500)
		return
	}
}

// AddProduct godoc
//
//	@Summary		Add a new product
//	@Description	Add a new product to the DB
//	@Tags			products
//	@Accept			json
//	@Produce		json
//	@Success		201	{object}	MessageWrapper
//	@Failure		500	{object}	HTTPError
//	@Router			/products/ [post]
func (p *Products) AddProduct(c *gin.Context) {
	product := c.Value(PRODUCT_KEY).(*data.Product)

	data.AddProduct(product)

	c.JSON(http.StatusCreated, gin.H{"message": "ok"})
}

// UpdateProduct godoc
//
//		@Summary		Update a product
//		@Description	Update a product in the DB
//		@Tags			products
//		@Accept			json
//		@Produce		json
//	 @Param          id path int true "Product ID"
//		@Success		202	{object}	MessageWrapper
//		@Failure		400	{object}	HTTPError
//		@Failure		404	{object}	HTTPError
//		@Failure		500	{object}	HTTPError
//		@Router			/products/ [put]
func (p *Products) UpdateProduct(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Not valid id"})
		return
	}

	product := c.Value(PRODUCT_KEY).(*data.Product)

	if err := data.UpdateProduct(id, product); err != nil {
		p.l.Println("Cannot updated product")
		switch err {
		case data.ErrProductNotFound:
			c.JSON(http.StatusNotFound, gin.H{"message": "Product not found"})
		default:
			c.JSON(http.StatusBadRequest, gin.H{"message": "Check request body"})
		}
		return
	}

	c.JSON(http.StatusAccepted, gin.H{"message": "ok"})
}

func (p *Products) GetProductMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		product := &data.Product{}

		// if err := product.FromJson(c.Request.Body); err != nil {
		if err := c.ShouldBind(product); err != nil {
			p.l.Println("Cannot unmarshall body, ", err)
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"message": fmt.Sprintf("Cannot unmarshall: %s", err)})
			return
		}

		c.Set(PRODUCT_KEY, product)

		c.Next()
	}
}

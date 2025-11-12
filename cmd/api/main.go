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

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/critma/prodapi/cmd/api/handlers"
	_ "github.com/critma/prodapi/cmd/docs"
)

//	@title			Swagger product API
//	@version		1.0
//	@description	This is a swagger docs for product API.
//	@termsOfService	http://swagger.io/terms/

//	@license.name	MIT
//	@license.url	https://mit-license.org/

// @host		localhost:8080
// @BasePath	/api/v1
func main() {
	l := log.New(os.Stdout, "product-api\t", log.LstdFlags)

	ph := handlers.NewProducts(l)

	r := gin.Default()

	v1 := r.Group("/api/v1")
	{
		v1.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
		products := v1.Group("/products")
		{
			products.GET("/", ph.GetProducts)
			products.POST("/", ph.GetProductMiddleware(), ph.AddProduct)
			product := products.Group("/:id")
			{
				product.Use(ph.GetProductMiddleware())
				product.PUT("/", ph.UpdateProduct)
			}
		}
	}

	s := http.Server{
		Addr:         ":8080",
		Handler:      r,
		IdleTimeout:  120 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
	}

	go func() {
		fmt.Println("Server is started on", s.Addr)
		err := s.ListenAndServe()
		switch err {
		case http.ErrServerClosed:
			fmt.Println("Server is closed")
		default:
			log.Fatal(err)
		}
	}()

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt)
	signal.Notify(sigChan, syscall.SIGTERM)

	sig := <-sigChan
	l.Println("Recieved terminate", sig)
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	s.Shutdown(ctx)
}

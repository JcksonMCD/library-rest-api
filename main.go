package main

import (
	"errors"
	"net/http"

	_ "example/go-rest-api/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

type book struct {
	// Capital field names allow for exportation.
	ID       string `json:"id"`
	Title    string `json:"title"`
	Author   string `json:"author"`
	Quantity int    `json:"quantity"`
}

type ErrorResponse struct {
	Message string `json:"message"`
}

// Data for API as not connecting to a db for this api.
var books = []book{
	{ID: "1", Title: "In Search of Lost Time", Author: "Marcel Proust", Quantity: 2},
	{ID: "2", Title: "The Great Gatsby", Author: "F. Scott Fitzgerald", Quantity: 5},
	{ID: "3", Title: "War and Peace", Author: "Leo Tolstoy", Quantity: 6},
}

// @Summary Get all books
// @Description Retrieve a list of all books
// @Tags Books
// @Accept  json
// @Produce  json
// @Success 200 {array} book
// @Router /books [get]
func getBooks(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, books)
}

// @Summary Get a book by ID
// @Description Retrieve a specific book by its ID
// @Tags Books
// @Accept  json
// @Produce  json
// @Param id path string true "Book ID"
// @Success 200 {object} book
// @Failure 404 {object} ErrorResponse
// @Router /books/{id} [get]
func bookById(c *gin.Context) {
	id := c.Param("id")
	book, err := getBookById(id)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not found!"})
		return
	}

	c.IndentedJSON(http.StatusOK, book)
}

func getBookById(id string) (*book, error) {
	for i, b := range books {
		if b.ID == id {
			return &books[i], nil
		}
	}

	return nil, errors.New("book not found")
}

// @Summary Checkout a book
// @Description Decrease the quantity of a book when checked out
// @Tags Books
// @Accept  json
// @Produce  json
// @Param id query string true "Book ID"
// @Success 200 {object} book
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /checkout [put]
func checkoutBook(c *gin.Context) {
	id, ok := c.GetQuery("id")

	if !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Missing id query parameter."})
		return
	}

	book, err := getBookById(id)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not found!"})
		return
	}

	if book.Quantity <= 0 {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Book not available."})
		return
	}

	book.Quantity -= 1
	c.IndentedJSON(http.StatusOK, book)
}

// @Summary Return a book
// @Description Incrase the quantity of a book when returned
// @Tags Books
// @Accept  json
// @Produce  json
// @Param id query string true "Book ID"
// @Success 200 {object} book
// @Failure 400 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /return [put]
func returnBook(c *gin.Context) {
	id, ok := c.GetQuery("id")

	if !ok {
		c.IndentedJSON(http.StatusBadRequest, gin.H{"message": "Missing id query parameter."})
		return
	}

	book, err := getBookById(id)

	if err != nil {
		c.IndentedJSON(http.StatusNotFound, gin.H{"message": "Book not found!"})
		return
	}

	book.Quantity += 1
	c.IndentedJSON(http.StatusOK, book)
}

// @Summary Create a new book
// @Description Add a new book to the library
// @Tags Books
// @Accept  json
// @Produce  json
// @Param book body book true "Book data"
// @Success 201 {object} book
// @Failure 400 {object} ErrorResponse
// @Router /books [post]
func createBook(c *gin.Context) {
	var newBook book

	if err := c.BindJSON(&newBook); err != nil {
		// BindJSON will handlesending error response so empty return can be used.
		return
	}

	books = append(books, newBook)
	c.IndentedJSON(http.StatusCreated, newBook)
}

func main() {
	router := gin.Default()
	router.GET("/books", getBooks)
	router.GET("/books/:id", bookById)
	router.POST("/books", createBook)
	router.PUT("/checkout", checkoutBook)
	router.PUT("/return", returnBook)
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	router.Run("localhost:8080")
}

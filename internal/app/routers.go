package app

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

// Route is the information for every URI.
type Route struct {
	// Name is the name of this Route.
	Name string
	// Method is the string for the HTTP method. ex) GET, POST etc..
	Method string
	// Pattern is the pattern of the URI.
	Pattern string
	// HandlerFunc is the handler function of this route.
	HandlerFunc gin.HandlerFunc
}

// Routes is the list of the generated Route.
type Routes []Route

// NewRouter returns a new router.
func NewRouter() *gin.Engine {
	router := gin.Default()

	for _, route := range routes {
		switch route.Method {
		case http.MethodGet:
			router.GET(route.Pattern, route.HandlerFunc)
		case http.MethodPost:
			router.POST(route.Pattern, route.HandlerFunc)
		case http.MethodPut:
			router.PUT(route.Pattern, route.HandlerFunc)
		case http.MethodDelete:
			router.DELETE(route.Pattern, route.HandlerFunc)
		}
	}

	return router
}

// Index is the index handler.
func Index(c *gin.Context) {
	c.String(http.StatusOK, "OK")
}

func Status(c *gin.Context) {
	c.String(http.StatusOK, "OK")
}

var routes = Routes{
	{
		"Index",
		http.MethodGet,
		"/library-service/v1/",
		Index,
	},

	{
		"Status",
		http.MethodGet,
		"/library-service/v1/status",
		Status,
	},

	{
		"CreateBook",
		http.MethodPost,
		"/library-service/v1/books",
		CreateBook,
	},

	{
		"DeleteBookById",
		http.MethodDelete,
		"/library-service/v1/books/:id",
		DeleteBookById,
	},

	{
		"GetBookById",
		http.MethodGet,
		"/library-service/v1/books/:id",
		GetBookById,
	},

	{
		"ListBooks",
		http.MethodGet,
		"/library-service/v1/books",
		ListBooks,
	},

	{
		"UpdateBookById",
		http.MethodPut,
		"/library-service/v1/books/:id",
		UpdateBookById,
	},
}

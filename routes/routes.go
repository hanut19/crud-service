package routes

import (
	"fmt"
	"log"
	"net/http"

	routeHandler "curd-service/handlers"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
)

var router *mux.Router

//----------------------ROUTES-------------------------------
//create a mux router
func CreateRouter() {
	router = mux.NewRouter()
}

//initialize all routes
func InitializeRoute() {
	router.HandleFunc("/", routeHandler.Index).Methods("GET")
	router.HandleFunc("/product", routeHandler.IsAuthorized(routeHandler.CreateProductHandler)).Methods("POST")
	//router.HandleFunc("/products", routeHandler.IsAuthorized(routeHandler.GetProductsHandler)).Methods("GET")
	router.HandleFunc("/products", routeHandler.IsAuthorized(routeHandler.GetProductsHandler)).Methods("GET")
	router.HandleFunc("/product/{id}", routeHandler.IsAuthorized(routeHandler.GetProductHandler)).Methods("GET")
	router.HandleFunc("/download/{id}", routeHandler.DownloadHandler).Methods("GET")
	router.HandleFunc("/product/{id}", routeHandler.IsAuthorized(routeHandler.DeleteProductHandler)).Methods("DELETE")

	router.HandleFunc("/product/{id}", routeHandler.IsAuthorized(routeHandler.UpdateProductHandler)).Methods("PUT")
	router.HandleFunc("/search/{search}", routeHandler.SearchHandler).Methods("GET")
	router.Methods("OPTIONS").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//w.Header().Set("Content-Type", "application/json")
		w.Header().Set("Access-Control-Allow-Origin", "")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Access-Control-Request-Headers, Access-Control-Request-Method, Connection, Host, Origin, User-Agent, Referer, Cache-Control, X-header")
	})
}

//start the server
func ServerStart() {
	fmt.Println("Server started at http://localhost:8081")
	err := http.ListenAndServe(":8081", handlers.CORS(handlers.AllowedHeaders([]string{"X-Requested-With", "Access-Control-Allow-Origin", "Content-Type", "Authorization"}), handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "HEAD", "OPTIONS"}), handlers.AllowedOrigins([]string{"*"}))(router))
	if err != nil {
		log.Fatal(err)
	}
}

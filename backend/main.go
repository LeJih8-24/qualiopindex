package main

import (
	"log"
	"net/http"
	"os"
	"path/filepath"

	"qualiopi-backend/database"
	"qualiopi-backend/routes"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {
	database.Connect()

	r := mux.NewRouter()
	routes.RegisterRoutes(r)

		spa := http.FileServer(http.Dir("./public"))
    	r.PathPrefix("/").HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
    		path := filepath.Join("./public", filepath.Clean(req.URL.Path))
    		if info, err := os.Stat(path); err != nil || info.IsDir() {
    			http.ServeFile(w, req, "./public/index.html")
    			return
    		}
    		spa.ServeHTTP(w, req)
    	})

	c := cors.New(cors.Options{
		AllowedOrigins:   []string{"http://localhost:4200"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
	})

//     log.Fatal(http.ListenAndServe(":8080", nil))

	handler := c.Handler(r)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server starting on port %s", port)
	log.Fatal(http.ListenAndServe(":"+port, handler))
}

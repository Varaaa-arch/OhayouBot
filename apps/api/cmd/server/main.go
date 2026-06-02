package main 

import ( 
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"
)

func main() {
	r := chi.NewRouter()

	r.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OhayouBot API Running"))
	})

	fmt.Println("Server is running on port :8080")
	http.ListenAndServe(":8080", r)
}
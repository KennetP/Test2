package main

import (
	"log"
	"mime"
	"net/http"
)

func middlewareA(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Executing middleware A")
		next.ServeHTTP(w, r)
		log.Println("Executing middleware A again (after the handler is finished)")
	})
}

func middlewareB(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Executing middleware B")
		if r.URL.Path == "/cherry" {
			return
		}
		next.ServeHTTP(w, r)
		log.Println("Executing middleware B again")
	})
}

func middlewareC(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Executing middleware C (on the way down)")
		next.ServeHTTP(w, r)
		log.Println("Executing middleware C again (on the way up)")
	})
}

func middlewareD(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Executing middleware D (before authentication)")
		// perform authentication here
		next.ServeHTTP(w, r)
	})
}

func middlewareE(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		log.Println("Executing middleware E (before authentication)")
		// perform authentication here
		next.ServeHTTP(w, r)
	})
}

func enforceJSONHandler(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		contentType := r.Header.Get("Content-Type")
		if contentType != "" {
			mt, _, err := mime.ParseMediaType(contentType)
			if err != nil {
				http.Error(w, "Malformed Content-Type header", http.StatusBadRequest)
				return
			}
			if mt != "application/json" {
				http.Error(w, "Content-Type header must be application/json", http.StatusUnsupportedMediaType)
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}

func final(w http.ResponseWriter, r *http.Request) {
	log.Println("Executing final handler")
	w.Write([]byte("OK"))
}

func main() {
	mux := http.NewServeMux()
	finalHandler := http.HandlerFunc(final)

	// Example 1 - Chaining middleware functions
	mux.Handle("/example1", middlewareA(middlewareB(middlewareC(finalHandler))))

	// Example 2 - Using middleware to perform authentication
	mux.Handle("/example2", middlewareD(middlewareE(finalHandler)))

	// Example 3 - Using middleware to enforce content type
	mux.Handle("/example3", enforceJSONHandler(finalHandler))

	log.Print("Starting server on :4000")
	err := http.ListenAndServe(":4000", mux)
	log.Fatal(err)
}

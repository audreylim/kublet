package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"os"
	"time"
)

func main() {
	mux := http.NewServeMux()
	mux.Handle("/public/", logging(public()))
	mux.Handle("/", logging(index()))
	mux.Handle("/contact", logging(contact()))
	mux.Handle("/privacy-policy", logging(privacy()))
	http.Handle("/public/", http.StripPrefix("/public/", http.FileServer(http.Dir("./public"))))

	port, ok := os.LookupEnv("PORT")
	if !ok {
		port = "8080"
	}

	addr := fmt.Sprintf(":%s", port)
	server := http.Server{
		Addr:         addr,
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  15 * time.Second,
	}
	log.Println("main: running simple server on port", port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("main: couldn't start simple server: %v\n", err)
	}
}

// logging is middleware for wrapping any handler we want to track response
// times for and to see what resources are requested.
func logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		req := fmt.Sprintf("%s %s", r.Method, r.URL)
		log.Println(req)
		next.ServeHTTP(w, r)
		log.Println(req, "completed in", time.Now().Sub(start))
	})
}

//var public = template.Must(template.ParseFiles("./public/index.html", "./public/body.html"))

// index is the handler responsible for rending the index page for the site.
func index() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tmpl := template.Must(template.ParseFiles("./layouts/index.html", "./layouts/home.html", "./layouts/forms.html"))
		err := tmpl.ExecuteTemplate(w, "index", nil)
		if err != nil {
			http.Error(w, fmt.Sprintf("index: couldn't parse template: %v", err), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	})
}

type ContactDetails struct {
	Name    string
	Email   string
	Subject string
	Message string
}

// contact is the handler responsible for rending the index page for the site.
func contact() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			return
		}

		//details := ContactDetails{
		//	Name:    r.FormValue("name"),
		//	Email:   r.FormValue("email"),
		//	Subject: r.FormValue("subject"),
		//	Message: r.FormValue("message"),
		//}

		fmt.Println("Email Sent Successfully!")

		return
	})
}

// privacy is the handler responsible for rending the index page for the site.
func privacy() http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var public = template.Must(template.ParseFiles("./layouts/index.html", "./layouts/privacy.html"))
		err := public.ExecuteTemplate(w, "privacy", nil)
		if err != nil {
			http.Error(w, fmt.Sprintf("index: couldn't parse template: %v", err), http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	})
}

// public serves static assets such as CSS and JavaScript to clients.
func public() http.Handler {
	return http.StripPrefix("/public/", http.FileServer(http.Dir("./public")))
}

package server

import (
	"log"
	"net/http"
	"os"

	"gopkg.in/yaml.v2"
)

func Start() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		out, err := yaml.Marshal(map[string]string{
			"foo": "bar",
		})
		if err != nil {
			http.Error(w, err.Error(), 500)
			return
		}

		log.Printf("%v", out)
		w.Write([]byte("Hello, world"))
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}
	http.ListenAndServe(":"+port, nil)
}

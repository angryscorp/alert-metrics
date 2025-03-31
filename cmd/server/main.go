package main

import (
	"github.com/angryscorp/alert-metrics/internal/infrastructure/metricrouter"
	"github.com/angryscorp/alert-metrics/internal/infrastructure/metricstorage"
	"net/http"
)

func main() {
	router := metricrouter.NewRouter(
		http.NewServeMux(),
		metricstorage.NewMemStorage(),
	)

	err := http.ListenAndServe(":8080", router)
	if err != nil {
		panic(err)
	}
}

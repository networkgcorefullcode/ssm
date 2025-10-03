package server

import (
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger"
)

func ServerSwagger() {
	// Expone la UI de Swagger usando http-swagger
	http.HandleFunc("/swagger-ui/", httpSwagger.WrapHandler)

	// Inicia el servidor en el puerto 8080
	http.ListenAndServe(":8080", nil)
}

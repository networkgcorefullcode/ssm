package server

import (
	"net/http"

	httpSwagger "github.com/swaggo/http-swagger"
)

func ServerSwagger() {
	// Expone la UI de Swagger usando http-swagger
	http.HandleFunc("/swagger-ui/", httpSwagger.WrapHandler)

	// Inicia el servidor en el puerto 9001
	http.ListenAndServe(":9001", nil)
}

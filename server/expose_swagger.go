package server

import (
	"net/http"

	constants "github.com/networkgcorefullcode/ssm/const"
	"github.com/swaggo/http-swagger"
)

func ServerSwagger() {
	// Servir el archivo swagger.yaml desde la carpeta ./docs
	http.Handle("/swagger/", http.StripPrefix("/swagger/", http.FileServer(http.Dir(constants.APP_PATH_SWAGGER))))

	// Redirigir a la interfaz Swagger UI usando http-swagger
	http.HandleFunc("/swagger-ui/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/swagger/index.html?url=/swagger/swagger.yaml", http.StatusFound)
	})

	// Iniciar el servidor en el puerto 8080
	http.ListenAndServe(":8080", nil)
}

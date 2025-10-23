package server

import (
	"net/http"

	"github.com/networkgcorefullcode/ssm/handlers"
	"github.com/networkgcorefullcode/ssm/logger"
)

func CreateEndpointHandlersPool() {

	// Set up HTTP handlers
	// Encrypt endpoints POST
	http.HandleFunc("/encrypt", func(w http.ResponseWriter, r *http.Request) {
		logger.AppLog.Debugf("Received /encrypt request")
		handlers.EncryptWithPool(w, r)
	})

	// Decrypt endpoints POST
	http.HandleFunc("/decrypt", func(w http.ResponseWriter, r *http.Request) {
		logger.AppLog.Debugf("Received /decrypt request")
		handlers.DecryptWithPool(w, r)
	})

	// Store Key endpoints POST
	http.HandleFunc("/store-key", func(w http.ResponseWriter, r *http.Request) {
		logger.AppLog.Debugf("Received /store-key request")
		handlers.StoreKeyWithPool(w, r)
	})

	// Generate Key endpoints POST
	http.HandleFunc("/generate-aes-key", func(w http.ResponseWriter, r *http.Request) {
		logger.AppLog.Debugf("Received /generate-aes-key request")
		handlers.GenerateAESKeyWithPool(w, r)
	})

	http.HandleFunc("/generate-des3-key", func(w http.ResponseWriter, r *http.Request) {
		logger.AppLog.Debugf("Received /generate-des3-key request")
		handlers.GenerateDES3KeyWithPool(w, r)
	})

	http.HandleFunc("/generate-des-key", func(w http.ResponseWriter, r *http.Request) {
		logger.AppLog.Debugf("Received /generate-des-key request")
		handlers.GenerateDESKeyWithPool(w, r)
	})

	// Syncronization handlers
	http.HandleFunc("/get-data-keys", func(w http.ResponseWriter, r *http.Request) {
		logger.AppLog.Debugf("Received /get-data-keys request")
		handlers.GetDataKeysWithPool(w, r)
	})

	http.HandleFunc("/get-key", func(w http.ResponseWriter, r *http.Request) {
		logger.AppLog.Debugf("Received /get-keys request")
		handlers.GetDataKeyWithPool(w, r)
	})

	http.HandleFunc("/get-all-keys", func(w http.ResponseWriter, r *http.Request) {
		logger.AppLog.Debugf("Received /get-all-keys request")
		handlers.GetAllKeysWithPool(w, r)
	})
}

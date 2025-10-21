package server

import (
	"net/http"

	"github.com/networkgcorefullcode/ssm/handlers"
	"github.com/networkgcorefullcode/ssm/logger"
)

func CreateEndpointHandlers(s *SSM) {
	// Set up HTTP handlers
	// Encrypt endpoints POST
	http.HandleFunc("/encrypt", func(w http.ResponseWriter, r *http.Request) {
		logger.AppLog.Debugf("Received /encrypt request")
		handlers.HandleEncrypt(s.mgr, w, r)
	})

	// Decrypt endpoints POST
	http.HandleFunc("/decrypt", func(w http.ResponseWriter, r *http.Request) {
		logger.AppLog.Debugf("Received /decrypt request")
		handlers.HandleDecrypt(s.mgr, w, r)
	})

	// Store Key endpoints POST
	http.HandleFunc("/store-key", func(w http.ResponseWriter, r *http.Request) {
		logger.AppLog.Debugf("Received /store-key request")
		handlers.HandleStoreKey(s.mgr, w, r)
	})

	// Generate Key endpoints POST
	http.HandleFunc("/generate-aes-key", func(w http.ResponseWriter, r *http.Request) {
		logger.AppLog.Debugf("Received /generate-aes-key request")
		handlers.HandleGenerateAESKey(s.mgr, w, r)
	})

	http.HandleFunc("/generate-des3-key", func(w http.ResponseWriter, r *http.Request) {
		logger.AppLog.Debugf("Received /generate-des3-key request")
		handlers.HandleGenerateDES3Key(s.mgr, w, r)
	})

	http.HandleFunc("/generate-des-key", func(w http.ResponseWriter, r *http.Request) {
		logger.AppLog.Debugf("Received /generate-des-key request")
		handlers.HandleGenerateDESKey(s.mgr, w, r)
	})
}

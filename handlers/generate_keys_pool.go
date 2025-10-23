package handlers

import (
	"encoding/json"
	"errors"
	"net/http"

	constants "github.com/networkgcorefullcode/ssm/const"
	"github.com/networkgcorefullcode/ssm/logger"
	"github.com/networkgcorefullcode/ssm/models"
	"github.com/networkgcorefullcode/ssm/pkcs11mgr"
)

// GenerateAESKeyWithPool handles AES key generation using connection pool
func GenerateAESKeyWithPool(w http.ResponseWriter, r *http.Request) error {
	logger.AppLog.Info("Processing AES key generation request with pool")

	var req models.GenAESKeyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.AppLog.Errorf("Failed to decode request body: %v", err)
		sendProblemDetails(w, "Bad Request", "El cuerpo de la petición no es válido JSON", "INVALID_JSON", http.StatusBadRequest, r.URL.Path)
		return err
	}

	if req.Id <= 0 {
		logger.AppLog.Error("ID is required but was empty")
		sendProblemDetails(w, "Bad Request", "El campo 'id' es requerido y no puede estar vacío", "MISSING_ID", http.StatusBadRequest, r.URL.Path)
		return errors.New("ID is required")
	}

	if req.Bits != 128 && req.Bits != 256 {
		logger.AppLog.Errorf("Invalid key size: %d bits", req.Bits)
		sendProblemDetails(w, "Bad Request", "El tamaño de clave debe ser 128 o 256 bits", "INVALID_KEY_SIZE", http.StatusBadRequest, r.URL.Path)
		return errors.New("invalid key size")
	}

	logger.AppLog.Infof("Generating AES key - ID: %s, Bits: %d", req.Id, req.Bits)

	var resp models.GenAESKeyResponse

	// Use connection pool
	err := pkcs11mgr.WithConnection(func(mgr *pkcs11mgr.Manager) error {
		var label string
		if req.Bits == 128 {
			label = constants.LABEL_ENCRYPTION_KEY_AES128
		} else if req.Bits == 256 {
			label = constants.LABEL_ENCRYPTION_KEY_AES256
		}

		handle, err := mgr.GenerateAESKey(label, req.Id, int(req.Bits))
		if err != nil {
			logger.AppLog.Errorf("AES key generation failed: %v", err)
			sendProblemDetails(w, "Key Generation Failed", "Error al generar la clave AES en el HSM", "KEY_GENERATION_ERROR", http.StatusInternalServerError, r.URL.Path)
			return err
		}

		logger.AppLog.Infof("AES key generated successfully - Handle: %d", handle)

		resp = models.GenAESKeyResponse{
			Handle: int32(handle),
			Id:     req.Id,
			Bits:   req.Bits,
		}

		return nil
	})

	if err != nil {
		// Error was already handled inside the WithConnection function
		return err
	}

	// Send successful response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		logger.AppLog.Errorf("Failed to encode response: %v", err)
		return err
	}

	return nil
}

// GenerateDESKeyWithPool handles DES key generation using connection pool
func GenerateDESKeyWithPool(w http.ResponseWriter, r *http.Request) error {
	logger.AppLog.Info("Processing DES key generation request with pool")

	var req models.GenDESKeyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.AppLog.Errorf("Failed to decode request body: %v", err)
		sendProblemDetails(w, "Bad Request", "El cuerpo de la petición no es válido JSON", "INVALID_JSON", http.StatusBadRequest, r.URL.Path)
		return err
	}

	if req.Id <= 0 {
		logger.AppLog.Error("ID is required but was empty")
		sendProblemDetails(w, "Bad Request", "El campo 'id' es requerido y no puede estar vacío", "MISSING_ID", http.StatusBadRequest, r.URL.Path)
		return errors.New("ID is required")
	}

	logger.AppLog.Infof("Generating DES key - ID: %d", req.Id)

	var resp models.GenDESKeyResponse

	// Use connection pool
	err := pkcs11mgr.WithConnection(func(mgr *pkcs11mgr.Manager) error {
		handle, err := mgr.GenerateDESKey(constants.LABEL_ENCRYPTION_KEY_DES, req.Id)
		if err != nil {
			logger.AppLog.Errorf("DES key generation failed: %v", err)
			sendProblemDetails(w, "Key Generation Failed", "Error al generar la clave DES en el HSM", "KEY_GENERATION_ERROR", http.StatusInternalServerError, r.URL.Path)
			return err
		}

		logger.AppLog.Infof("DES key generated successfully - Handle: %d", handle)

		resp = models.GenDESKeyResponse{
			Handle: int32(handle),
			Id:     req.Id,
		}

		return nil
	})

	if err != nil {
		// Error was already handled inside the WithConnection function
		return err
	}

	// Send successful response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		logger.AppLog.Errorf("Failed to encode response: %v", err)
		return err
	}

	return nil
}

// GenerateDES3KeyWithPool handles DES3 key generation using connection pool
func GenerateDES3KeyWithPool(w http.ResponseWriter, r *http.Request) error {
	logger.AppLog.Info("Processing DES3 key generation request with pool")

	var req models.GenDES3KeyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.AppLog.Errorf("Failed to decode request body: %v", err)
		sendProblemDetails(w, "Bad Request", "El cuerpo de la petición no es válido JSON", "INVALID_JSON", http.StatusBadRequest, r.URL.Path)
		return err
	}

	if req.Id <= 0 {
		logger.AppLog.Error("ID is required but was empty")
		sendProblemDetails(w, "Bad Request", "El campo 'id' es requerido y no puede estar vacío", "MISSING_ID", http.StatusBadRequest, r.URL.Path)
		return errors.New("ID is required")
	}

	logger.AppLog.Infof("Generating DES3 key - ID: %d", req.Id)

	var resp models.GenDES3KeyResponse

	// Use connection pool
	err := pkcs11mgr.WithConnection(func(mgr *pkcs11mgr.Manager) error {
		handle, err := mgr.GenerateDES3Key(constants.LABEL_ENCRYPTION_KEY_DES3, req.Id)
		if err != nil {
			logger.AppLog.Errorf("DES3 key generation failed: %v", err)
			sendProblemDetails(w, "Key Generation Failed", "Error al generar la clave DES3 en el HSM", "KEY_GENERATION_ERROR", http.StatusInternalServerError, r.URL.Path)
			return err
		}

		logger.AppLog.Infof("DES3 key generated successfully - Handle: %d", handle)

		resp = models.GenDES3KeyResponse{
			Handle: int32(handle),
			Id:     req.Id,
		}

		return nil
	})

	if err != nil {
		// Error was already handled inside the WithConnection function
		return err
	}

	// Send successful response
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		logger.AppLog.Errorf("Failed to encode response: %v", err)
		return err
	}

	return nil
}

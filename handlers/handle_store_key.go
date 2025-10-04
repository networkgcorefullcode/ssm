package handlers

import (
	"encoding/base64"
	"encoding/json"
	"net/http"

	"github.com/miekg/pkcs11"
	constants "github.com/networkgcorefullcode/ssm/const"
	"github.com/networkgcorefullcode/ssm/logger"
	"github.com/networkgcorefullcode/ssm/models"
	"github.com/networkgcorefullcode/ssm/pkcs11mgr"
)

// HandleStoreKey maneja las peticiones de almacenamiento de claves
// @Summary Almacenar clave
// @Description Almacena una clave en el HSM y opcionalmente la encripta
// @Tags Key Management
// @Accept json
// @Produce json
// @Param request body models.StoreKeyRequest true "Datos de la clave a almacenar"
// @Success 200 {object} models.StoreKeyResponse "Clave almacenada exitosamente"
// @Failure 400 {object} models.ProblemDetails "Petición inválida"
// @Failure 500 {object} models.ProblemDetails "Error interno del servidor"
// @Router /store-key [post]
func HandleStoreKey(mgr *pkcs11mgr.Manager, w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		postStoreKey(mgr, w, r)
	default:
		sendProblemDetails(w, "Method Not Allowed", "El método HTTP no está permitido para este endpoint", "METHOD_NOT_ALLOWED", http.StatusMethodNotAllowed, r.URL.Path)
	}
}

func postStoreKey(mgr *pkcs11mgr.Manager, w http.ResponseWriter, r *http.Request) {
	logger.AppLog.Info("Processing store key request")

	var req models.StoreKeyRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.AppLog.Errorf("Failed to decode request body: %v", err)
		sendProblemDetails(w, "Bad Request", "El cuerpo de la petición no es válido JSON", "INVALID_JSON", http.StatusBadRequest, r.URL.Path)
		return
	}

	label := req.KeyLabel
	id := req.Id
	logger.AppLog.Infof("Decoding key value for label: %s, ID: %s", label, id)
	key_value, err := base64.StdEncoding.DecodeString(req.KeyValue)
	if err != nil {
		logger.AppLog.Errorf("Failed to decode base64 key value: %v", err)
		sendProblemDetails(w, "Bad Request", "El valor de la clave en base64 no es válido", "INVALID_BASE64", http.StatusBadRequest, r.URL.Path)
		return
	}

	logger.AppLog.Infof("Storing key in HSM - Label: %s", label)
	handle, err := mgr.StoreKey(label, key_value, []byte(id))
	if err != nil {
		logger.AppLog.Errorf("Failed to store key: %v", err)
		sendProblemDetails(w, "Key Storage Failed", "Error al almacenar la clave en el HSM", "KEY_STORAGE_ERROR", http.StatusInternalServerError, r.URL.Path)
		return
	}

	logger.AppLog.Infof("Key stored successfully - Handle: %d", handle)

	resp := models.StoreKeyResponse{
		Handle:    uint(handle),
		CipherKey: nil, // Inicialmente nil, se asignará si se puede encriptar
	}

	// Intentar encontrar la clave de encriptación para encriptar el valor almacenado
	logger.AppLog.Infof("Looking for encryption key: %s", constants.LABEL_ENCRYPTION_KEY)
	findHandle, err := mgr.FindKeyByLabel(constants.LABEL_ENCRYPTION_KEY)
	if err != nil || findHandle == 0 {
		logger.AppLog.Warnf("Encryption key not found or error: %v. Returning response without encrypted key", err)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		if encodeErr := json.NewEncoder(w).Encode(resp); encodeErr != nil {
			logger.AppLog.Errorf("Failed to encode response: %v", encodeErr)
		}
		return
	}

	// Encriptar el valor de la clave almacenada
	logger.AppLog.Info("Encrypting stored key value")
	cipher, err := mgr.EncryptKey(findHandle, nil, key_value, pkcs11.CKM_AES_CBC_PAD)
	if err != nil {
		logger.AppLog.Errorf("Failed to encrypt key value: %v. Returning response without encrypted key", err)
		resp.CipherKey = nil
	} else {
		logger.AppLog.Info("Key value encrypted successfully")
		encryptedKeyB64 := base64.StdEncoding.EncodeToString(cipher)
		resp.CipherKey = &encryptedKeyB64
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if encodeErr := json.NewEncoder(w).Encode(resp); encodeErr != nil {
		logger.AppLog.Errorf("Failed to encode response: %v", encodeErr)
	}

}

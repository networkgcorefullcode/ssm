package handlers

import (
	"encoding/hex"
	"encoding/json"
	"net/http"
	"time"

	"github.com/miekg/pkcs11"
	constants "github.com/networkgcorefullcode/ssm/const"
	"github.com/networkgcorefullcode/ssm/factory"
	"github.com/networkgcorefullcode/ssm/logger"
	"github.com/networkgcorefullcode/ssm/models"
	"github.com/networkgcorefullcode/ssm/pkcs11mgr"
	"github.com/networkgcorefullcode/ssm/safe"
)

// HandleEncrypt handles encryption requests
// @Summary Encrypt data
// @Description Encrypts data using an AES key stored in the HSM
// @Tags Encryption
// @Accept json
// @Produce json
// @Param request body models.EncryptRequest true "Data to encrypt"
// @Success 200 {object} models.EncryptResponse "Data encrypted successfully"
// @Failure 400 {object} models.ProblemDetails "Invalid request"
// @Failure 404 {object} models.ProblemDetails "Key not found"
// @Failure 500 {object} models.ProblemDetails "Internal server error"
// @Router /encrypt [post]
func HandleEncrypt(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodPost:
		postEncrypt(w, r)
	default:
		sendProblemDetails(w, "Method Not Allowed", "The HTTP method is not allowed for this endpoint", "METHOD_NOT_ALLOWED", http.StatusMethodNotAllowed, r.URL.Path)
	}
}

func postEncrypt(w http.ResponseWriter, r *http.Request) {

	// init the pkcs manager
	mgr, err := pkcs11mgr.New(factory.SsmConfig.Configuration.PkcsPath,
		uint(factory.SsmConfig.Configuration.LotsNumber),
		factory.SsmConfig.Configuration.Pin)
	if err != nil {
		logger.AppLog.Errorf("Failed to create PKCS11 manager: %v", err)
		sendProblemDetails(w, "Internal Server Error", "Failed to initialize PKCS11 manager", "PKCS_INIT_ERROR", http.StatusInternalServerError, r.URL.Path)
		return
	}

	err = mgr.OpenSession()
	if err != nil {
		logger.AppLog.Errorf("Failed to OpenSession PKCS11 manager: %v", err)
		sendProblemDetails(w, "Internal Server Error", "The pkcs session have a error during stablishment", "PKCS_ERROR", http.StatusInternalServerError, r.URL.Path)
		return
	}

	defer mgr.CloseSession()
	defer mgr.Finalize()

	logger.AppLog.Info("Processing encrypt request")

	var req models.EncryptRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		logger.AppLog.Errorf("Failed to decode request body: %v", err)
		sendProblemDetails(w, "Bad Request", "The request body is not valid JSON", "INVALID_JSON", http.StatusBadRequest, r.URL.Path)
		return
	}

	logger.AppLog.Infof("Decoding hex plaintext for key: %s", req.KeyLabel)
	pt, err := hex.DecodeString(req.Plain)
	if err != nil {
		logger.AppLog.Errorf("Failed to decode hex plaintext: %v", err)
		sendProblemDetails(w, "Bad Request", "The hex data is not valid", "INVALID_HEX", http.StatusBadRequest, r.URL.Path)
		return
	}

	logger.AppLog.Infof("Finding key by label: %s", req.KeyLabel)
	keyHandle, err := mgr.FindKeyLabelReturnRandom(req.KeyLabel)
	if err != nil {
		logger.AppLog.Errorf("Key not found: %s, error: %v", req.KeyLabel, err)
		sendProblemDetails(w, "Key Not Found", "The specified key does not exist in the HSM", "KEY_NOT_FOUND", http.StatusNotFound, r.URL.Path)
		return
	}
	atrr, err := mgr.GetObjectAttributes(keyHandle)
	if err != nil {
		logger.AppLog.Errorf("Atributes not found: %s, error: %v", req.KeyLabel, err)
		sendProblemDetails(w, "Atributes Not Found", "The specified key does not exist in the HSM", "KEY_NOT_FOUND", http.StatusNotFound, r.URL.Path)
		return
	}

	logger.AppLog.Info("Generating initialization vector (IV)")
	var size int
	if req.EncryptionAlgorithm == constants.ALGORITHM_DES3_OurUsers || req.EncryptionAlgorithm == constants.ALGORITHM_DES_OurUsers {
		size = 8
	} else if req.EncryptionAlgorithm == constants.ALGORITHM_AES128_OurUsers || req.EncryptionAlgorithm == constants.ALGORITHM_AES256_OurUsers {
		size = 16
	}
	iv := make([]byte, size)
	if err := safe.RandRead(iv); err != nil {
		logger.AppLog.Errorf("Failed to generate IV: %v", err)
		sendProblemDetails(w, "Internal Server Error", "Error generating initialization vector", "IV_GENERATION_FAILED", http.StatusInternalServerError, r.URL.Path)
		return
	}

	logger.AppLog.Info("Encrypting data")

	var ciphertext []byte
	switch req.EncryptionAlgorithm {
	case constants.ALGORITHM_AES128_OurUsers, constants.ALGORITHM_AES256_OurUsers:
		ciphertext, err = mgr.EncryptKey(keyHandle, iv, pt, pkcs11.CKM_AES_CBC_PAD)
		if err != nil {
			logger.AppLog.Errorf("Encryption failed: %v", err)
			ciphertext, err = mgr.EncryptKey(keyHandle, iv, pt, pkcs11.CKM_AES_CBC)
			if err != nil {
				logger.AppLog.Errorf("Encryption failed: %v", err)
				sendProblemDetails(w, "Encryption Failed", "Error during encryption process", "ENCRYPTION_ERROR", http.StatusInternalServerError, r.URL.Path)
				return
			}
		}
	case constants.ALGORITHM_DES3_OurUsers:
		ciphertext, err = mgr.EncryptKey(keyHandle, iv, pt, pkcs11.CKM_DES3_CBC_PAD)
		if err != nil {
			logger.AppLog.Errorf("Encryption failed: %v", err)
			ciphertext, err = mgr.EncryptKey(keyHandle, iv, pt, pkcs11.CKM_DES3_CBC)
			if err != nil {
				logger.AppLog.Errorf("Encryption failed: %v", err)
				sendProblemDetails(w, "Encryption Failed", "Error during encryption process", "ENCRYPTION_ERROR", http.StatusInternalServerError, r.URL.Path)
				return
			}
		}
	case constants.ALGORITHM_DES_OurUsers:
		ciphertext, err = mgr.EncryptKey(keyHandle, iv, pt, pkcs11.CKM_DES_CBC_PAD)
		if err != nil {
			logger.AppLog.Errorf("Encryption failed: %v", err)
			ciphertext, err = mgr.EncryptKey(keyHandle, iv, pt, pkcs11.CKM_DES_CBC)
			if err != nil {
				logger.AppLog.Errorf("Encryption failed: %v", err)
				sendProblemDetails(w, "Encryption Failed", "Error during encryption process", "ENCRYPTION_ERROR", http.StatusInternalServerError, r.URL.Path)
				return
			}
		}
	}

	safe.Zero(pt) // Clear sensitive data from memory

	logger.AppLog.Info("Encryption completed successfully")

	ciphertextStr := hex.EncodeToString(ciphertext)
	ivStr := hex.EncodeToString(iv)
	ok := true

	timeCreated := time.Now()
	timeUpdated := timeCreated

	// Create response using the structured model
	resp := models.EncryptResponse{
		Cipher:      ciphertextStr,
		Iv:          ivStr,
		Ok:          ok,
		TimeCreated: timeCreated,
		TimeUpdated: timeUpdated,
		Id:          atrr.Id,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	if err := json.NewEncoder(w).Encode(resp); err != nil {
		logger.AppLog.Errorf("Failed to encode response: %v", err)
	}
}

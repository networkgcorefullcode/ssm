package k4opt

import (
	"encoding/json"
	"net/http"

	"github.com/networkgcorefullcode/ssm/models"
	"github.com/networkgcorefullcode/ssm/pkcs11mgr"
)

func HandleGenerateAESKey(mgr *pkcs11mgr.Manager, w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	// Validar método HTTP
	if r.Method != http.MethodPost {
		sendErrorResponse(w, "Method Not Allowed", "Only POST method is allowed", http.StatusMethodNotAllowed)
		return
	}

	// Leer y parsear el cuerpo JSON
	var req models.GenAESKeyRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		sendErrorResponse(w, "Invalid JSON", "Failed to parse request body: "+err.Error(), http.StatusBadRequest)
		return
	}

	// Validar campos requeridos
	if req.Label == "" {
		sendErrorResponse(w, "Validation Error", "Label is required", http.StatusBadRequest)
		return
	}

	if req.ID == "" {
		sendErrorResponse(w, "Validation Error", "ID is required", http.StatusBadRequest)
		return
	}

	// Validar bits (solo 128, 192, 256)
	if req.Bits != 128 && req.Bits != 192 && req.Bits != 256 {
		sendErrorResponse(w, "Validation Error", "Bits must be 128, 192, or 256", http.StatusBadRequest)
		return
	}

	// Generar la clave AES
	handle, err := mgr.GenerateAESKey(req.Label, []byte(req.ID), req.Bits)
	if err != nil {
		sendErrorResponse(w, "Key Generation Failed", "Failed to generate AES key: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Respuesta exitosa
	response := models.GenAESKeyResponse{
		Handle: uint(handle),
		Label:  req.Label,
		ID:     req.ID,
		Bits:   req.Bits,
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(response)
}

func sendErrorResponse(w http.ResponseWriter, title, detail string, status int) {
	problemDetails := models.ProblemDetails{
		Title:  title,
		Detail: detail,
		Status: status,
	}

	w.WriteHeader(status)
	json.NewEncoder(w).Encode(problemDetails)
}

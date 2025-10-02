package server

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/networkgcorefullcode/ssm/pkcs11mgr"
	"github.com/networkgcorefullcode/ssm/safe"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type SSM struct {
	mgr         *pkcs11mgr.Manager
	mongoclient *mongo.Client
}

func New(mgr *pkcs11mgr.Manager, mongoURI string) (*SSM, error) {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(mongoURI))
	if err != nil {
		return nil, err
	}
	return &SSM{mgr: mgr, mongoclient: client}, nil
}

func (s *SSM) Start(socketPath string) error {
	// remove old socket
	_ = os.Remove(socketPath)
	l, err := net.Listen("unix", socketPath)
	if err != nil {
		return err
	}
	http.HandleFunc("/encrypt", s.handleEncrypt)
	http.HandleFunc("/decrypt", s.handleDecrypt)
	go func() {
		log.Printf("SSM listening on unix socket %s", socketPath)
		if err := http.Serve(l, nil); err != nil {
			log.Fatalf("server error: %v", err)
		}
	}()
	return nil
}

type encryptRequest struct {
	KeyLabel string `json:"key_label"`
	PlainB64 string `json:"plain_b64"`
}

func (s *SSM) handleEncrypt(w http.ResponseWriter, r *http.Request) {
	var req encryptRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}

	// Decodifica plaintext (base64 -> []byte)
	pt, err := base64.StdEncoding.DecodeString(req.PlainB64)
	if err != nil {
		http.Error(w, "bad base64", http.StatusBadRequest)
		return
	}

	// Aquí deberías buscar el handle por etiqueta; ejemplo simplificado:
	// handle := s.mgr.FindKeyByLabel(req.KeyLabel)  // implementar
	// for demo, asumimos que ya tienes handle:
	// iv := random 16 bytes
	iv := make([]byte, 16)
	// TODO: usar crypto/rand.Read(iv)
	// Usa manager para encrypt
	ciphertext, err := s.mgr.EncryptWithAESKey( /*handle*/ 1, iv, pt)
	// Scrub plain text
	safe.Zero(pt)

	if err != nil {
		http.Error(w, "encrypt error", 500)
		return
	}

	// Almacena ciphertext + iv + metadata en Mongo
	// Ejemplo simplificado:
	coll := s.mongoclient.Database("aether").Collection("k4")
	doc := map[string]interface{}{
		"key_label": req.KeyLabel,
		"iv":        base64.StdEncoding.EncodeToString(iv),
		"cipher":    base64.StdEncoding.EncodeToString(ciphertext),
		"alg":       "AES-CBC-PAD", // cambiar cuando uses GCM
	}
	_, _ = coll.InsertOne(context.Background(), doc)

	fmt.Fprintf(w, `{"ok": true}`)
}

type decryptRequest struct {
	KeyLabel  string `json:"key_label"`
	CipherB64 string `json:"cipher_b64"`
	IVB64     string `json:"iv_b64"`
}

func (s *SSM) handleDecrypt(w http.ResponseWriter, r *http.Request) {
	var req decryptRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "bad request", http.StatusBadRequest)
		return
	}
	cipher, _ := base64.StdEncoding.DecodeString(req.CipherB64)
	iv, _ := base64.StdEncoding.DecodeString(req.IVB64)

	// Obtener handle por label (implementar)
	plaintext, err := s.mgr.DecryptWithAESKey( /*handle*/ 1, iv, cipher)
	if err != nil {
		http.Error(w, "decrypt error", 500)
		return
	}
	// A estas alturas, el plaintext está en memoria: usar mlock antes de usarlo y zero luego
	// Enviar result (pero idealmente devolver solo resultados, no la clave)
	resp := map[string]string{
		"plain_b64": base64.StdEncoding.EncodeToString(plaintext),
	}
	// scrubbear
	safe.Zero(plaintext)
	_ = json.NewEncoder(w).Encode(resp)
}

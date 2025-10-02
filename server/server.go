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
	"path/filepath"

	"github.com/networkgcorefullcode/ssm/factory"
	"github.com/networkgcorefullcode/ssm/logger"
	"github.com/networkgcorefullcode/ssm/pkcs11mgr"
	"github.com/networkgcorefullcode/ssm/safe"
	"github.com/urfave/cli/v3"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type SSM struct {
	mgr         *pkcs11mgr.Manager
	mongoclient *mongo.Client
}

type (
	// Config information.
	Config struct {
		cfg string
	}
)

func New(mgr *pkcs11mgr.Manager, mongoURI string) (*SSM, error) {
	client, err := mongo.Connect(context.Background(), options.Client().ApplyURI(mongoURI))
	if err != nil {
		return nil, err
	}
	return &SSM{mgr: mgr, mongoclient: client}, nil
}

var config Config

var ssmCLi = []cli.Flag{
	&cli.StringFlag{
		Name:     "cfg",
		Usage:    "ssm config file",
		Required: true,
	},
}

func (ssm *SSM) GetCliCmd() (flags []cli.Flag) {
	return ssmCLi
}

func (ssm *SSM) Initialize(c *cli.Command) error {
	config = Config{
		cfg: c.String("cfg"),
	}

	absPath, err := filepath.Abs(config.cfg)
	if err != nil {
		logger.CfgLog.Errorln(err)
		return err
	}

	if err := factory.InitConfigFactory(absPath); err != nil {
		return err
	}

	ssm.setLogLevel()

	if err := factory.CheckConfigVersion(); err != nil {
		return err
	}

	factory.SsmConfig.CfgLocation = absPath
	return nil
}

func (ausf *SSM) setLogLevel() {
	if factory.SsmConfig.Logger == nil {
		logger.InitLog.Warnln("SSM config without log level setting")
		return
	}

	if factory.SsmConfig.Logger.SSM != nil {
		if factory.SsmConfig.Logger.SSM.DebugLevel != "" {
			if level, err := zapcore.ParseLevel(factory.SsmConfig.Logger.SSM.DebugLevel); err != nil {
				logger.InitLog.Warnf("SSM Log level [%s] is invalid, set to [info] level",
					factory.SsmConfig.Logger.SSM.DebugLevel)
				logger.SetLogLevel(zap.InfoLevel)
			} else {
				logger.InitLog.Infof("SSM Log level is set to [%s] level", level)
				logger.SetLogLevel(level)
			}
		} else {
			logger.InitLog.Warnln("SSM Log level not set. Default set to [info] level")
			logger.SetLogLevel(zap.InfoLevel)
		}
	}
}

func (ausf *SSM) FilterCli(c *cli.Command) (args []string) {
	for _, flag := range ausf.GetCliCmd() {
		name := flag.Names()[0]
		value := fmt.Sprint(c.Generic(name))
		if value == "" {
			continue
		}

		args = append(args, "--"+name, value)
	}
	return args
}

func (s *SSM) Start() error {
	// remove old socket
	socketPath := factory.SsmConfig.Configuration.SocketPath

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

package server

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/networkgcorefullcode/ssm/factory"
	"github.com/networkgcorefullcode/ssm/logger"
)

func startHTTPServer(router *gin.Engine) error {
	// Start HTTPS or HTTP server based on configuration
	if factory.SsmConfig.Configuration.IsHttps == nil || *factory.SsmConfig.Configuration.IsHttps {
		certFile := factory.SsmConfig.Configuration.CertFile
		keyFile := factory.SsmConfig.Configuration.KeyFile
		caFile := factory.SsmConfig.Configuration.CAFile
		bindAddr := factory.SsmConfig.Configuration.BindAddr

		if certFile == "" || keyFile == "" || caFile == "" {
			logger.AppLog.Error("HTTPS is enabled but certFile, keyFile or caFile is not set")
			return fmt.Errorf("missing certFile, keyFile or caFile")
		}

		logger.AppLog.Infof("SSM listening api https (mTLS enforced) on %s", bindAddr)

		// 1️⃣ Load the CA certificate that will verify client certificates
		caCert, err := os.ReadFile(caFile)
		if err != nil {
			return fmt.Errorf("failed to read CA file: %v", err)
		}
		caCertPool := x509.NewCertPool()
		if ok := caCertPool.AppendCertsFromPEM(caCert); !ok {
			return fmt.Errorf("failed to append CA certificate")
		}

		// 2️⃣ Configure TLS to require and verify client certificates
		tlsConfig := &tls.Config{
			ClientCAs:    caCertPool,
			ClientAuth:   tls.RequireAndVerifyClientCert, // 🔒 mTLS obligatory
			MinVersion:   tls.VersionTLS13,               // Force TLS 1.3
			Certificates: make([]tls.Certificate, 1),
		}

		// 3️⃣ Load the server certificate
		serverCert, err := tls.LoadX509KeyPair(certFile, keyFile)
		if err != nil {
			return fmt.Errorf("failed to load server certificate/key: %v", err)
		}
		tlsConfig.Certificates[0] = serverCert

		// 4️⃣ Create the HTTPS server with this configuration
		srv := &http.Server{
			Addr:      bindAddr,
			Handler:   router,
			TLSConfig: tlsConfig,
		}

		// 5️⃣ Start the server with mTLS
		if err := srv.ListenAndServeTLS("", ""); err != nil {
			logger.AppLog.Errorf("Server error: %v", err)
			return err
		}
		return nil

	} else {
		logger.AppLog.Infof("SSM listening api http %s", factory.SsmConfig.Configuration.BindAddr)
		srv := &http.Server{Addr: factory.SsmConfig.Configuration.BindAddr, Handler: router}
		if err := srv.ListenAndServe(); err != nil {
			logger.AppLog.Errorf("Server error: %v", err)
			return err
		}
		return nil
	}
}

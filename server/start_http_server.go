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

		logger.AppLog.Debugf("Starting HTTPS server configuration - CertFile: %s, KeyFile: %s, CAFile: %s", certFile, keyFile, caFile)

		if certFile == "" || keyFile == "" || caFile == "" {
			logger.AppLog.Error("HTTPS is enabled but certFile, keyFile or caFile is not set")
			return fmt.Errorf("missing certFile, keyFile or caFile")
		}

		logger.AppLog.Infof("SSM listening api https (mTLS enforced) on %s", bindAddr)

		// 1️⃣ Load the CA certificate that will verify client certificates
		logger.AppLog.Debugf("Step 1/5: Loading CA certificate from: %s", caFile)
		caCert, err := os.ReadFile(caFile)
		if err != nil {
			logger.AppLog.Errorf("Failed to read CA file %s: %v", caFile, err)
			return fmt.Errorf("failed to read CA file: %v", err)
		}
		logger.AppLog.Debugf("CA certificate file loaded successfully, size: %d bytes", len(caCert))

		caCertPool := x509.NewCertPool()
		if ok := caCertPool.AppendCertsFromPEM(caCert); !ok {
			logger.AppLog.Error("Failed to parse and append CA certificate to pool")
			return fmt.Errorf("failed to append CA certificate")
		}
		logger.AppLog.Debug("CA certificate successfully added to certificate pool for client verification")

		// 2️⃣ Configure TLS to require and verify client certificates
		logger.AppLog.Debug("Step 2/5: Configuring TLS settings - ClientAuth: RequireAndVerifyClientCert, MinVersion: TLS 1.3")
		tlsConfig := &tls.Config{
			ClientCAs:    caCertPool,
			ClientAuth:   tls.RequireAndVerifyClientCert, // 🔒 mTLS obligatory
			MinVersion:   tls.VersionTLS13,               // Force TLS 1.3
			Certificates: make([]tls.Certificate, 1),
		}
		logger.AppLog.Debug("TLS configuration created - mTLS will be enforced for all connections")

		// 3️⃣ Load the server certificate
		logger.AppLog.Debugf("Step 3/5: Loading server certificate pair - Cert: %s, Key: %s", certFile, keyFile)
		serverCert, err := tls.LoadX509KeyPair(certFile, keyFile)
		if err != nil {
			logger.AppLog.Errorf("Failed to load server certificate/key pair: %v", err)
			return fmt.Errorf("failed to load server certificate/key: %v", err)
		}
		tlsConfig.Certificates[0] = serverCert
		logger.AppLog.Debug("Server certificate and private key loaded successfully")

		// 4️⃣ Create the HTTPS server with this configuration
		logger.AppLog.Debugf("Step 4/5: Creating HTTPS server instance with mTLS configuration on %s", bindAddr)
		srv := &http.Server{
			Addr:      bindAddr,
			Handler:   router,
			TLSConfig: tlsConfig,
		}
		logger.AppLog.Debug("HTTPS server instance created with TLS configuration")

		// 5️⃣ Start the server with mTLS
		logger.AppLog.Debugf("Step 5/5: Starting HTTPS server with mTLS on %s", bindAddr)
		logger.AppLog.Info("Server will require client certificates signed by the configured CA")
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

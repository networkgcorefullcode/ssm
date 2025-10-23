package server

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"

	"github.com/networkgcorefullcode/ssm/factory"
	"github.com/networkgcorefullcode/ssm/handlers"
	"github.com/networkgcorefullcode/ssm/logger"
	"github.com/networkgcorefullcode/ssm/pkcs11mgr"
	"github.com/urfave/cli/v3"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type SSM struct{}

var SsmServer = &SSM{}

// TODO: create a proper server struct to hold server config if needed
var main_server http.Server

type (
	// Config information.
	Config struct {
		cfg string
	}
)

func New() (*SSM, error) {
	return &SSM{}, nil
}

func Get() *SSM {
	return SsmServer
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

	main_server = http.Server{
		Addr: factory.SsmConfig.Configuration.BindAddr,
	}
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

	logger.AppLog.Infof("Removing old socket at %s if exists", socketPath)
	_ = os.Remove(socketPath)

	logger.AppLog.Infof("Starting to listen on unix socket %s", socketPath)
	l, err := net.Listen("unix", socketPath)
	if err != nil {
		logger.AppLog.Errorf("Failed to listen on socket %s: %v", socketPath, err)
		return err
	}

	// Initialize PKCS11 connection pool
	poolConfig := pkcs11mgr.DefaultPoolConfig()
	poolConfig.PkcsPath = factory.SsmConfig.Configuration.PkcsPath
	poolConfig.SlotNumber = uint(factory.SsmConfig.Configuration.LotsNumber)
	poolConfig.Pin = factory.SsmConfig.Configuration.Pin
	poolConfig.MaxSize = 10 // Configure based on expected load
	poolConfig.MinSize = 2  // Minimum connections to maintain

	logger.AppLog.Info("Initializing PKCS11 connection pool...")
	if err := pkcs11mgr.InitializeGlobalPool(poolConfig); err != nil {
		logger.AppLog.Errorf("Failed to initialize PKCS11 connection pool: %v", err)
		return err
	}
	logger.AppLog.Info("PKCS11 connection pool initialized successfully")

	// init the pkcs manager
	// PkcsManager, err = pkcs11mgr.New(factory.SsmConfig.Configuration.PkcsPath,
	// 	uint(factory.SsmConfig.Configuration.LotsNumber),
	// 	factory.SsmConfig.Configuration.Pin)

	// SsmServer.mgr = PkcsManager

	// if err != nil {
	// 	logger.AppLog.Errorf("Failed to initialize PKCS11 manager: %v", err)
	// 	return err
	// }

	// err = PkcsManager.OpenSession()

	// if err != nil {
	// 	logger.AppLog.Errorf("Failed to OpenSession PKCS11 manager: %v", err)
	// 	return err
	// }

	// Pool monitoring endpoint
	http.HandleFunc("/pool/stats", func(w http.ResponseWriter, r *http.Request) {
		logger.AppLog.Debugf("Received /pool/stats request")
		handlers.HandlePoolStats(w, r)
	})

	// HealthCheck endpoint
	http.HandleFunc("/health-check", func(w http.ResponseWriter, r *http.Request) {
		logger.AppLog.Debugf("Received /health-check request")
		handlers.HandleHealthCheck(w, r)
	})

	if factory.SsmConfig.Configuration.HandlersPoolConect {
		CreateEndpointHandlersPool()
	} else {
		CreateEndpointHandlers()
	}

	// Serve HTTP requests in a separate goroutine
	logger.AppLog.Infof("SSM listening on unix socket %s", socketPath)
	go func() error {
		if err := http.Serve(l, nil); err != nil {
			logger.AppLog.Errorf("Server error: %v", err)
			return err
		}
		return nil
	}()

	if *factory.SsmConfig.Configuration.ExposeSwaggerUi {
		go func() {
			logger.AppLog.Infof("Swagger UI available at http://localhost:9001/swagger-ui")
			ServerSwagger()
		}()
	}

	// Start HTTPS or HTTP server based on configuration
	if factory.SsmConfig.Configuration.IsHttps == nil || *factory.SsmConfig.Configuration.IsHttps {
		// HTTPS server
		certFile := factory.SsmConfig.Configuration.CertFile
		keyFile := factory.SsmConfig.Configuration.KeyFile
		if certFile == "" || keyFile == "" {
			logger.AppLog.Error("HTTPS is enabled but certFile or keyFile is not set in the configuration")
			return fmt.Errorf("certFile or keyFile not set")
		}
		logger.AppLog.Infof("SSM listening api https %s", factory.SsmConfig.Configuration.BindAddr)
		// Use ListenAndServeTLS to handle HTTPS connections
		if err := http.ListenAndServeTLS(factory.SsmConfig.Configuration.BindAddr, certFile, keyFile, nil); err != nil {
			logger.AppLog.Errorf("Server error: %v", err)
			return err
		}
		return nil
	} else {
		logger.AppLog.Infof("SSM listening api http %s", factory.SsmConfig.Configuration.BindAddr)
		// Use ListenAndServe to handle HTTP connections
		if err := http.ListenAndServe(factory.SsmConfig.Configuration.BindAddr, nil); err != nil {
			logger.AppLog.Errorf("Server error: %v", err)
			return err
		}
	}

	// Close PKCS11 connection pool
	if pool := pkcs11mgr.GetGlobalPool(); pool != nil {
		logger.AppLog.Info("Closing PKCS11 connection pool...")
		pool.Close()
	}

	// PkcsManager.CloseSession()
	// PkcsManager.Finalize()

	logger.AppLog.Info("SSM server stopped gracefully")
	return nil
}

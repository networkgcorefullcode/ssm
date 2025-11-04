package server

import (
	"fmt"
	"net"
	"net/http"
	"os"
	"path/filepath"

	"github.com/networkgcorefullcode/ssm/database"
	"github.com/networkgcorefullcode/ssm/factory"
	"github.com/networkgcorefullcode/ssm/handlers"
	"github.com/networkgcorefullcode/ssm/logger"
	"github.com/networkgcorefullcode/ssm/pkcs11mgr"
	"github.com/networkgcorefullcode/ssm/server/middleware"
	"github.com/urfave/cli/v3"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type SSM struct{}

var SsmServer = &SSM{}

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

	// init the pkcs manager
	pkcsManager, err := pkcs11mgr.New(factory.SsmConfig.Configuration.PkcsPath,
		uint(factory.SsmConfig.Configuration.LotsNumber),
		factory.SsmConfig.Configuration.Pin)
	if err != nil {
		logger.AppLog.Errorf("Failed to initialize PKCS11 manager: %v", err)
		return err
	}

	pkcsManager.CloseAllSessions()
	pkcs11mgr.SetChanMaxSessions(factory.SsmConfig.Configuration.MaxSessions)

	handlers.SetPKCS11Manager(pkcsManager)
	middleware.SetPKCS11Manager(pkcsManager)
	pkcs11mgr.SetPKCS11Manager(pkcsManager)
	database.SetPKCS11Manager(pkcsManager)

	// Initialize PKCS11 functions and constants
	pkcs11mgr.InitPKCS11()
	// Initialize the database functions
	database.InitDB()

	// Build Gin router with all endpoints
	router := CreateGinRouter()

	// Serve HTTP requests in a separate goroutine
	logger.AppLog.Infof("SSM listening on unix socket %s", socketPath)
	go func() error {
		if err := http.Serve(l, router); err != nil {
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

	startHTTPServer(router)

	pkcsManager.CloseAllSessions()
	pkcsManager.Finalize()

	logger.AppLog.Info("SSM server stopped gracefully")
	return nil
}

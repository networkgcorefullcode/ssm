package server

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"

	"github.com/networkgcorefullcode/ssm/factory"
	"github.com/networkgcorefullcode/ssm/k4opt"
	"github.com/networkgcorefullcode/ssm/logger"
	"github.com/networkgcorefullcode/ssm/pkcs11mgr"
	"github.com/urfave/cli/v3"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type SSM struct {
	mgr *pkcs11mgr.Manager
}

var SsmServer = &SSM{}

type (
	// Config information.
	Config struct {
		cfg string
	}
)

func New(mgr *pkcs11mgr.Manager) (*SSM, error) {
	return &SSM{mgr: mgr}, nil
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

	_ = os.Remove(socketPath)
	l, err := net.Listen("unix", socketPath)
	if err != nil {
		return err
	}
	http.HandleFunc("/encrypt", k4opt.HandleEncryptK4)
	http.HandleFunc("/decrypt", func(w http.ResponseWriter, r *http.Request) {
		k4opt.HandleDecryptK4(s.mgr, w, r)
	})

	log.Printf("SSM listening on unix socket %s", socketPath)
	if err := http.Serve(l, nil); err != nil {
		log.Fatalf("server error: %v", err)
	}

	return nil
}

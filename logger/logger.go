package logger

import (
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	log         *zap.Logger
	AppLog      *zap.SugaredLogger
	InitLog     *zap.SugaredLogger
	CfgLog      *zap.SugaredLogger
	HttpLog     *zap.SugaredLogger
	ProducerLog *zap.SugaredLogger
	CommLog     *zap.SugaredLogger
	CallbackLog *zap.SugaredLogger
	UtilLog     *zap.SugaredLogger
	ConsumerLog *zap.SugaredLogger
	GinLog      *zap.SugaredLogger
	atomicLevel zap.AtomicLevel
)

const (
	FieldRanAddr     string = "ran_addr"
	FieldRanId       string = "ran_id"
	FieldAmfUeNgapID string = "amf_ue_ngap_id"
	FieldSupi        string = "supi"
	FieldSuci        string = "suci"
)

func init() {
	atomicLevel = zap.NewAtomicLevelAt(zap.InfoLevel)
	config := zap.Config{
		Level:            atomicLevel,
		Development:      false,
		Encoding:         "console",
		EncoderConfig:    zap.NewProductionEncoderConfig(),
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}

	config.EncoderConfig.TimeKey = "timestamp"
	config.EncoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	config.EncoderConfig.LevelKey = "level"
	config.EncoderConfig.EncodeLevel = zapcore.CapitalLevelEncoder
	config.EncoderConfig.CallerKey = "caller"
	config.EncoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	config.EncoderConfig.MessageKey = "message"
	config.EncoderConfig.StacktraceKey = ""

	var err error
	log, err = config.Build()
	if err != nil {
		panic(err)
	}

	AppLog = log.Sugar().With("component", "SSM", "category", "App")
	InitLog = log.Sugar().With("component", "SSM", "category", "Init")
	CfgLog = log.Sugar().With("component", "SSM", "category", "CFG")
	HttpLog = log.Sugar().With("component", "SSM", "category", "HTTP")
	ProducerLog = log.Sugar().With("component", "SSM", "category", "Producer")
	CommLog = log.Sugar().With("component", "SSM", "category", "Comm")
	CallbackLog = log.Sugar().With("component", "SSM", "category", "Callback")
	UtilLog = log.Sugar().With("component", "SSM", "category", "Util")
	ConsumerLog = log.Sugar().With("component", "SSM", "category", "Consumer")
	GinLog = log.Sugar().With("component", "SSM", "category", "GIN")
}

func GetLogger() *zap.Logger {
	return log
}

// SetLogLevel: set the log level (panic|fatal|error|warn|info|debug)
func SetLogLevel(level zapcore.Level) {
	CfgLog.Infoln("set log level:", level)
	atomicLevel.SetLevel(level)
}

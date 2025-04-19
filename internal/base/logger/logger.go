package logger

import (
	"graph-med/internal/base/conf"
	"io"
	"log"
	"os"
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var (
	logger      *zap.Logger
	currentDate string
	logConfig   *conf.AllConfig
	logMutex    sync.Mutex
)

func Initialize(config *conf.AllConfig) {
	logMutex.Lock()
	defer logMutex.Unlock()

	// 保存配置以便后续使用
	logConfig = config
	dir := config.Log.Output

	if err := os.MkdirAll(dir, os.ModePerm); err != nil {
		log.Fatal(err)
	}

	// 以当前时间作为日志文件名
	currentDate = time.Now().Format("2006-01-02")
	filename := dir + "/" + currentDate + ".log"

	// 日志输出到文件
	fileWriterSyncer := getLogWriter(filename)
	encoder := getEncoder()
	core := zapcore.NewCore(encoder, fileWriterSyncer, zapcore.DebugLevel)
	logger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))

	Info("日志初始化成功")
}

func getEncoder() zapcore.Encoder {
	return zapcore.NewJSONEncoder(zap.NewProductionEncoderConfig())
}

func getLogWriter(filename string) zapcore.WriteSyncer {
	// 使用追加模式打开文件，如果不存在则创建
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	ws := io.MultiWriter(file, os.Stdout)
	return zapcore.AddSync(ws)
}

// 检查日期是否变化，如果变化则重新初始化日志文件
func checkDate() {
	today := time.Now().Format("2006-01-02")
	if today != currentDate {
		// 日期已变化，需要重新初始化日志文件
		if logConfig != nil {
			Initialize(logConfig)
		}
	}
}

func Info(msg string, fields ...zap.Field) {
	checkDate()
	logger.Info(msg, fields...)
}

func Error(msg string, fields ...zap.Field) {
	checkDate()
	logger.Error(msg, fields...)
}

func Sync() error {
	return logger.Sync()
}

package main

import (
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"os"
	"time"

	"go.uber.org/zap"
)

const (
	kibanaConfigFile = "./config/kibana.json"
	logPath          = "./logs/go.log"
)

func main() {
	// Setup dashboards in Kibana (not required step)
	if err := setupDashboards(); err != nil {
		fmt.Printf("failed to setup Kibana dashboards, error: %s\n", err.Error())
	}

	// logger := WriteLog()

	// logger.Debug("i am debug", zap.String("key", "debug"))
	// logger.Info("i am info", zap.String("key", "info"))
	// logger.Error("i am error", zap.String("key", "error"))

	os.OpenFile(logPath, os.O_RDONLY|os.O_CREATE, 0666)
	c := zap.NewProductionConfig()
	c.OutputPaths = []string{"stdout", logPath}
	l, err := c.Build()
	if err != nil {
		panic(err)
	}
	i := 0
	for {
		i++
		time.Sleep(time.Second * 3)
		if rand.Intn(10) == 1 {
			l.Error("test error", zap.Error(fmt.Errorf("error because test: %d", i)))
		} else {
			l.Info(fmt.Sprintf("test log: %d", i))
		}
	}
}

// func getEncoder() zapcore.Encoder {
// 	return zapcore.NewConsoleEncoder(zapcore.EncoderConfig{
// 		MessageKey:   "message",
// 		TimeKey:      "time",
// 		LevelKey:     "level",
// 		CallerKey:    "caller",
// 		EncodeLevel:  CustomLevelEncoder,         //Format cách hiển thị level log
// 		EncodeTime:   SyslogTimeEncoder,          //Format hiển thị thời điểm log
// 		EncodeCaller: zapcore.ShortCallerEncoder, //Format dòng code bắt đầu log
// 	})
// }

// func SyslogTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
// 	enc.AppendString(t.Format("2006-01-02 15:04:05"))
// }

// func CustomLevelEncoder(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
// 	enc.AppendString("[" + level.CapitalString() + "]")
// }

// func logErrorWriter() zapcore.WriteSyncer {
// 	errFileLog, _ := os.OpenFile("./logs/error.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

// 	return zapcore.NewMultiWriteSyncer(
// 		zapcore.AddSync(&lumberjack.Logger{
// 			Filename: errFileLog.Name(),
// 			MaxSize:  500, // megabytes
// 			MaxAge:   30,  // days
// 		}),
// 		zapcore.AddSync(os.Stdout))
// }

// func logInfoWriter() zapcore.WriteSyncer {
// 	infoFileLog, _ := os.OpenFile("./logs/info.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

// 	return zapcore.NewMultiWriteSyncer(
// 		zapcore.AddSync(&lumberjack.Logger{
// 			Filename: infoFileLog.Name(),
// 			MaxSize:  500, // megabytes
// 			MaxAge:   30,  // days
// 		}),
// 		zapcore.AddSync(os.Stdout))
// }

// func logDebugWriter() zapcore.WriteSyncer {
// 	debugFileLog, _ := os.OpenFile("./logs/debug.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)

// 	return zapcore.NewMultiWriteSyncer(
// 		zapcore.AddSync(&lumberjack.Logger{
// 			Filename: debugFileLog.Name(),
// 			MaxSize:  500, // megabytes
// 			MaxAge:   30,  // days
// 		}),
// 		zapcore.AddSync(os.Stdout))
// }

// // Write log to file by level log and console
// func WriteLog() *zap.Logger {
// 	highWriteSyncer := logErrorWriter()
// 	averageWriteSyncer := logDebugWriter()
// 	lowWriteSyncer := logInfoWriter()

// 	encoder := getEncoder()

// 	highPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
// 		return lev >= zap.ErrorLevel
// 	})

// 	lowPriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
// 		return lev < zap.ErrorLevel && lev > zap.DebugLevel
// 	})

// 	averagePriority := zap.LevelEnablerFunc(func(lev zapcore.Level) bool {
// 		return lev < zap.ErrorLevel && lev < zap.InfoLevel
// 	})

// 	lowCore := zapcore.NewCore(encoder, lowWriteSyncer, lowPriority)
// 	averageCore := zapcore.NewCore(encoder, averageWriteSyncer, averagePriority)
// 	highCore := zapcore.NewCore(encoder, highWriteSyncer, highPriority)

// 	logger := zap.New(zapcore.NewTee(lowCore, averageCore, highCore), zap.AddCaller())
// 	return logger
// }

// setupDashboards put graphs and dashboards inside kibana
func setupDashboards() error {
	f, err := os.Open(kibanaConfigFile)
	if err != nil {
		return err
	}
	defer f.Close()
	url := "http://localhost:5601/api/kibana/dashboards/import"
	req, err := http.NewRequest("POST", url, f)
	if err != nil {
		return err
	}

	req.Header.Add("Kbn-Xsrf", "true")
	req.Header.Add("Content-Type", "application/json")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}

	if res.StatusCode != 200 {
		return fmt.Errorf("error responce from kibana: %s", string(body))
	}
	return nil
}

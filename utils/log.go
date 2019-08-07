package utils

import (
	"encoding/json"
	"fmt"
	"os"
	"strings"

	"github.com/sirupsen/logrus"
)

const logLevel = logrus.DebugLevel
const MAXLOGLINES = 100

type (
	logStream struct {
		buf []string
	}

	logData struct {
		Level string `json:"level"`
		Msg   string `json:"msg"`
		Time  string `json:"time"`
	}
)

var stream logStream
var logFile *os.File

func (w *logStream) Write(p []byte) (n int, err error) {
	var data logData
	json.Unmarshal(p, &data)
	str := fmt.Sprintf("%s %s %s\n", data.Time, strings.ToUpper(data.Level), data.Msg)
	os.Stdout.Write([]byte(str))
	if w.buf == nil {
		w.buf = []string{}
	}
	if logFile == nil {
		logFile, err = os.OpenFile("logrus.log", os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0777)
		if err != nil {
			logrus.Info("Failed to log to file")
		}
	} else {
		logFile.WriteString(str)
	}
	str = string(p)
	w.buf = append(w.buf, str)
	if len(w.buf) > MAXLOGLINES {
		w.buf = w.buf[len(w.buf)-MAXLOGLINES:]
	}
	return len(p), nil
}

func (w *logStream) Read() string {
	if w.buf == nil {
		w.buf = []string{}
	}
	return strings.Join(w.buf, "")
}

func ReadAllLogs() string {
	return stream.Read()
}

func CreateLogger(module ...string) *logrus.Entry {
	log := logrus.New()
	log.SetFormatter(&logrus.JSONFormatter{
		TimestampFormat: "01-02 15:04:05",
	})
	// log.SetReportCaller(true)
	log.SetLevel(logLevel)
	log.Out = &stream
	return logrus.NewEntry(log)
	// if len(module) > 0 {
	// 	return logrus.WithField("module", module[0])
	// }
	// return logrus.WithField("module", nil)
}

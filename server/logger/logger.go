package logger
import (
	"os"
	"github.com/sirupsen/logrus"
)
var log = logrus.New()

func InitLogger() (*os.File, error){
	removeLogFile("app.log")
	file, err := os.OpenFile("app.log", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err == nil {
        log.Out = file
    } else {
        log.Info("Failed to log to file, using default stderr")
		return nil, err
    }
	log.SetLevel(logrus.InfoLevel)
	log.SetFormatter(&logrus.TextFormatter{
		FullTimestamp: true,
	})
	return file, nil
}

func removeLogFile(fileName string) {
    if _, err := os.Stat(fileName); err == nil {
        err := os.Remove(fileName)
        if err != nil {
            log.Printf("Failed to remove old log file: %v\n", err)
        }
    }
}

func Info(info... any) {
	log.Info(info...)
}

func Fatal(info... any) {
	log.Fatal(info...)
}
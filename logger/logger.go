package logger

import (
	"fmt"
	"log"
	"os"
	"runtime"
)

const LogFileName string = "blueprint-linter.log"

var (
	traceLogger    *log.Logger
	infoLogger     *log.Logger
	warningLogger  *log.Logger
	fatalLogger    *log.Logger
	LoggingEnabled bool
)

func InitLogger() {
	fileName := ""
	switch runtime.GOOS {
	case "windows":
		logsDir := fmt.Sprintf("%s\\%s", os.Getenv("TEMP"), "GenesysCloud")
		if err := mkdirIfNotExist(logsDir); err != nil {
			return
		}
		fileName = fmt.Sprintf("%s\\%s", logsDir, LogFileName)
	case "darwin":
		homeDir, _ := os.UserHomeDir()
		logsDir := fmt.Sprintf("%s/Library/Logs/GenesysCloud", homeDir)
		if err := mkdirIfNotExist(logsDir); err != nil {
			return
		}
		fileName = fmt.Sprintf("%s/%s", logsDir, LogFileName)
	default:
		logsDir := "/tmp/GenesysCloud"
		if err := mkdirIfNotExist(logsDir); err != nil {
			return
		}
		fileName = fmt.Sprintf("%s/%s", logsDir, LogFileName)
	}

	// If the file doesn't exist, create it or append to the file
	file, err := os.OpenFile(fileName, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return
	}

	flags := log.Ldate | log.Ltime
	traceLogger = log.New(file, "TRACE: ", flags)
	infoLogger = log.New(file, "INFO: ", flags)
	warningLogger = log.New(file, "WARNING: ", flags)
	fatalLogger = log.New(file, "FATAL: ", flags)

	traceLogger.SetOutput(file)
	infoLogger.SetOutput(file)
	warningLogger.SetOutput(file)
	traceLogger.SetOutput(file)

	log.SetOutput(file)
}

func Trace(v ...interface{}) {
	fmt.Fprint(os.Stderr, v...)
	if traceLogger != nil && LoggingEnabled {
		traceLogger.Println(v...)
	}
}

func Tracef(format string, v ...interface{}) {
	fmt.Fprintf(os.Stderr, format, v...)
	if traceLogger != nil && LoggingEnabled {
		traceLogger.Printf(format, v...)
	}
}

func Info(v ...interface{}) {
	if infoLogger != nil && LoggingEnabled {
		infoLogger.Println(v...)
	}
}

func Infof(format string, v ...interface{}) {
	if infoLogger != nil && LoggingEnabled {
		infoLogger.Printf(format, v...)
	}
}

func Warn(v ...interface{}) {
	if warningLogger != nil && LoggingEnabled {
		warningLogger.Println(v...)
	}
}

func Warnf(format string, v ...interface{}) {
	if warningLogger != nil && LoggingEnabled {
		warningLogger.Printf(format, v...)
	}
}

func Fatal(v ...interface{}) {
	fmt.Fprint(os.Stderr, v...)
	if fatalLogger != nil && LoggingEnabled {
		fatalLogger.Fatal(v...)
	} else {
		os.Exit(1)
	}
}

func Fatalf(format string, v ...interface{}) {
	fmt.Fprintf(os.Stderr, format, v...)
	if fatalLogger != nil && LoggingEnabled {
		fatalLogger.Fatalf(format, v...)
	} else {
		os.Exit(1)
	}
}

func mkdirIfNotExist(directory string) error {
	if _, err := os.Stat(directory); os.IsNotExist(err) {
		err = os.Mkdir(directory, 0755)
		if err != nil {
			return err
		}
	}

	return nil
}

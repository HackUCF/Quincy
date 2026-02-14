/*
Package log configures a global logger and exports common logging functions. It is a very light wrapper over [Zap].
This logger utilizes an init() function and has no internal dependencies. This makes it easy to import and start using.
Logs are always written to stdout. If the QU_LOG_FILE environment variable is set, logs are also appended to that file.
Each process instance is tagged with a unique "restart" UUID so log entries can be correlated across restarts.
This package might need more environment variables to configure some common settings, like log level and stack traces.

[Zap]: https://github.com/uber-go/zap
*/
package log

import (
	"os"

	"github.com/google/uuid"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

const level = zap.DebugLevel

var logger *zap.SugaredLogger

// build the base logger object
// takes an optional filepath that can be used
func buildLogger(level zapcore.Level, logFilePath string) *zap.Logger {
	encoderConfig := zap.NewProductionEncoderConfig()

	// list of destinations for logs
	cores := []zapcore.Core{
		// always contains stdout
		zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), zapcore.AddSync(os.Stdout), level),
	}

	// optionally contains a file
	if logFilePath != "" {
		file, err := os.OpenFile(logFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			panic(err)
		}
		cores = append(cores, zapcore.NewCore(zapcore.NewJSONEncoder(encoderConfig), zapcore.AddSync(file), level))
	}

	return zap.New(zapcore.NewTee(cores...),
		zap.AddCaller(),
		zap.AddCallerSkip(1),
	)
}

func init() {
	baseLogger := buildLogger(level, os.Getenv("QU_LOG_FILE"))
	// add a uuid to identify logs between restarts
	instanceID := uuid.New().String()
	logger = baseLogger.Sugar().With("restart", instanceID)
	Info("logger initialized")
}

/*
Debug logs a message with level == 'debug'.
Takes arbitrary key value pairs which are added as extra fields in the output.

	log.Debug(
	  "nothing happened",  // message
	  "honk", "mimimimi",  // key value pairs
	  "number", 1,
	  "key_without_value", // ignored
	)
*/
func Debug(msg string, keysAndValues ...any) {
	logger.Debugw(msg, keysAndValues...)
}

/*
Info logs a message with level == 'info'.
Takes arbitrary key value pairs which are added as extra fields in the output.

	log.Info(
	  "something happended", // message
	  "nothing", "gets",     // key value pairs
	  "past", "my bow!",
	  "key_without_value",   // ignored
	)
*/
func Info(msg string, keysAndValues ...any) {
	logger.Infow(msg, keysAndValues...)
}

/*
Warn logs a message with level == 'warn'.
Takes arbitrary key value pairs which are added as extra fields in the output.

	log.Warn(
	  "something weird",   // message
	  "error", err,        // key value pairs
	  "number", 2,
	  "key_without_value", // ignored
	)
*/
func Warn(msg string, keysAndValues ...any) {
	logger.Warnw(msg, keysAndValues...)
}

/*
Error logs a message with level == 'error'.
Takes arbitrary key value pairs which are added as extra fields in the output.

	log.Error(
	  "something failed",  // message
	  "error", err,        // key value pairs
	  "number", 3,
	  "key_without_value", // ignored
	)
*/
func Error(msg string, keysAndValues ...any) {
	logger.Errorw(msg, keysAndValues...)
}

/*
Panic logs a message with level == 'panic'. It immediately calls panic() after.
This is safe to call in DB functions and Gin Routes; the recovery middleware will handle it.
Takes arbitrary key value pairs which are added as extra fields in the output.

	log.Panic(
	  "i died",            // message
	  "error", err,        // key value pairs
	  "number", 4,
	  "key_without_value", // ignored
	)
*/
func Panic(msg string, keysAndValues ...any) {
	logger.Panicw(msg, keysAndValues...)
}

/*
Package log configures a global logger and exports common logging functions. It is a very light wrapper over [Zap].
This logger utilizes an init() function and has no internal dependencies. This makes it easy to import and start using.
This package might need environment variables to configure some common settings, like log level and stack traces.

[Zap]: https://github.com/uber-go/zap
*/
package log

import (
	"go.uber.org/zap"
)

var logger *zap.SugaredLogger

func init() {
	zapConfig := zap.NewProductionConfig()

	// we don't want big ugly stack traces
	zapConfig.DisableStacktrace = true

	// show all logs for now
	zapConfig.Level = zap.NewAtomicLevelAt(zap.DebugLevel)

	baseLogger, err := zapConfig.Build(
		zap.AddCallerSkip(1), // skip marking calls from this file
	)
	if err != nil {
		panic(err)
	}

	logger = baseLogger.Sugar()

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

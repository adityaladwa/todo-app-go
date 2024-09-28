package logger

import (
    "net/http"
    "strings"
    "time"

    "github.com/sirupsen/logrus"
)


func NewLogger(level string) *logrus.Logger {
    log := logrus.New()
    lvl, err := logrus.ParseLevel(strings.ToLower(level))
    if err != nil {
        log.SetLevel(logrus.InfoLevel)
    } else {
        log.SetLevel(lvl)
    }
    log.SetFormatter(&logrus.TextFormatter{
        ForceColors: true,
    })
    return log
}

func LogrusMiddleware(logger *logrus.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Call the next handler in the chain
			next.ServeHTTP(w, r)

			// Log the request details using the custom logrus logger
			logger.WithFields(logrus.Fields{
				"method":  r.Method,
				"url":     r.URL.Path,
				"elapsed": time.Since(start),
			}).Info("Handled request")
		})
	}
}
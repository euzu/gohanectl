package utils

import (
	"fmt"
	"github.com/go-chi/chi/middleware"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"io"
	"net/http"
	"os"
	"strconv"
	"time"
)

func getLogLevel(lvl *string) zerolog.Level {
	switch *lvl {
	case "panic":
		return zerolog.PanicLevel
	case "fatal":
		return zerolog.FatalLevel
	case "error":
		return zerolog.ErrorLevel
	case "warn":
		return zerolog.WarnLevel
	case "info":
		return zerolog.InfoLevel
	default:
		return zerolog.DebugLevel
	}
}

var myHttpLogger *zerolog.Logger

func InitLogger(logLevel *string, logFileName string, logConsole bool) {
	log.Info().Msgf("Logging level: %s,  file: %s", getLogLevel(logLevel), logFileName)
	zerolog.TimeFieldFormat = time.RFC3339 // "2006-01-02 15:04:05"
	zerolog.SetGlobalLevel(getLogLevel(logLevel))
	if logFileName != "" {
		f, err := os.OpenFile(logFileName, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)
		if err == nil {
			var wrt io.Writer = f
			if logConsole {
				wrt = io.MultiWriter(os.Stderr, f)
			}
			log.Logger = zerolog.New(wrt).With().Timestamp().Logger()
			logger := zerolog.New(wrt).With().Timestamp().Logger()
			myHttpLogger = &logger
		}
	}

	if myHttpLogger == nil {
		logger := zerolog.New(os.Stderr).With().Timestamp().Logger()
		myHttpLogger = &logger
	}
}

type customLogEntry struct {
	request *http.Request
}

func (e *customLogEntry) Write(status, bytes int, header http.Header, elapsed time.Duration, extra interface{}) {
	if myHttpLogger.GetLevel() == zerolog.DebugLevel {
		r := e.request
		l := myHttpLogger.With().Logger()
		l.UpdateContext(func(c zerolog.Context) zerolog.Context {
			return c.Str("at", time.Now().Format("2006-01-02 15:04:05")).
				Str("method", r.Method).
				Str("uri", r.URL.RequestURI()).
				Str("proto", r.Proto).
				Str("host", r.Host).
				Str("remoteAddr", r.RemoteAddr).
				Str("status", strconv.Itoa(status)).
				Str("bytes", strconv.Itoa(bytes)).
				Str("took", fmt.Sprintf("%d", elapsed.Milliseconds()))
		})

		l.Debug().Msg("")
	}
}

func (e *customLogEntry) Panic(v interface{}, stack []byte) {
	//panicEntry := l.NewLogEntry(l.request).(*defaultLogEntry)
	//cW(panicEntry.buf, l.useColor, bRed, "panic: %+v", v)
	//l.Logger.Print(panicEntry.buf.String())
	//l.Logger.Print(string(stack))
}

type customLogFormatter struct {
}

// NewLogEntry creates a new LogEntry for the request.
func (l *customLogFormatter) newLogEntry(r *http.Request) middleware.LogEntry {
	entry := &customLogEntry{
		request: r,
	}
	return entry
}

func customRequestLogger(f *customLogFormatter) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		fn := func(w http.ResponseWriter, r *http.Request) {
			entry := f.newLogEntry(r)
			ww := middleware.NewWrapResponseWriter(w, r.ProtoMajor)

			t1 := time.Now()
			defer func() {
				entry.Write(ww.Status(), ww.BytesWritten(), nil, time.Since(t1), nil)
			}()

			next.ServeHTTP(ww, middleware.WithLogEntry(r, entry))
		}
		return http.HandlerFunc(fn)
	}
}

func HttpLogger(next http.Handler) http.Handler {
	httpLogger := customRequestLogger(&customLogFormatter{})
	return httpLogger(next)
}

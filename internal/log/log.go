package log

import (
	"app/internal/config"
	"context"
	"fmt"
	"log/slog"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/dusted-go/logging/prettylog"
)

const (
    LevelDebug = slog.LevelDebug
    LevelInfo = slog.LevelInfo
    LevelWarn = slog.LevelWarn
    LevelError = slog.LevelError
    LevelPanic = slog.Level(12)
    LevelFatal = slog.Level(16)
)

var logger *slog.Logger = nil
var dbWritter *DBLogWritter = nil

// type CustomHandler struct {
// 	slog.JSONHandler
// }
// 
// func (h *CustomHandler) Handle(ctx context.Context, r slog.Record) error {
// 	switch r.Level {
// 	case LevelPanic:
// 		r.Level = LevelPanic
//         r.Add("level", "PANIC")
// 		//r.AddAttrs(slog.String("level", "PANIC"))
// 	case LevelFatal:
// 		r.Level = LevelFatal
//         r.Add("level", "FATAL")
// 		//r.AddAttrs(slog.String("level", "FATAL"))
// 	}
// 	return h.JSONHandler.Handle(ctx, r)
// }

func Init() {
    logLevel := slog.LevelInfo
    if config.Log.Debug {
        logLevel = slog.LevelDebug
    }
    handlerOption := &slog.HandlerOptions{
        Level: logLevel,

		ReplaceAttr: func(groups []string, a slog.Attr) slog.Attr {
			// if a.Key == slog.TimeKey {
			// 	return slog.Attr{}
			// }

            if a.Key == slog.LevelKey {
				a.Key = "level"
				level := a.Value.Any().(slog.Level)
				switch level {
                case LevelDebug:
                    a.Value = slog.StringValue("DEBUG")
                case LevelInfo:
                    a.Value = slog.StringValue("INFO")
                case LevelWarn:
                    a.Value = slog.StringValue("WARN")
                case LevelError:
                    a.Value = slog.StringValue("ERROR")
                case LevelPanic:
                    a.Value = slog.StringValue("PANIC")
                case LevelFatal:
                    a.Value = slog.StringValue("FATAL")
				}
			}

			return a
		},
    }

    var handler slog.Handler
    switch strings.ToLower(config.Log.Style) {
    case "text":
        handler = slog.NewTextHandler(os.Stdout, handlerOption)
    case "json":
        handler = slog.NewJSONHandler(os.Stdout, handlerOption)
    case "pretty":
        handler = prettylog.NewHandler(handlerOption) 
    default:
        handler = NewHandler(handlerOption)
        // handler = &CustomHandler{*slog.NewJSONHandler(os.Stdout, handlerOption)}
    }

    logger = slog.New(handler)

    if config.Log.DumpSqlite {
        dbWritter = NewDBWritter()
    }
}

func log(ctx context.Context, level slog.Level, msg string, args ...any) {
    if logger == nil { Init() }
    requestId, ok := ctx.Value("requestId").(string)
    if ok {
        args = append(args, slog.String("requestId", requestId))
    }
    switch level {
        case LevelDebug:
            logger.Debug(msg, args...)
        case LevelInfo:
            logger.Info(msg, args...)
        case LevelWarn:
            logger.Warn(msg, args...)
        case LevelError:
            logger.Error(msg, args...)
        case LevelPanic:
            logger.Log(ctx, LevelPanic, msg, args...)
            panic(msg)
        case LevelFatal:
            logger.Log(ctx, LevelFatal, msg, args...)
            os.Exit(1)
        default: 
            logger.Info(msg, args...)
    }

    timeStamp := time.Now().UTC()

    go writeToDB(requestId, level.String(), msg, timeStamp, args)
}

func writeToDB(requestId string, level string, msg string, timeStamp time.Time, args []any) {
    attrs := map[string]any{} 
    
    for i := 0; i < len(args); i++ {
        if a, ok := args[i].(slog.Attr); ok {
            attrs[a.Key] = a.Value.Any()
        } else if a, ok := args[i].(string); ok {
            i++
            if i < len(args) {
                attrs[a] = args[i]
            } else {
                attrs[a] = ""
            }
        }
    }

    dbWritter.Write(requestId, level, msg, timeStamp, attrs)
}

func Info(ctx context.Context, msg string, args ...any) {
    log(ctx, LevelInfo, msg, args...)
}

func Error(ctx context.Context, msg string, args ...any) {
    log(ctx, LevelError, msg, args...)
}

func Warn(ctx context.Context, msg string, args ...any) {
    log(ctx, LevelWarn, msg, args...)
}

func Debug(ctx context.Context, msg string, args ...any) {
    log(ctx, LevelDebug, msg, args...)
}

func Panic(ctx context.Context, msg string, args ...any) {
    log(ctx, LevelPanic, msg, args...)
}

func Fatal(ctx context.Context, msg string, args ...any) {
    log(ctx, LevelFatal, msg, args...)
}

func EnterMethod(ctx context.Context, args ...any) {
    msg := "Enter method: " + getMethodName()
    log(ctx, LevelInfo, msg, args...)
}

func ExitMethod(ctx context.Context, args ...any) {
    msg := "Exit method: " + getMethodName()
    log(ctx, LevelInfo, msg, args...)
}

func getMethodName() string {
	pc, _, _, _ := runtime.Caller(2)
	return fmt.Sprintf("%s", runtime.FuncForPC(pc).Name())
}

package logger

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/gookit/slog"
	"github.com/gookit/slog/handler"
)

const (
	Template string = "[{{datetime}}] [{{level}}] {{message}} {{data}} {{extra}}\n"
)

type MyError struct {
	Message string
}

func (e MyError) Error() string {
	return e.Message
}

func NewError(msg string) error {
	return MyError{Message: msg}
}

func Err(err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf(getCallerInfo() + " " + err.Error())
}

func Info(err error) {
	defer slog.MustClose()
	l := setSettingsInfo()
	l.Info(getCallerInfo(), err)
}

func Error(err error) {
	defer slog.MustClose()
	l := setSettingsError()
	l.Error(getCallerInfo(), err)
}

func Fatal(err error) {
	defer slog.MustClose()
	l := setSettingsError()
	l.Fatal(getCallerInfo(), err)
}

func setSettingsError() *slog.Logger {
	f := slog.NewTextFormatter(Template)
	filename := time.Now().Format("01.01.2000") + "-errors"
	h, _ := handler.NewFileHandler("errors\\"+filename+".log", handler.WithLogLevels(slog.DangerLevels))
	h.SetFormatter(f)
	l := slog.NewWithHandlers(h)
	return l
}

func setSettingsInfo() *slog.Logger {
	f := slog.NewTextFormatter(Template)
	h := handler.NewConsoleHandler(slog.NormalLevels)
	h.SetFormatter(f)
	l := slog.NewWithHandlers(h)
	return l
}

func getCallerInfo() string {
	_, file, line, _ := runtime.Caller(2)
	parts := strings.Split(file, "/")
	file = parts[len(parts)-1]
	return "[" + file + ":" + strconv.Itoa(line) + "]"
}

package logger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"text/template"
	"time"
)

type Level uint32

const (
	LevelDefault Level = iota
	LevelFatal
	LevelError
	LevelWarning
	LevelInfo
	LevelDebug
	LevelTrace

	LevelCritical = LevelFatal
)

func (l Level) String() string {
	switch l {
	case LevelDefault:
		return "Default"

	case LevelFatal:
		return "CRITICAL"

	case LevelError:
		return "ERROR"

	case LevelWarning:
		return "WARNING"

	case LevelInfo:
		return "INFO"

	case LevelDebug:
		return "DEBUG"

	case LevelTrace:
		return "TRACE"
	}

	return "Unknown"
}

func StringLevel(s string) Level {
	switch strings.ToLower(s) {
	case "default":
		return LevelDefault

	case "critical":
		return LevelCritical

	case "fatal":
		return LevelFatal

	case "error":
		return LevelError

	case "warning":
		return LevelWarning

	case "info":
		return LevelInfo

	case "debug":
		return LevelDebug

	case "trace":
		return LevelTrace
	}

	return LevelDefault
}

type Context struct {
	wtr  io.Writer
	tmpl *template.Template
	lvl  Level
}

var (
	DefaultFormat  = `{"timestamp":"{{eenTimeStamp .Now}}","logLevel":"{{.Level}}","logFacility":"main","triggerLabel":null,"function":"{{.Function}}","file":"{{.Filename}}","lineNo":{{.LineNo}},"message":"{{.Message}}","extra":{{toJson .Extra}}}` //"{{eenTimeStamp .Now}}[{{.Level}}][NULL, {{.Function}}({{.Filename}}:{{.LineNo}})][{{.Trigger}}]: {{.Message}} "
	DefaultContext = NewContext(os.Stderr, DefaultFormat, LevelWarning)
)

type Logger struct {
	*Context
	subfacility string
	lvl         Level
	buf         []byte
	mux         sync.Mutex
}

func eenTimeStamp(t time.Time) string {
	return t.Format("2006-01-02T15:04:05.000Z")
}

func toJson(v interface{}) (string, error) {
	bytes, err := json.Marshal(v)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
func NewContext(w io.Writer, fmt string, lvl Level) *Context {
	tmpl := template.New("LogPrefix")
	tmpl = tmpl.Funcs(map[string]interface{}{"eenTimeStamp": eenTimeStamp, "toJson": toJson})
	tmpl = template.Must(tmpl.Parse(fmt))
	return &Context{wtr: w,
		tmpl: tmpl,
		lvl:  lvl}
}

func (ctx *Context) SetLevel(lvl Level) {
	atomic.StoreUint32((*uint32)(&ctx.lvl), uint32(lvl))
}

func (ctx *Context) Level() Level {
	return Level(atomic.LoadUint32((*uint32)(&ctx.lvl)))
}

// Get the logger context
func (ctx *Context) GetLogger(sub string, lvl Level) *Logger {
	return &Logger{Context: ctx,
		subfacility: sub,
		lvl:         lvl,
		mux:         sync.Mutex{}}
}

func (l *Logger) SetLevel(lvl Level) {
	l.mux.Lock()
	defer l.mux.Unlock()
	l.lvl = lvl
}

func (l *Logger) Level() Level {
	l.mux.Lock()
	defer l.mux.Unlock()
	return l.lvl
}

func (l *Logger) Output(calldepth int, lvl Level, trigger string, m string) error {
	now := time.Now()
	pc, filename, lineno, _ := runtime.Caller(calldepth)
	f := runtime.FuncForPC(pc)
	if trigger == "" {
		trigger = "NULL"
	}

	//Variable replacement
	vars := map[string]interface{}{"Now": now,
		"Level":    lvl,
		"Function": f.Name(),
		"Filename": filepath.Base(filename),
		"LineNo":   lineno,
		"Trigger":  trigger,
		"Message":  strings.ReplaceAll(m, `"`, `\"`)}

	//Execute template
	l.mux.Lock()
	defer l.mux.Unlock()
	var tbuf bytes.Buffer
	l.buf = l.buf[:0]
	if err := l.tmpl.ExecuteTemplate(&tbuf, "LogPrefix", vars); err != nil {
		return err
	}

	l.buf = append(l.buf, tbuf.Bytes()...)
	if len(l.buf) > 0 && l.buf[len(l.buf)-1] != '\n' {
		l.buf = append(l.buf, '\n')
	}

	_, err := l.wtr.Write(l.buf)
	return err
}

func (l *Logger) CheckLevel(lvl Level) bool {
	llvl := l.lvl
	if llvl == LevelDefault {
		llvl = l.Context.Level()
	}

	return llvl >= lvl
}

func (l *Logger) Log(lvl Level, v ...interface{}) {
	if lvl == LevelFatal || l.CheckLevel(lvl) {
		err := l.Output(2, lvl, "", fmt.Sprint(v...))
		if err != nil {
			panic(fmt.Sprint(v...))
		}
		if lvl == LevelFatal {
			panic(fmt.Sprint(v...))
		}
	}
}

func (l *Logger) Logf(lvl Level, str string, v ...interface{}) {
	if lvl == LevelFatal || l.CheckLevel(lvl) {
		err := l.Output(2, lvl, "", fmt.Sprintf(str, v...))
		if err != nil {
			panic(fmt.Sprintf(str, v...))
		}
		if lvl == LevelFatal {
			panic(fmt.Sprintf(str, v...))
		}
	}
}

func (l *Logger) Trigger(lvl Level, trigger string, v ...interface{}) {
	msg := fmt.Sprint(v...)
	err := l.Output(2, lvl, trigger, msg)
	if err != nil {
		panic(msg)
	}
	if lvl == LevelFatal {
		panic(msg)
	}
}

func (l *Logger) Triggerf(lvl Level, trigger string, msg string, v ...interface{}) {
	msg = fmt.Sprintf(msg, v...)
	err := l.Output(2, lvl, trigger, msg)
	if err != nil {
		panic(msg)
	}

	if lvl == LevelFatal {
		panic(msg)
	}
}

func (l *Logger) Fatal(v ...interface{}) {
	s := fmt.Sprint(v...)
	err := l.Output(2, LevelFatal, "", s)
	if err != nil {
		panic(s)
	}

	panic(s)
}

func (l *Logger) Fatalf(str string, v ...interface{}) {
	s := fmt.Sprintf(str, v...)
	err := l.Output(2, LevelFatal, "", s)
	if err != nil {
		panic(s)
	}

	panic(s)
}

func (l *Logger) Error(v ...interface{}) {
	if l.CheckLevel(LevelError) {
		err := l.Output(2, LevelError, "", fmt.Sprint(v...))
		if err != nil {
			panic(v)
		}

	}
}

func (l *Logger) Errorf(str string, v ...interface{}) {
	if l.CheckLevel(LevelError) {
		err := l.Output(2, LevelError, "", fmt.Sprintf(str, v...))
		if err != nil {
			panic(v)
		}

	}
}

func (l *Logger) Warning(v ...interface{}) {
	if l.CheckLevel(LevelWarning) {
		err := l.Output(2, LevelWarning, "", fmt.Sprint(v...))
		if err != nil {
			panic(v)
		}

	}
}

func (l *Logger) Warningf(str string, v ...interface{}) {
	if l.CheckLevel(LevelWarning) {
		err := l.Output(2, LevelWarning, "", fmt.Sprintf(str, v...))
		if err != nil {
			panic(v)
		}

	}
}

func (l *Logger) Info(v ...interface{}) {
	if l.CheckLevel(LevelInfo) {
		err := l.Output(2, LevelInfo, "", fmt.Sprint(v...))
		if err != nil {
			panic(v)
		}

	}
}

func (l *Logger) Infof(str string, v ...interface{}) {
	if l.CheckLevel(LevelInfo) {
		err := l.Output(2, LevelInfo, "", fmt.Sprintf(str, v...))
		if err != nil {
			panic(v)
		}

	}
}

func (l *Logger) Debug(v ...interface{}) {
	if l.CheckLevel(LevelDebug) {
		err := l.Output(2, LevelDebug, "", fmt.Sprint(v...))
		if err != nil {
			panic(v)
		}

	}
}

func (l *Logger) Debugf(str string, v ...interface{}) {
	if l.CheckLevel(LevelDebug) {
		err := l.Output(2, LevelDebug, "", fmt.Sprintf(str, v...))
		if err != nil {
			panic(v)
		}

	}
}

func (l *Logger) Trace(v ...interface{}) {
	if l.CheckLevel(LevelTrace) {
		err := l.Output(2, LevelTrace, "", fmt.Sprint(v...))
		if err != nil {
			panic(v)
		}

	}
}

func (l *Logger) Tracef(str string, v ...interface{}) {
	if l.CheckLevel(LevelTrace) {
		err := l.Output(2, LevelTrace, "", fmt.Sprintf(str, v...))
		if err != nil {
			panic(v)
		}

	}
}

func (l *Logger) Panic(v ...interface{}) {
	s := fmt.Sprint(v...)
	err := l.Output(2, LevelFatal, "", s)
	if err != nil {
		panic(s)
	}

	panic(s)
}

func (l *Logger) Panicf(str string, v ...interface{}) {
	s := fmt.Sprintf(str, v...)
	err := l.Output(2, LevelFatal, "", s)
	if err != nil {
		panic(s)
	}

	panic(s)
}

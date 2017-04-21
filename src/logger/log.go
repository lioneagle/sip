package logger

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path"
	"runtime"
	"sync"
	"time"
)

const STACK_BUFFER_SIZE int = 8192
const TIMESTAMP_FMT string = "2006-01-02 15:04:05.000"

const (
	EMERGENCY = 0
	ALERT     = 1
	CRITICAL  = 2
	ERROR     = 3
	WARNING   = 4
	NOTICE    = 5
	INFO      = 6
	DEBUG     = 7
)

type Level int

var levelNames = []string{"EMERGENCY", "ALERT", "CRITICAL", "ERROR", "WARNING", "NOTICE", "INFO", "DEBUG"}

func (level Level) String() string {
	if level < 0 || int(level) >= len(levelNames) {
		return "???"
	}
	return levelNames[level]
}

type Logger struct {
	showShorFile      bool
	shorFileNameDepth int
	showPackage       bool
	showFuncName      bool
	out               io.Writer  // destination for output
	mu                sync.Mutex // ensures atomic writes; protects the following fields
	level             Level
	stackTraceLevel   Level
	modules           map[int]string
}

func NewLogger(out io.Writer) *Logger {
	logger := &Logger{out: out}
	logger.init()
	return logger
}

func (this *Logger) init() {
	this.level = WARNING
	this.stackTraceLevel = EMERGENCY
	this.showShorFile = true
	this.shorFileNameDepth = 1
	this.showPackage = true
	this.showFuncName = false
}

func (this *Logger) SetLevel(level Level) {
	this.level = level
}

func (this *Logger) SetStackTraceLevel(level Level) {
	this.stackTraceLevel = level
}

func (this *Logger) log(level Level, msg string, args ...interface{}) {
	if level > this.level {
		return
	}

	// Write all the data into a buffer.
	// Format is:
	// <timestamp> [level][<file>:<line>:<function>]: <message>
	now := time.Now()

	var buffer bytes.Buffer
	buffer.WriteString(now.Format(TIMESTAMP_FMT))
	buffer.WriteString(" ")
	buffer.WriteString(fmt.Sprintf("[%s]", level.String()))
	buffer.WriteString(this.fileInfo(4))
	buffer.WriteString(": ")
	buffer.WriteString(fmt.Sprintf(msg, args...))
	buffer.WriteString("\n")

	if level <= this.stackTraceLevel {
		buffer.WriteString("--- BEGIN stacktrace: ---\n")
		buffer.Write(stackTrace())
		buffer.WriteString("--- END stacktrace ---\n\n")
	}

	this.output(buffer.Bytes())
}

func (this *Logger) print(msg string, args ...interface{}) {
	var buffer bytes.Buffer
	buffer.WriteString(fmt.Sprintf(msg, args...))
	buffer.WriteString("\n")
	this.output(buffer.Bytes())
}

func (this *Logger) output(msg []byte) {
	this.mu.Lock()
	defer this.mu.Unlock()
	this.out.Write(msg)
}

func (this *Logger) fileInfo(depth int) string {
	stackInfo := "[???]"
	if pc, fileName, line, ok := runtime.Caller(depth); ok {

		if this.showShorFile {
			fileName = extractFileName(fileName, this.shorFileNameDepth)
		}

		if this.showFuncName {
			funcName := runtime.FuncForPC(pc).Name()
			if !this.showPackage {
				funcName = path.Base(funcName)
			}
			stackInfo = fmt.Sprintf("[%s:%d:%s]", fileName, line, funcName)
		} else {
			stackInfo = fmt.Sprintf("[%s:%d]", fileName, line)
		}

	}
	return stackInfo
}

func stackTrace() []byte {
	trace := make([]byte, STACK_BUFFER_SIZE)
	count := runtime.Stack(trace, true)
	return trace[:count]
}

func extractFileName(fileName string, shorFileNameDepth int) string {
	if shorFileNameDepth < 1 {
		shorFileNameDepth = 1
	}

	i := 0
	depth := 0
	for i = len(fileName) - 1; i >= 0; i-- {
		if fileName[i] == '/' {
			depth++
			if depth >= shorFileNameDepth {
				break
			}
		}
	}
	return fileName[i+1 : len(fileName)]
}

func (this *Logger) Emergency(msg string, args ...interface{}) {
	this.log(EMERGENCY, msg, args...)
}

func (this *Logger) Alert(msg string, args ...interface{}) {
	this.log(ALERT, msg, args...)
}

func (this *Logger) Critical(msg string, args ...interface{}) {
	this.log(CRITICAL, msg, args...)
}

func (this *Logger) Error(msg string, args ...interface{}) {
	this.log(ERROR, msg, args...)
}

func (this *Logger) Warning(msg string, args ...interface{}) {
	this.log(WARNING, msg, args...)
}

func (this *Logger) Notice(msg string, args ...interface{}) {
	this.log(NOTICE, msg, args...)
}

func (this *Logger) Info(msg string, args ...interface{}) {
	this.log(INFO, msg, args...)
}

func (this *Logger) Debug(msg string, args ...interface{}) {
	this.log(DEBUG, msg, args...)
}

func (this *Logger) Print(msg string) {
	this.print(msg)
}

func (this *Logger) PrintStack() {
	this.output(stackTrace())
}

var defaultLogger *Logger = nil

func getDefaultLogger() *Logger {
	if defaultLogger != nil {
		return defaultLogger
	}
	defaultLogger = NewLogger(os.Stderr)
	return defaultLogger
}

func Emergency(msg string, args ...interface{}) {
	getDefaultLogger().Emergency(msg, args...)
}

func Alert(msg string, args ...interface{}) {
	getDefaultLogger().Alert(msg, args...)
}

func Critical(msg string, args ...interface{}) {
	getDefaultLogger().Critical(msg, args...)
}

func Error(msg string, args ...interface{}) {
	getDefaultLogger().Error(msg, args...)
}

func Warning(msg string, args ...interface{}) {
	getDefaultLogger().Warning(msg, args...)
}

func Notice(msg string, args ...interface{}) {
	getDefaultLogger().Notice(msg, args...)
}

func Info(msg string, args ...interface{}) {
	getDefaultLogger().Info(msg, args...)
}

func Debug(msg string, args ...interface{}) {
	getDefaultLogger().Debug(msg, args...)
}

func Print(msg string) {
	getDefaultLogger().Print(msg)
}

func PrintStack() {
	getDefaultLogger().PrintStack()
}

func SetLevel(level Level) {
	getDefaultLogger().SetLevel(level)
}

func SetStackTraceLevel(level Level) {
	getDefaultLogger().SetStackTraceLevel(level)
}

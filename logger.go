package logger

import (
	"fmt"
	"log"
	"net"
	"runtime"
	"strings"
	"time"
)

type Logger interface {
	Log(level string, fmt string, v ...interface{})
	LogE(fmt string, v ...interface{})
	LogI(fmt string, v ...interface{})
	LogW(fmt string, v ...interface{})
}

type logger struct {
	addr    *net.UDPAddr
	conn    *net.UDPConn
	console bool
	tag     string
	level   int
}

const (
	INFO  = 3
	WARN  = 4
	ERROR = 5
)

func levelStr2N(str string) int {
	rv := -1
	switch {
	case strings.EqualFold(str, "INFO"):
		rv = INFO
	case strings.EqualFold(str, "WARN"):
		rv = WARN
	case strings.EqualFold(str, "ERROR"):
		rv = ERROR
	}
	return rv
}

func NewLoggerEx(tag string, addr string, console bool, level string) Logger {
	ret := new(logger)
	ret.console = console
	ret.tag = tag
	if addr == "" {
		return ret
	}

	ret.level = levelStr2N(level)
	if ret.level == -1 {
		log.Printf("unkown level(%s)\n", level)
		return nil
	}

	serverAddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		log.Println("can not resolve udp address", addr, err.Error())
		return nil
	}

	ret.addr = serverAddr
	ret.conn, err = net.DialUDP("udp", nil, ret.addr)
	if err != nil {
		log.Println("can not dial udp to", addr, err.Error())
		return nil
	}
	return ret
}

func NewLogger(tag string, addr string, console bool, level string) Logger {
	return NewLoggerEx(tag, addr, console, "INFO")
}

func (this logger) Log(level string, f string, v ...interface{}) {
	if !this.console && this.addr == nil {
		return
	}

	curLevel := levelStr2N(level)
	if curLevel < this.level {
		return
	}

	now := time.Now()
	msg := fmt.Sprintf(f, v...)
	//log formate [APP name] [Log time] [Log message]
	head := fmt.Sprintf("[%s] [%s] [%s] ",
		this.tag,
		now.Format("2006-01-02 15:04:05.000"),
		level)
	var logPos string
	if strings.ToLower(level) == "error" || strings.ToLower(level) == "warn" {
		pc, file, line, _ := runtime.Caller(3)
		funcInfo := runtime.FuncForPC(pc)
		logPos = fmt.Sprintf(" %s@%s:%d ", funcInfo.Name(), file, line)
	}
	if this.console {
		log.Println(head, msg, logPos)
	}

	if this.addr != nil {
		go func() {
			data := []byte(head)
			data = append(data, []byte(msg+logPos)...)
			this.conn.Write(data)
		}()
	}
}

func (this logger) LogE(fmt string, v ...interface{}) {
	this.Log("ERROR", fmt, v...)
}

func (this logger) LogW(fmt string, v ...interface{}) {
	this.Log("WARN", fmt, v...)
}

func (this logger) LogI(fmt string, v ...interface{}) {
	this.Log("INFO", fmt, v...)
}

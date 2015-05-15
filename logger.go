package logger

import (
	"fmt"
	"log"
	"net"
	"runtime"
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
}

func NewLogger(tag string, addr string, console bool) Logger {
	ret := new(logger)
	ret.console = console
	ret.tag = tag
	if addr == "" {
		return ret
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

func (this logger) Log(level string, f string, v ...interface{}) {
	if !this.console && this.addr == nil {
		return
	}
	now := time.Now()
	msg := fmt.Sprintf(f, v...)
	//log formate [APP name] [Log time] [Log message]
	head := fmt.Sprintf("[%s] [%s] [%s] ",
		this.tag,
		now.Format("2006-01-02 15:04:05.000"),
		level)

	pc, file, line, _ := runtime.Caller(1)
	funcInfo := runtime.FuncForPC(pc)
	logPos := fmt.Sprintf("FUNC:%s@LINE:%d@FILE:%s", funcInfo.Name(), line, file)

	if this.console {
		log.Println(head, logPos, msg)
	}

	if this.addr != nil {
		go func() {
			data := []byte(head)
			data = append(data, []byte(logPos+msg)...)
			//data = append(data, )
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

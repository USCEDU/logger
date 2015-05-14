package logger

import (
	"testing"
)

func TestNewLogger(t *testing.T){
	l := NewLogger("test", "192.168.5.70:3232", true)
	if l == nil{
		t.Error("construct logger failed")
	}else{
		t.Log("create logger success")
	}
}

func TestHelloWorld(t *testing.T){
	l := NewLogger("test", "", true)
	l.LogI("Hello, %s", "world")
}

func TestUDPHelloWorld(t *testing.T){
	l := NewLogger("test", "192.168.6.40:9090", true)
	l.LogE("Hello, it is a udp message")
}

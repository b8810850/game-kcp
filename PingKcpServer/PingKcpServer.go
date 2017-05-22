package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"strconv"
	"time"

	"github.com/xtaci/kcp-go"
	//	"strconv"
)

type Config struct {
	Port string `xml:"port"`
}

func main() {

	content, err := ioutil.ReadFile("serverconf.xml")
	if err != nil {
		log.Fatal(err)
	}
	var result Config
	err = xml.Unmarshal(content, &result)
	if err != nil {
		log.Fatal(err)
	}

	l := kcp.ListenKcp(result.Port)
	l.SetFastMode()
	go func() {
		kcplistener := l
		kcplistener.SetReadBuffer(4 * 1024 * 1024)
		kcplistener.SetWriteBuffer(4 * 1024 * 1024)
		kcplistener.SetDSCP(46)
		for {
			s, err := l.Accept()
			if err != nil {
				fmt.Println(err.Error())
				return
			}

			// coverage test
			s.(*kcp.UDPSession).SetReadBuffer(4 * 1024 * 1024)
			s.(*kcp.UDPSession).SetWriteBuffer(4 * 1024 * 1024)

			go handleEcho(s.(*kcp.UDPSession))
		}
	}()
	for {
		time.Sleep(time.Hour * 24)
	}
}

func handleEcho(conn *kcp.UDPSession) {
	conn.SetStreamMode(true)
	conn.SetWindowSize(4096, 4096)
	conn.SetNoDelay(1, 10, 2, 1)
	conn.SetDSCP(46)
	conn.SetMtu(1400)
	conn.SetACKNoDelay(false)
	conn.SetReadDeadline(time.Now().Add(time.Hour))
	conn.SetWriteDeadline(time.Now().Add(time.Hour))
	buf := make([]byte, 65536)
	defer func() {
		conn.Close()
		if err := recover(); err != nil {
			fmt.Println(err)
		}
	}()
	for {

		n, err := conn.Read(buf)
		if err != nil {
			fmt.Println("n=" + strconv.Itoa(n))
			panic(err)
		}
		n, err = conn.Write(buf[:n])
		if nil != err {
			panic(err)
		}
	}
}
func CheckError(err error, info string) (res bool) {
	if err != nil {
		fmt.Println("start server fail")
		fmt.Println(info + err.Error())
		return false
	}
	return true
}

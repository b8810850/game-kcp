package main

import (
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net"
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
	fmt.Println(result.Port)
	StartTpcServer(result.Port)
}

func CheckError(err error, info string) (res bool) {
	if err != nil {
		fmt.Println("start server fail")
		fmt.Println(info + err.Error())
		return false
	}
	return true
}

func StartTpcServer(port string) {
	service := ":" + port

	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	if !CheckError(err, "ResolveTCPAddr") {
		return
	}

	l, err := net.ListenTCP("tcp", tcpAddr)
	if !CheckError(err, "ListenTCP") {
		return
	}

	for {
		fmt.Println("Listening...")
		conn, err := l.Accept()
		if !CheckError(err, "Accept") {
			continue
		}
		fmt.Println("Accepting...")

		go Handler(conn)
	}
}

func Handler(conn net.Conn) {
	buf := make([]byte, 4096)
	for {
		readLenth, err1 := conn.Read(buf)
		if nil != err1 {
			fmt.Println(err1.Error())
			return
		}
		_, err := conn.Write(buf[0:readLenth])

		//		fmt.Println("i receive message from client:" + string(buf[0:readLenth]))
		if nil != err {
			fmt.Println(err.Error())
			return
		}
		//		fmt.Println("i send message to client, len:" + strconv.Itoa(lenth))
	}
}

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
	StartUdpServer(result.Port)
}

func CheckError(err error, info string) (res bool) {
	if err != nil {
		fmt.Println("start server fail")
		fmt.Println(info + err.Error())
		return false
	}
	return true
}

func StartUdpServer(port string) {
	service := ":" + port

	udpAddr, err := net.ResolveUDPAddr("udp", service)
	if !CheckError(err, "ResolveUDPAddr") {
		return
	}

	l, err := net.ListenUDP("udp", udpAddr)
	if !CheckError(err, "ListenUDP") {
		return
	}
	fmt.Println("listening")
	Handler(l)
}

func Handler(conn *net.UDPConn) {
	buf := make([]byte, 2048)
	for {

		readLenth, raddr, err1 := conn.ReadFromUDP(buf)
		if nil != err1 {
			fmt.Println(err1.Error())
			return
		}
		_, err := conn.WriteToUDP(buf[0:readLenth], raddr)

		//		fmt.Println("i receive message from client:" + string(buf[0:readLenth]))
		if nil != err {
			fmt.Println(err.Error())
			continue
		}
		//		fmt.Println("i send message to client, len:" + strconv.Itoa(lenth))
	}
}

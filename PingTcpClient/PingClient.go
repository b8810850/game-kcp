package main

import (
	"bytes"
	"encoding/binary"
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net"
	"os"
	"strconv"
	"strings"
	"time"
	"unsafe"
)

type Config struct {
	Port string `xml:"port"`
	Freq string `xml:"freq"`
}

var pingFreq int

func main() {
	fmt.Println("start")
	content, err := ioutil.ReadFile("clientconf.xml")
	if err != nil {
		log.Fatal(err)
	}
	var result Config
	err = xml.Unmarshal(content, &result)
	if err != nil {
		log.Fatal(err)
	}
	pingFreq, err = strconv.Atoi(result.Freq)
	if nil != err {
		log.Fatal(err)
	}
	StartTcpClient(result.Port)
}

func PathExists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return false, err
}

func CheckError(err error, info string) (res bool) {
	if err != nil {
		fmt.Println("start client fail")
		fmt.Println(info + err.Error())
		return false
	}
	return true
}

type myPkg struct {
	id     int
	myTime time.Time
}

func StartTcpClient(tcpArrd string) {

	addr, err := net.ResolveTCPAddr("tcp", tcpArrd)
	if !CheckError(err, "ResolveTCPAddr") {
		return
	}
	conn, err := net.DialTCP("tcp", nil, addr)
	if !CheckError(err, "DialUDP") {
		return
	}
	pingPkg := make(chan myPkg, 1)
	go Ping(conn, pingPkg)
	recvPkg := make(chan myPkg, 1)
	go Recv(conn, recvPkg)
	pingMap := make(map[int]time.Time)
	recvPkgBuffer := make(map[int]myPkg)
	sendPkgNum := 0
	recvPkgNum := 0

	result, err := PathExists("./log")
	if !result {
		err := os.Mkdir("./Log", os.ModeDir) //在当前目录下生成log目录
		if nil != err {
			fmt.Println(err.Error())
			os.Exit(0)
		}
	}

	now := time.Now().String()
	now = strings.Replace(now, ":", "_", -1)
	f, err := os.Create("./Log/" + now) //创建文件
	if nil != err {
		fmt.Println(err.Error())
		os.Exit(0)
	}

	defer f.Close()
	var writer io.Writer = f

	var totalDelay time.Duration = 0
	for {
		select {
		case e1 := <-pingPkg:
			//如果ch1通道成功读取数据，则执行该case处理语句
			pkgId := e1.id
			pingTime := e1.myTime
			pingMap[pkgId] = pingTime
			sendPkgNum += 1
			break
		case e2 := <-recvPkg:
			if 0 != len(recvPkgBuffer) {
				for _, pkg := range recvPkgBuffer {
					pkgId := pkg.id
					var receiveTime time.Time = pkg.myTime
					pingTime, ok := pingMap[pkgId]
					if ok {
						fmt.Println(len(recvPkgBuffer))
						delay := receiveTime.Sub(pingTime)

						recvPkgNum += 1
						totalDelay += (delay)
						fmt.Println("receive pkg ,delay :" + delay.String() + " 发包数:" + strconv.Itoa(sendPkgNum) + " 接包数:" + strconv.Itoa(recvPkgNum) + " 丢包率:" + strconv.FormatFloat(float64(sendPkgNum-recvPkgNum)/float64(sendPkgNum), 'f', -1, 32) + " avg delay:" + time.Duration(float64(totalDelay)/float64(recvPkgNum)).String())
						writer.Write([]byte("receive pkg ,delay :" + delay.String() + " 发包数:" + strconv.Itoa(sendPkgNum) + " 接包数:" + strconv.Itoa(recvPkgNum) + " 丢包率:" + strconv.FormatFloat(float64(sendPkgNum-recvPkgNum)/float64(sendPkgNum), 'f', -1, 32) + "avg delay:" + time.Duration(float64(totalDelay)/float64(recvPkgNum)).String() + "\n"))

						delete(recvPkgBuffer, pkgId)
						delete(pingMap, pkgId)
					}
				}
			}

			//如果ch2通道成功读取数据，则执行该case处理语句
			pkgId := e2.id
			var receiveTime time.Time = e2.myTime
			pingTime, ok := pingMap[pkgId]
			if !ok {
				recvPkgBuffer[pkgId] = e2
				break
			}
			delay := receiveTime.Sub(pingTime)
			recvPkgNum += 1
			totalDelay += (delay)
			fmt.Println("receive pkg ,delay :" + delay.String() + " 发包数:" + strconv.Itoa(sendPkgNum) + " 接包数:" + strconv.Itoa(recvPkgNum) + " 丢包率:" + strconv.FormatFloat(float64(sendPkgNum-recvPkgNum)/float64(sendPkgNum), 'f', -1, 32) + " avg delay:" + time.Duration(float64(totalDelay)/float64(recvPkgNum)).String())
			writer.Write([]byte("receive pkg ,delay :" + delay.String() + " 发包数:" + strconv.Itoa(sendPkgNum) + " 接包数:" + strconv.Itoa(recvPkgNum) + " 丢包率:" + strconv.FormatFloat(float64(sendPkgNum-recvPkgNum)/float64(sendPkgNum), 'f', -1, 32) + "avg delay:" + time.Duration(float64(totalDelay)/float64(recvPkgNum)).String() + "\n"))
			delete(pingMap, pkgId)
			break
		}
	}
}
func Recv(conn *net.TCPConn, pkg chan myPkg) {
	buf := make([]byte, 4096)
	for {
		lenth, readError := conn.Read(buf)
		if nil != readError {
			fmt.Println(readError.Error())
			fmt.Println("Read server go wrong ,exit.")
			return
		}
		receiveTime := time.Now()

		var msgLen int32
		if lenth < int(unsafe.Sizeof(msgLen)) {
			fmt.Println("parse msg go wrong")
			continue
		}
		parseIndex := 0
		for {
			buffer := bytes.NewBuffer(buf[parseIndex:])
			binary.Read(buffer, binary.BigEndian, &msgLen)
			diffIndex := int(unsafe.Sizeof(msgLen))
			parseIndex += diffIndex
			if parseIndex+int(msgLen) > lenth {
				fmt.Println("解包错误")
				break
			}
			pingNum, err := strconv.Atoi(string(buf[parseIndex : parseIndex+int(msgLen)]))
			if nil != err {
				fmt.Println("parse server go wrong .")
				break
			}
			pkg <- myPkg{pingNum, receiveTime}
			parseIndex += int(msgLen)
			if parseIndex >= lenth {
				break
			}
		}
	}
}
func Ping(conn *net.TCPConn, pkg chan myPkg) {
	pingNum := 0
	for {
		pingNum += 1

		var msgLenth int32 = int32(len([]byte(strconv.Itoa(pingNum))))

		b_buf := bytes.NewBuffer([]byte{})
		binary.Write(b_buf, binary.BigEndian, msgLenth)

		buf := append(b_buf.Bytes(), []byte(strconv.Itoa(pingNum))...)
		_, err := conn.Write(buf)
		if nil != err {
			fmt.Println(err.Error())
			conn.Close()
			break
		}
		cur := time.Now()
		pkg <- myPkg{pingNum, cur}
		time.Sleep(time.Duration(pingFreq) * time.Millisecond)
	}
}

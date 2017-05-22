package main

import (
	"encoding/xml"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/xtaci/kcp-go"
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

	fmt.Println(result.Port)
	StartKcpClient(result.Port)
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

func StartKcpClient(tcpArrd string) {
	sess, err := kcp.DiagKcp(tcpArrd)
	if !CheckError(err, "DialUDP") {
		return
	}

	sess.SetNoDelay(1, 10, 2, 1)
	sess.SetACKNoDelay(false)
	sess.SetWriteDelay(false)

	pingPkg := make(chan myPkg, 1)
	go Ping(sess, pingPkg)
	recvPkg := make(chan myPkg, 1)
	go Recv(sess, recvPkg)
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
func Recv(sess *kcp.UDPSession, pkg chan myPkg) {
	buf := make([]byte, 4096)
	for {
		lenth, readError := sess.Read(buf)
		if nil != readError {
			sess.Close()
			fmt.Println(readError.Error())
			fmt.Println("Read server go wrong ,exit.")
			return
		}
		receiveTime := time.Now()
		pingNum, err := strconv.Atoi(string(buf[0:lenth]))
		if nil != err {
			fmt.Println("parse server go wrong ,exit.")
			return
		}
		pkg <- myPkg{pingNum, receiveTime}
	}
}
func Ping(sess *kcp.UDPSession, pkg chan myPkg) {
	pingNum := 0
	for {
		pingNum += 1

		_, err := sess.Write([]byte(strconv.Itoa(pingNum)))
		cur := time.Now()
		pkg <- myPkg{pingNum, cur}
		if nil != err {
			fmt.Println(err.Error())
			sess.Close()
			break
		}
		time.Sleep(time.Duration(pingFreq) * time.Millisecond)
	}
}

package kcp

import (
	"flag"
	"fmt"
	"net"
	"os"
	"runtime"
	runtimePprof "runtime/pprof"

	"crypto/sha1"

	"golang.org/x/crypto/pbkdf2"
)

var flagCpuprofile = flag.String("cpuprofile", "cpuProfile.prof", "write cpu profile to file")

const portGame = "127.0.0.1:20000"

var key = []byte("testkey")
var fec = 4
var pass = pbkdf2.Key(key, []byte(portGame), 4096, 32, sha1.New)

func DiagKcp(port string) (*UDPSession, error) {
	//block, _ := NewSalsa20BlockCrypt(pass)
	sess, err := DialWithOptions(port, nil, 10, 3)
	if err != nil {
		panic(err)
	}

	sess.SetStreamMode(false)
	sess.SetWindowSize(4096, 4096)
	sess.SetReadBuffer(4 * 1024 * 1024)
	sess.SetWriteBuffer(4 * 1024 * 1024)
	sess.SetNoDelay(1, 10, 2, 1)
	sess.SetMtu(1400)
	sess.SetACKNoDelay(false)
	sess.SetWriteDelay(false)
	sess.SetDUP(1)
	//sess.SetDeadline(time.Now().Add(time.Minute))
	return sess, err
}

func ListenKcp(port string) *Listener {
	runtime.GOMAXPROCS(8)
	if *flagCpuprofile != "" {

		f, _ := os.Create(*flagCpuprofile)

		runtimePprof.StartCPUProfile(f)

		defer runtimePprof.StopCPUProfile()

	}

	//	block, _ := NewSalsa20BlockCrypt(pass)
	l, err := ListenWithOptions(port, nil, 10, 3)
	if err != nil {
		panic(err)
	}
	listener := l
	fmt.Println("listen start...")
	return listener
}

func listenGameForTest() (net.Listener, error) {
	block, _ := NewSalsa20BlockCrypt(pass)
	return ListenWithOptions(portGame, block, 10, 3)
}

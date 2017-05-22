package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
)

func appendToFile(fileName string, content string) error {
	// 以只写的模式，打开文件
	f, err := os.OpenFile(fileName, os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("cacheFileList.yml file create failed. err: " + err.Error())
	} else {
		// 查找文件末尾的偏移量
		n, _ := f.Seek(0, os.SEEK_END)
		// 从末尾的偏移量开始写入内容
		_, err = f.WriteAt([]byte(content), n)
	}
	defer f.Close()
	return err
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

func main() {

	// flag.Args方式
	flag.Parse()
	var ch []string = flag.Args()
	if ch != nil && len(ch) > 0 {
		fmt.Println("Hello ", ch[0]) // 第一个参数开始
	}
	isExit, _ := PathExists("statistic")
	if false == isExit {
		os.Create("statistic")
	}
	dir, error := os.OpenFile("../log", os.O_RDONLY, os.ModeDir)
	if error != nil {
		defer dir.Close()
		fmt.Println(error.Error())
		return
	}
	fileinfo, _ := dir.Stat()
	fmt.Println(fileinfo.IsDir())
	names, _ := dir.Readdir(-1)
	for _, name := range names {
		file, _ := os.OpenFile("../log/"+name.Name(), os.O_RDONLY, os.ModeDevice)
		bfRd := bufio.NewReader(file)
		var lastStr string
		for {
			line, err := bfRd.ReadString('\n')

			if err != nil { //遇到任何错误立即返回，并忽略 EOF 错误信息
				if err == io.EOF {
					appendToFile("statistic", ch[0]+" "+lastStr)
					return
				}
			}
			lastStr = line

		}

	}
}

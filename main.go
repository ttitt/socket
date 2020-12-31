package main

import (
	"fmt"
	"log"
	"net"
	"strconv"
)

const (
	CONSTHEADER       = "H"
	CONSTHEADERLENGTH = len(CONSTHEADER)
	CONSTMSGLENGTH    = 4
)

func main() {
	Listening()
}

func Listening() {
	tcpListen, err := net.Listen("tcp", ":2333")

	if err != nil {
		panic(err)
	}

	for {
		conn, err := tcpListen.Accept()
		if err != nil {
			log.Println(err)
			continue
		}
		fmt.Println(conn.RemoteAddr())
		fmt.Println("已经连接")
		go connHandle(conn)
	}
}
func connHandle(conn net.Conn) {
	defer fmt.Println(conn.RemoteAddr())
	defer fmt.Println("关闭了连接")
	defer conn.Close()

	var errs error
	receive := true
	readBuff := make([]byte, 20)
	tempBuff := make([]byte, 0)
	data := make([]byte, 20)
	log.Println(string(tempBuff))
	for {
		if receive {
			n, err := conn.Read(readBuff)
			if err != nil {
				return
			}
			tempBuff = append(tempBuff, readBuff[:n]...)
		}

		tempBuff, data, errs, receive = Depack(tempBuff) //对缓冲区的包进行分包处理
		if errs != nil {
			return
		}
		if len(data) == 0 {
			continue
		}
		fmt.Println(string(data))
	}
}

//拆包
func Depack(buff [] byte) ([]byte, []byte, error, bool) {
	bufflen := len(buff)
	log.Println(string(buff))
	//头部长度不完整
	if bufflen < CONSTHEADERLENGTH+CONSTMSGLENGTH {
		log.Println("数据包头部长度不完整")
		return buff, nil, nil, true
	}

	//是不是有头
	if string(buff[:CONSTHEADERLENGTH]) != CONSTHEADER {
		log.Println("数据包并不是头")
		return buff, nil, nil, true
	}

	//内容不完整
	msgLength, _ := strconv.Atoi(string(buff[CONSTHEADERLENGTH : CONSTHEADERLENGTH+CONSTMSGLENGTH]))
	if bufflen < CONSTHEADERLENGTH+CONSTMSGLENGTH+msgLength {
		log.Println("内容不完整")
		return buff, nil, nil, true
	}

	data := buff[CONSTHEADERLENGTH+CONSTMSGLENGTH : CONSTHEADERLENGTH+CONSTMSGLENGTH+msgLength]
	buffs := buff[CONSTHEADERLENGTH+CONSTMSGLENGTH+msgLength:]
	if len(buffs) > CONSTHEADERLENGTH+CONSTMSGLENGTH {
		return buffs, data, nil, false
	} else {
		return buffs, data, nil, true
	}

}

package main

import (
	"./resp"
	"crypto/sha256"
	"crypto/tls"
	"encoding/base64"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strings"
)

const (
	Ping       string = "-PING Empty request."
	NoAuth     string = "-NOAUTH Authentication required."
	BadRequest string = "-ERR there was an error processing the request"
	BadAuth    string = "-ERR invalid password"
	NoServer   string = "-ERR Server unavailable"
)

func abort(conn net.Conn, msg string) {
	log.Printf("Connection failed: %s", msg)
	conn.Write([]byte(msg + "\r\n"))
	conn.Close()
}

func auth(conn net.Conn, secret string) {
	rd := resp.NewReader(conn)
	v, _, err := rd.ReadValue()

	if err == io.EOF {
		abort(conn, Ping)
		return
	}

	if err != nil {
		log.Print("err != nil")
		log.Printf("err: %s", err)
		abort(conn, BadRequest)
		return
	}

	if v.Type() != resp.Array {
		log.Print("v.Type() != resp.Array ")
		log.Printf("v.Type(): %s", v.Type())
		log.Printf("v: %s", v)
		abort(conn, NoAuth)
		return
	}

	cmdParts := v.Array()
	if len(cmdParts) != 2 {
		log.Print("cmdParts != 2")
		log.Printf("cmdParts: %s", cmdParts)
		abort(conn, NoAuth)
		return
	}

	if strings.ToUpper(cmdParts[0].String()) != "AUTH" {
		log.Print("cmdParts[0] != AUTH")
		abort(conn, NoAuth)
		return
	}

	auth, err := base64.StdEncoding.DecodeString(cmdParts[1].String())
	if err != nil {
		abort(conn, BadAuth)
		return
	}

	authSplit := strings.Split(string(auth), " ")
	if len(authSplit) != 2 {
		abort(conn, BadAuth)
		return
	}
	hostPort := authSplit[0]
	signature := authSplit[1]

	hash := sha256.Sum256([]byte(hostPort + secret))
	hexdigest := fmt.Sprintf("%x", hash)
	if hexdigest != signature {
		abort(conn, BadAuth)
		return
	}

	log.Printf("Authentication succeeded, connecting client to %s", hostPort)

	// Flush the parsing buffer to the client, in case we ate too much
	trailing, _ := rd.Peek(rd.Buffered())
	conn.Write(trailing)

	forward(conn, hostPort)
}

func forward(conn net.Conn, hostPort string) {
	client, err := net.Dial("tcp", hostPort)
	if err != nil {
		abort(conn, NoServer)
		return
	}

	log.Printf("Connected to backend %s\n", hostPort)
	conn.Write([]byte("+OK\r\n"))

	go func() {
		defer client.Close()
		defer conn.Close()
		io.Copy(client, conn)
	}()
	go func() {
		defer client.Close()
		defer conn.Close()
		io.Copy(conn, client)
	}()
}

func main() {
	if len(os.Args) < 3 {
		log.Fatalf("Usage %s listen:port secret [cert_file] [cert_key_file]\n", os.Args[0])
		return
	}

	var listener net.Listener
	var err error
	if len(os.Args) == 5 {
		cer, err := tls.LoadX509KeyPair(os.Args[3], os.Args[4])
		if err != nil {
			log.Fatalf("Failed to load certificates: %v", err)
			return
		}

		config := &tls.Config{Certificates: []tls.Certificate{cer}}
		listener, err = tls.Listen("tcp", os.Args[1], config)
	} else {
		listener, err = net.Listen("tcp", os.Args[1])
	}

	if err != nil {
		log.Fatalf("Failed to setup listener: %v", err)
	}

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("ERROR: failed to accept listener: %v", err)
		}
		go auth(conn, os.Args[2])
	}
}

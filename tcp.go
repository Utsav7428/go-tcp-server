package main

import (
	"bufio"
    "fmt"
	"io"
	"log"
	"net"
	"strings"
	"time"
    "regexp"
	"strconv"
)
const MaxLineLenBytes = 1024
const ReadWriteTimeout = time.Minute

func main() {
    //create a server
	lis, err := net.Listen("tcp", ":1234")
	if err != nil {
		log.Fatalf("failed to start listener: %v", err)
	}

	for {
        //accepting connections
		conn, err := lis.Accept()
		if err != nil {
			log.Printf("failed to accept conn: %v", err)
			continue
		}
     //once a connection is accepted spawn a goroutine to accept it
		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {

    log.Printf("accepted connection from %s", conn.RemoteAddr())

    //Any IO error from a socket means that the connection is unstable, and should be closed.
    defer func() {
		_ = conn.Close()
		log.Printf("closed connection from %s", conn.RemoteAddr())
	}()

	done := make(chan struct{})

    // time out one minute from now if no
	// data is received. the error can be
	// safely ignored.
	_ = conn.SetReadDeadline(time.Now().Add(ReadWriteTimeout))

    go func() {
        // limit the maximum line length (in bytes)
        lim := &io.LimitedReader{
            R: conn,
            N: MaxLineLenBytes,
        }
        scan := bufio.NewScanner(lim)
        for scan.Scan() {
            input := strings.TrimSpace(scan.Text())

            if input == "" {
				continue
			}

            output := evalSimpleExpression(input)
            if _, err := conn.Write([]byte(output+"\n")); err != nil {
                log.Printf("failed to write output: %v", err)
                return
            }
            log.Printf("wrote response: %s", output)
            // reset the number of bytes remaining in the LimitReader
            lim.N = MaxLineLenBytes
            
            _ = conn.SetReadDeadline(time.Now().Add(ReadWriteTimeout))
        }
    
        done <- struct{}{}
    }()

	<-done
}

func evalSimpleExpression(input string) string {
	// Use regex to extract parts: number, operator, number
	re := regexp.MustCompile(`^\s*(\d+)\s*([+\-*\/])\s*(\d+)\s*$`)
	matches := re.FindStringSubmatch(input)
	if len(matches) != 4 {
		return "invalid expression"
	}

	a, _ := strconv.Atoi(matches[1])
	op := matches[2]
	b, _ := strconv.Atoi(matches[3])

	switch op {
	case "+":
		return fmt.Sprintf("%d", a+b)
	case "*":
		return fmt.Sprintf("%d", a*b)
	case "-":
		return fmt.Sprintf("%d", a-b)
	case "/":
		if b == 0 {
			return "division by zero"
		}
		return fmt.Sprintf("%d", a/b)
	default:
		return "unsupported operator"
	}
}


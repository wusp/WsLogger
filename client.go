package main

import (
	"context"
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var (
	wg sync.WaitGroup
)

type Client struct {
	// The websocket connection.
	conn *websocket.Conn
	// The raw comm saved file
	file *os.File
	// channel used to send heartbeat
	heartCh chan int
}

func (c *Client) saveMessageToFile() {
	// send signal to waitGroup at the end of this function
	defer func() {
		time.Sleep(500 * time.Millisecond)
		wg.Done()
	}()

	log.Println("Start loop messages...")
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			//errStr := fmt.Sprintln("Prepare to exit due to", err)
			//log.Print(errStr)
			break
		} else {
			_, err = c.file.WriteString(fmt.Sprintf("%s %s\n", time.Now().Format("2006-01-02-15:04:05"), string(message)))
			if err != nil {
				errStr := fmt.Sprintln("Write string to file error.", err)
				log.Print(errStr)
				break
			}
			// send heartbeat
			c.heartCh <- 1
		}
	}
}

func openFile(dir string, fileName string) *os.File {
	log.Println("Create file with name:", fileName, "\n on dir:", dir)
	filePathStr := filepath.Join(dir, filepath.Base(fileName))
	file, err := os.Create(filePathStr)
	if err != nil {
		return nil
	}
	return file
}

func connectWs(urlStr string, logDir string, readTimeout int) {
	log.Println("Start websocket message loop-reading on", urlStr, "readTimeout(s):", readTimeout)
	heartbeatCh := make(chan int, 100)

	conn, _, err := websocket.DefaultDialer.Dial(urlStr, nil)
	if err != nil {
		log.Print("dial err:", err)
		return
	}

	client := &Client{
		conn:    conn,
		file:    openFile(logDir, fmt.Sprintf("rawcomm-%s.log", time.Now().Format("2006-01-02-15:04:05"))),
		heartCh: heartbeatCh}

	// 资源清理Callback
	closeResourceF := func() {
		// close connection when some errors happening...
		log.Println("Close heartbeat channel")
		close(client.heartCh)
		log.Println("Close log file")
		_ = client.file.Close()
		log.Println("Close ws connection")
		_ = client.conn.Close()
		log.Println("All close!")
		fmt.Print("\n\n")
	}

	wg.Add(1)
	readTimeoutDuration := time.Duration(readTimeout)
	go heartbeatChecker(heartbeatCh, time.Second*readTimeoutDuration, closeResourceF)
	go client.saveMessageToFile()
	wg.Wait()
}

func heartbeatChecker(ch chan int, heartbeatTimeout time.Duration, closeResourceF func()) {
	defer closeResourceF()

	ctx, cancel := context.WithTimeout(context.Background(), heartbeatTimeout)

	for {
		select {
		case <-ch:
			cancel()
			ctx, cancel = context.WithTimeout(context.Background(), heartbeatTimeout)

		case <-ctx.Done():
			log.Println("Timeout triggered, program into panicking...")
			return
		}
	}
}

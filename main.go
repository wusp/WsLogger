package main

import (
	"flag"
	"fmt"
	"log"
	"net/url"
	"os"
	"path/filepath"
	"time"
)

var (
	runningDir string
)

func main() {
	runningDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Fatalln("some unknown error occurs:", err, "exit(1)")
		os.Exit(1)
	}

	var addr = flag.String("a", "localhost:8081", "WebSocket url")
	var path = flag.String("p", "/rawcomm", "RawComm url path")
	var dir = flag.String("d", fmt.Sprintf("%s%s", runningDir, "/log"), "Log's file storage，default is '$PWD'/log.\n")
	var restart = flag.Bool("c", false, "Restart a new task when connection errors happened.")
	var readOT = flag.Int("t", 60, "WebSocket read timeout(in seconds)")
	flag.Parse()

	// 1. create ws url
	u := url.URL{Scheme: "ws", Host: *addr, Path: *path}
	// 2. ensure dir exists
	if _, err := os.Stat(*dir); err == nil {
		// path/to/whatever exists
		fmt.Println("Log folder exists on:", *dir)
	} else if os.IsNotExist(err) {
		// path/to/whatever does *not* exist
		fmt.Println("Create log folder on:", *dir)
		e := os.MkdirAll(*dir, 0700)
		if e != nil {
			log.Fatal(e)
		}
	} else {
		// Schrodinger: file may or may not exist. See err for details.
		// Therefore, do *NOT* use !os.IsNotExist(err) to test for file existence
		fmt.Println("Unknown error occurs when checking log dir, exit(1)")
		os.Exit(1)
	}

	if *restart {
		// 自动重连WS的情况下，每次碰到错误退出则重新执行
		for {
			connectWs(u.String(), *dir, *readOT)
			time.Sleep(1 * time.Second)
		}
	} else {
		// 不重连的情况下，只要碰到错误退出就退出本程序
		connectWs(u.String(), *dir, *readOT)
	}

	log.Print("wsRawComm exit(0).")
}

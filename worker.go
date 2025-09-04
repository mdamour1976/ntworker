package main

import (
	"flag"
	"io"
	"log"
	"net"
	"os"
	"strconv"
	"time"
)

var (
	version       = "dev"
	commit        = "none"
	date          = "unknown"
	latestVersion string
)

func pollForUpdates(address string, interval time.Duration, returnImmediately bool) {
	for {
		func() {
			conn, err := net.Dial("tcp", address)
			if err != nil {
				log.Println(err)
				// log but don't prevent a failed connection break the app
				return
			}
			defer conn.Close()
			// read version
			version, err := io.ReadAll(conn)
			if err != nil {
				log.Println(err)
			}
			latestVersion = string(version)
			log.Println("Client received latest version:", latestVersion)
		}()
		if returnImmediately {
			break
		}
		time.Sleep(interval)
	}
}

func doSomething() {
	// endlessly pull work from an imaginary queue
	// terminate only when the version number changes
	for {
		if version != latestVersion {
			log.Println("New version available: " + latestVersion)
			os.Exit(0)
		}
		log.Println("Working...")
		time.Sleep(10 * time.Second)
		//os.Exit(1)
	}
}

func main() {
	log.Println("Program starting....")

	log.Printf("Version: %s Commit: %s Built: %s", version, commit, date)

	updateIntervalFlag := flag.Int("update-interval", 60, "Specify the update checking interval in minutes")
	portFlag := flag.Int("ipc-port", 9999, "Specify the update checking interval in hours")
	flag.Parse()
	flag.VisitAll(func(f *flag.Flag) {
		log.Printf("Flag: -%s=%s (default: %s)\n", f.Name, f.Value, f.DefValue)
	})

	address := "127.0.0.1:" + strconv.Itoa(*portFlag)

	pollForUpdates(address, time.Duration(*updateIntervalFlag)*time.Minute, true)

	go doSomething()

	pollForUpdates(address, time.Duration(*updateIntervalFlag)*time.Minute, false)

	// should be unreachable
	log.Println("Program terminating normally...")
}

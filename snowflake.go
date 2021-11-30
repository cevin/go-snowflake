package main

import (
	"flag"
	"fmt"
	"github.com/sony/sonyflake"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var (
	addr string
	port uint
	machineId uint

	quit = make(chan os.Signal)
)

func main() {

	flag.StringVar(&addr, "addr", "localhost", "listen address")
	flag.UintVar(&port, "port", 8080, "listen port")
	flag.UintVar(&machineId, "machine-id", 1, "machine id")

	flag.Parse()

	var uintMachineId = uint16(machineId)


	var snowflake = sonyflake.NewSonyflake(sonyflake.Settings{
		StartTime: time.Now(),
		MachineID: func() (uint16, error) {
			return uintMachineId, nil
		},
	})

	http.HandleFunc("/getId", func(writer http.ResponseWriter, request *http.Request) {

		writer.Header().Set("Content-Type", "application/json; charset=utf-8")

		if id, err := snowflake.NextID(); err == nil {
			writer.Write([]byte(fmt.Sprintf(`{"id":%d, machineId:%d}`, id, uintMachineId)))
		} else {
			writer.WriteHeader(http.StatusInternalServerError)
		}
	})

	log.Printf("PID: %d", os.Getpid())
	log.Printf("work with machine id: %d", uintMachineId)
	log.Printf("listen %s:%d", addr, port)
	log.Printf("visit /getId")


	signal.Notify(quit, syscall.SIGTERM, syscall.SIGQUIT, os.Interrupt)

	go func() {
		select {
		case <-quit:
			log.Printf("quit")
			os.Exit(0)
		}
	}()

	if err := http.ListenAndServe(fmt.Sprintf("%s:%d", addr, port), nil); err != nil {
		log.Fatal(err)
	}
}

package main

import (
	"lxm-oil-prices/service"
	"context"
	"flag"
	"os"
	"os/signal"
	"syscall"

	"fmt"
)


func main() {

	var port int
	flag.IntVar(&port, "p", 8000, "端口号")
	flag.Parse()

	sigs := make(chan os.Signal, 1)
	done := make(chan bool)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		service.Run(ctx, port)
	}()

	go func() {
		sig := <-sigs
		fmt.Println(sig)
		cancel()
		done <- true
	}()

	<-done
	fmt.Println("[IAM]=>stop service.")
}

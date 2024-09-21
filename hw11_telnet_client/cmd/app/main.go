package main

import (
	"context"
	"flag"
	"fmt"
	"net"
	"os"
	"os/signal"
	"time"

	"github.com/AnnDutova/otus_go_hw/hw11_telnet_client/internal/client/telnet"
)

var timeout time.Duration

func init() {
	flag.DurationVar(&timeout, "timeout", 10*time.Second, "telnet client default timeout")
}

func main() {
	flag.Parse()

	if len(os.Args[2:]) < 2 {
		fmt.Println("Invalid argument count")
		os.Exit(1)
	}

	addr := net.JoinHostPort(os.Args[2], os.Args[3])
	tc := telnet.NewTelnetClient(addr, timeout, os.Stdin, os.Stdout)

	if err := tc.Connect(); err != nil {
		fmt.Printf("Error occurred during creation Telnet client connection, err: %v", err)
		os.Exit(1)
	}
	defer tc.Close()

	ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer cancel()

	go func() {
		tc.Receive()
	}()

	go func() {
		tc.Send()
	}()

	<-ctx.Done()
}

package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

type (
	TelnetClient interface {
		Connect() error
		io.Closer
		Send() error
		Receive() error
	}

	telnetClt struct {
		conn    net.Conn
		addr    string
		timeout time.Duration
		in      io.ReadCloser
		out     io.Writer
	}
)

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &telnetClt{
		addr:    address,
		timeout: timeout,
		in:      in,
		out:     out,
	}
}

func (clt *telnetClt) Connect() error {
	conn, err := net.DialTimeout("tcp", clt.addr, clt.timeout)
	if err != nil {
		return err
	}
	clt.conn = conn

	fmt.Fprintf(os.Stderr, "...Connected to %s\n", clt.addr)

	return nil
}

func (clt *telnetClt) Close() error {
	fmt.Fprintf(os.Stderr, "...Connection was closed by peer\n")
	return clt.conn.Close()
}

func (clt *telnetClt) Send() error {
	if _, err := io.Copy(clt.conn, clt.in); err != nil {
		return fmt.Errorf("failed to send: %w", err)
	}
	return nil
}

func (clt *telnetClt) Receive() error {
	if _, err := io.Copy(clt.out, clt.conn); err != nil {
		return fmt.Errorf("failed to receive: %w", err)
	}
	return nil
}

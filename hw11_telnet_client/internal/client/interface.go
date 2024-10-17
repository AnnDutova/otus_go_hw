package client

import (
	"io"
)

type (
	Client interface {
		Connect() error
		io.Closer
		Send() error
		Receive() error
	}
)

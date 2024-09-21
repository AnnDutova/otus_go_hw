package telnet

import (
	"bytes"
	"io"
	"net"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/require"

	"github.com/AnnDutova/otus_go_hw/hw11_telnet_client/internal/terrors"
)

func TestTelnetClient(t *testing.T) {
	t.Run("basic", func(t *testing.T) {
		l, err := net.Listen("tcp", "127.0.0.1:")
		require.NoError(t, err)
		defer func() { require.NoError(t, l.Close()) }()

		var wg sync.WaitGroup
		wg.Add(2)

		go func() {
			defer wg.Done()

			in := &bytes.Buffer{}
			out := &bytes.Buffer{}

			timeoutFlag, err := time.ParseDuration("10s")
			require.NoError(t, err)

			client := NewTelnetClient(l.Addr().String(), timeoutFlag, io.NopCloser(in), out)
			require.NoError(t, client.Connect())
			defer func() { require.NoError(t, client.Close()) }()

			in.WriteString("hello\n")
			err = client.Send()
			require.NoError(t, err)

			err = client.Receive()
			require.NoError(t, err)
			require.Equal(t, "world\n", out.String())
		}()

		go func() {
			defer wg.Done()

			conn, err := l.Accept()
			require.NoError(t, err)
			require.NotNil(t, conn)
			defer func() { require.NoError(t, conn.Close()) }()

			request := make([]byte, 1024)
			n, err := conn.Read(request)
			require.NoError(t, err)
			require.Equal(t, "hello\n", string(request)[:n])

			n, err = conn.Write([]byte("world\n"))
			require.NoError(t, err)
			require.NotEqual(t, 0, n)
		}()

		wg.Wait()
	})

	t.Run("invalid telnet client connection", func(t *testing.T) {
		in := &bytes.Buffer{}
		out := &bytes.Buffer{}

		timeoutFlag, err := time.ParseDuration("10s")
		require.NoError(t, err)

		client := NewTelnetClient("invalidhost:4242", timeoutFlag, io.NopCloser(in), out)
		require.Error(t, client.Connect())
	})

	t.Run("Send/Receive to close connection", func(t *testing.T) {
		l, err := net.Listen("tcp", "127.0.0.1:")
		require.NoError(t, err)
		defer func() { require.NoError(t, l.Close()) }()

		var wg sync.WaitGroup
		wg.Add(2)

		go func() {
			defer wg.Done()

			in := &bytes.Buffer{}
			out := &bytes.Buffer{}

			timeoutFlag, err := time.ParseDuration("10s")
			require.NoError(t, err)

			client := NewTelnetClient(l.Addr().String(), timeoutFlag, io.NopCloser(in), out)
			require.NoError(t, client.Connect())

			in.WriteString("hello\n")
			err = client.Send()
			require.NoError(t, err)

			require.NoError(t, client.Close())

			in.WriteString("End\n")
			err = client.Send()
			require.Error(t, err)

			var cErr terrors.ClientErr
			require.ErrorAs(t, err, &cErr)

			err = client.Receive()
			require.Error(t, err)
			require.ErrorAs(t, err, &cErr)
		}()

		go func() {
			defer wg.Done()

			conn, err := l.Accept()
			require.NoError(t, err)
			require.NotNil(t, conn)

			request := make([]byte, 1024)
			n, err := conn.Read(request)
			require.NoError(t, err)
			require.Equal(t, "hello\n", string(request)[:n])

			conn.Close()
		}()

		wg.Wait()
	})
}

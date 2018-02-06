package core

import (
	"errors"
	"fmt"
	"io"
	"net"
)

// GetError is called to return an appropriate error given some other error.
func GetError(err error) error {
	if err == nil {
		return nil
	}

	if err == io.EOF {
		return ErrConnClosed{}
	}

	switch err.(type) {
	case *net.OpError:
		return ErrOffline{
			err: err.Error(),
		}
	default:
		return err
	}
}

// ErrNoImpl occurs when a package or piece of logic is not implemented.
var ErrNoImpl = errors.New("not implemented")

// ErrInvalidSig occurs when an invalid protocol signature is received.
var ErrInvalidSig = errors.New("invalid protocol signature")

// ErrConnClosed indicates that a remote endpoint closed a connection.
//
// This error should be used when an io.EOF error is returned from a net operation.
type ErrConnClosed struct {
	err string
}

func (e ErrConnClosed) Error() string {
	return "network connection closed"
}

// An ErrOffline error occurs when a remote endpoint crashes / goes offline.
//
// This error should be used when the following error is returned from a net read operation:
//   + err.(*net.OpError)
//     - err.Op = read
//   + err.Err.(*os.SyscallError)
//	   - err.Err.Syscall = wsarecv
//   + err.Err.Err.(*syscall.Errno)
//	   - err.Err.Err.Error() = "An existing connection was forcibly closed by the remote host."
type ErrOffline struct {
	err string
}

func (e ErrOffline) Error() string {
	return fmt.Sprintf("network endpoint has gone offline. %s", e.err)
}

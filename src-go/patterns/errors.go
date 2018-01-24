package patterns

import (
	"errors"
	"fmt"
	"io"
	"net"
)

// Error is called to return an appropriate error given some other error.
func Error(err error) error {
	if err == nil {
		return nil
	}

	if err == io.EOF {
		return ErrConnLost{
			error: err,
		}
	}

	switch err.(type) {
	case *net.OpError:
		return ErrOffline{
			error: err,
		}
	default:
		return err
	}
}

// ErrInvalidSig occurs when an invalid protocol signature is received.
var ErrInvalidSig = errors.New("invalid protocol signature")

// ErrConnLost indicates that a remote endpoint dropped a connection.
//
// This error should be used when an io.EOF error is returned from a net operation.
type ErrConnLost struct {
	error
}

func (e ErrConnLost) String() string {
	return "network connection to lost"
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
	error
}

func (e ErrOffline) String() string {
	return fmt.Sprintf("network endpoint has gone offline. %s", e.error)
}

# Go Source Code

## Server Errors

drops connection:

+ err == io.EOF

server crashes:

+ err.(*net.OpError)
  - err.Op = read
+ err.Err.(*os.SyscallError)
  - err.Err.Syscall = wsarecv
+ err.Err.Err.(*syscall.Errno)
  - err.Err.Err.Error() = "An existing connection was forcibly closed by the remote host."
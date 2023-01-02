# Example: Send a file descriptor over a Unix-domain socket

Usage:

```shell
# Terminal A:
go run ./recv -socket test.sock

# Terminal B:
go run ./send -file go.sum -socket test.sock
```

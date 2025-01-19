# Example: Call a subprocess with a service pipe

Usage:

```shell
go build -o phost ./host
go build -o papp ./app

./phost -- ./papp a b c
```

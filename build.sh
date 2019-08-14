export GOPATH=$(dirname $(readlink -f $0))
CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o bin/main .

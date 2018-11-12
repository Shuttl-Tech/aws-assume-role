from golang:1.11.2-stretch as base
workdir /go/src/Shuttl-Tech/aws-assume-role
add . .
run go get -d -v ./...
run CGO_ENABLED=0 go build -a -tags netgo -ldflags '-w -extldflags "-static"' -o /plugin.bin

from scratch
copy --from=base /plugin.bin /plugin

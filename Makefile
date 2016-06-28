APP = x-conf

build:
	go build -o ../bin/${APP} -ldflags '-s -w'

linux:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ../bin/${APP} -ldflags '-s -w'

run:
	@go run *.go

clean:
	@rm ../bin/${APP}

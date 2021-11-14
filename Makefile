.PHONY: build clean deploy

build: clean 
	export GO111MODULE=on
	env CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/generate-ticket generate-ticket/main.go
	env CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/complete-ticket complete-ticket/main.go
	env CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/pay-ticket pay-ticket/main.go
	env CGO_ENABLED=0 GOARCH=amd64 GOOS=linux go build -ldflags="-s -w" -o bin/process-events process-events/main.go

clean:
	rm -rf ./bin 

deploy: clean build
	sls deploy --verbose
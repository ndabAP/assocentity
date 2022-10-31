build:
	env GOOS=linux GOARCH=amd64 go build -o bin/assocentity cmd/main.go
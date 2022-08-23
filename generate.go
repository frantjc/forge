package forge

//go:generate buf format -w

//go:generate buf generate .

//go:generate go fmt ./...

// TODO go:generate golangci-lint run --fix

//go:generate env GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" -o ./internal/bin/shim ./internal/cmd/shim

//go:generate upx ./internal/bin/shim

//go:generate go mod tidy

install: 
	go mod tidy
	go install

test:
	go test -race -timeout 5s -v ./...

coverage:
	go test -race -timeout 5s -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out


fmt:
	gofmt -w .
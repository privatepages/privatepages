cd src

go mod verify
go build -v ./...
go vet ./...
go install honnef.co/go/tools/cmd/staticcheck@latest
staticcheck ./...
go install golang.org/x/lint/golint@latest
golint -set_exit_status -min_confidence 0.5 ./...
go test -race -vet=off ./...
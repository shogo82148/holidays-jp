.PHONY: build-TriggerFunction
build-TriggerFunction:
	GOOS=linux GOARCH=arm64 go build -o $(ARTIFACTS_DIR)/bootstrap -tags lambda.norpc main.go

.PHONY: test
test:
	go test -v ./...

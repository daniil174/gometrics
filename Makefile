GOLANGCI_LINT_CACHE?=/tmp/praktikum-golangci-lint-cache
SERVER_PORT := 37797
ADDRESS := "localhost:37797"
TEMP_FILE := "./temp"

.PHONY : all
all: preproc build-all

.PHONY : preproc
preproc: clean fmt _golangci-lint-run test

.PHONY : build-all
build-all: clean server agent

server:
	go build -o ./bin/server ./cmd/server/main.go

agent:
	go build -o ./bin/agent ./cmd/agent/main.go


test:
	go test ./... -race -coverprofile=cover.out -covermode=atomic

.PHONY : clean
clean:
	rm -f ./bin/agent
	rm -f  ./bin/server
	rm -f  ./cover.out

.PHONY : check-coverage
check-coverage:
	go tool cover -html cover.out

.PHONY : fmt
fmt:
	go fmt ./...
	goimports -l -w .

#.PHONY : lint
#lint: golangci-lint run ./...



.PHONY : run-autotests
run-autotests: iter1 iter2 iter3 iter4 iter5 #iter6

.PHONY : iter1
iter1:
	./bin/metricstest -test.run=^TestIteration1$$ -binary-path=./bin/server

.PHONY : iter2
iter2:
	/bin/metricstest -test.run=^TestIteration2A$$ -source-path=. -agent-binary-path=./bin/agent

.PHONY : iter3
iter3:
	./bin/metricstest -test.run=^TestIteration3A$$ -source-path=. -agent-binary-path=./bin/agent -binary-path=./bin/server
	./bin/metricstest -test.run=^TestIteration3B$$ -source-path=. -agent-binary-path=./bin/agent -binary-path=./bin/server

.PHONY : iter4
iter4:
	./bin/metricstest -test.run=^TestIteration4$$ -source-path=. -agent-binary-path=./bin/agent -binary-path=./bin/server -server-port=$(SERVER_PORT)

.PHONY : iter5
iter5:
	./bin/metricstest -test.run=^TestIteration5$$ -agent-binary-path=./bin/agent -binary-path=./bin/server -server-port=$(SERVER_PORT) -source-path=.

.PHONY : iter6
iter6:
	./bin/metricstest -test.run=^TestIteration6$$ -agent-binary-path=./bin/agent -binary-path=./bin/server -server-port=$(SERVER_PORT) -source-path=.

.PHONY: golangci-lint-run
golangci-lint-run: _golangci-lint-rm-unformatted-report

.PHONY: _golangci-lint-reports-mkdir
_golangci-lint-reports-mkdir:
	mkdir -p ./golangci-lint

.PHONY: _golangci-lint-run
_golangci-lint-run: _golangci-lint-reports-mkdir
	-docker run --rm \
    -v $(shell pwd):/app \
    -v $(GOLANGCI_LINT_CACHE):/root/.cache \
    -w /app \
    golangci/golangci-lint:v1.57.2 \
        golangci-lint run \
            -c .golangci.yml \
	> ./golangci-lint/report-unformatted.json

.PHONY: _golangci-lint-format-report
_golangci-lint-format-report: _golangci-lint-run
	cat ./golangci-lint/report-unformatted.json | jq > ./golangci-lint/report.json

.PHONY: _golangci-lint-rm-unformatted-report
_golangci-lint-rm-unformatted-report: _golangci-lint-format-report
	rm ./golangci-lint/report-unformatted.json

.PHONY: golangci-lint-clean
golangci-lint-clean:
	sudo rm -rf ./golangci-lint

.PHONY: run fmt test cover infra-up infra-up infra-test infra-local infra-local-down infra-upd lint setup

# This assumes tflocal is installed https://github.com/localstack/terraform-local

all: infra-down infra-up infra-test
infra-upd:
	cd Docker && docker-compose -f docker-compose.yaml up -d

infra-up:
	cd Docker && docker-compose -f docker-compose.yaml up 

infra-down:
	cd Docker && docker-compose -f docker-compose.yaml down 

infra-test:
	sleep 5 && aws --region us-east-1 dynamodb list-tables --endpoint-url http://localhost:4566 --no-cli-pager

infra-local:
	cd terraform && \
	export TF_LOG=INFO && \
	tflocal init && \
	tflocal apply -auto-approve

infra-local-down:
	cd terraform && export TF_LOG=INFO && tflocal destroy -auto-approve

fmt:
	go fmt ./... && cd terraform && terraform fmt

lint: 
	golangci-lint run

test:
	go test -timeout 10000ms -v ./... -covermode=count -coverprofile=cover.out && go tool cover -func=cover.out

cover: test
	go tool cover -html=cover.out -o coverage.html

run: fmt lint
	go run ./cmd/simplechat/main.go

setup:  infra-upd infra-local run 

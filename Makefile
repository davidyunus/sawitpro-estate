

.PHONY: clean all init generate generate_mocks

all: build/main

build/main: main.go generated
	@echo "Building..."
	go build -o $@ $<

clean:
	rm -rf generated

init: clean generate
	go mod tidy
	go mod vendor

mockgen:
	mockgen -source=src/domain/estate.go -destination=src/mock/estate.go
	mockgen -source=src/domain/palm_tree.go -destination=src/mock/palm_tree.go

test:
	go clean -testcache
	go test -short -coverprofile coverage.out -short -v ./...

test_api:
	go clean -testcache
	go test -v -coverprofile=cover.out ./...

generate: generated generate_mocks

generated: api.yml
	@echo "Generating files..."
	mkdir generated || true
	oapi-codegen --package generated -generate types,server,spec $< > generated/api.gen.go

INTERFACES_GO_FILES := $(shell find src/domain/estate.go)
INTERFACES_GEN_GO_FILES := $(INTERFACES_GO_FILES:%.go=%.mock.gen.go)

generate_mocks: $(INTERFACES_GEN_GO_FILES)
$(INTERFACES_GEN_GO_FILES): %.mock.gen.go: %.go
	@echo "Generating mocks $@ for $<"
	mockgen -source=$< -destination=$@ -package=$(shell basename $(dir $<))
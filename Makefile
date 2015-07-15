NAME := xhyvectl

VERSION := 0.1.0

PROJECT := github.com/ailispaw/$(NAME)

WORKSPACE := `godep path`

GIT_COMMIT := `git rev-parse --short HEAD`

all: build

get:
	godep get ./...

fmt:
	go fmt -x ./...

vet:
	go vet -x ./...

test: restore
	godep go test ./...

build: fmt vet
	go build -v

install: build
	@install -CSv $(NAME) $(GOPATH)/bin/

uninstall:
	go clean -x -i

clean:
	go clean -x
	$(RM) -r "$(WORKSPACE)"

save:
	godep save

update:
	godep update ...
	$(RM) -r "$(WORKSPACE)"

restore:
	GOPATH="$(WORKSPACE)" godep restore

.PHONY: all get fmt vet test build install uninstall clean save update restore

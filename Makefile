.PHONY: clean dep dep-vendor-only test

clean:
	@rm -rf bin

dep:
	@dep ensure -v

dep-vendor-only:
	@dep ensure -v -vendor-only

test:
	@go test -v -race -cover ./...

build:
	CGO_ENABLED=0 \
	go build \
		-o bin/pdb-monitor-controller \
		-ldflags "-X main.version=$(VERSION)" \
		main.go

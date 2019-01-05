.PHONY: build build-environment-build lambda-environment-build

build:
	docker-compose run --rm build-environment make build-environment-build
	docker-compose run --rm lambda-environment make lambda-environment-build

build-environment-build:
	rm -rf ./build
	CGO_CFLAGS='$(shell pkg-config --cflags librsvg-2.0)' CGO_LDFLAGS='-L/usr/lib $(shell pkg-config --libs librsvg-2.0) -Wl,-rpath,$$ORIGIN/lib' go build -o ./build/main .
	./build/main --help
	mkdir ./build/lib
	cp \
		/usr/lib/*.so \
		/usr/lib/*.so.* \
		./build/lib
	patchelf --set-rpath '$$ORIGIN' ./build/lib/*.so.*

lambda-environment-build:
	./build/main --help || (ldd ./build/main && exit 1)

vendor: ../Gopkg.lock ../Gopkg.toml
	dep ensure

test: vendor build-environment-build
	cd .. && CGO_CFLAGS='$(shell pkg-config --cflags librsvg-2.0)' CGO_LDFLAGS='$(shell pkg-config --libs librsvg-2.0) -Wl,-rpath,$$ORIGIN/lib' go test -cover -race -v ./...

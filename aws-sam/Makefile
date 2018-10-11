.PHONY: build build-environment-build lambda-environment-build

build:
	docker-compose run --rm build-environment make build-environment-build
	docker-compose run --rm lambda-environment make lambda-environment-build

build-environment-build:
	rm -rf ./build
	CGO_CFLAGS='$(shell pkg-config --cflags librsvg-2.0)' CGO_LDFLAGS='$(shell pkg-config --libs librsvg-2.0) -Wl,-rpath,$$ORIGIN/lib' go build -o ./build/main .
	./build/main --help
	mkdir ./build/lib
	cp \
		/usr/lib/x86_64-linux-gnu/librsvg-2.so.* \
		/usr/lib/x86_64-linux-gnu/libgdk_pixbuf-2.0.so.* \
		/usr/lib/x86_64-linux-gnu/libpangocairo-1.0.so.* \
		/usr/lib/x86_64-linux-gnu/libpango-1.0.so.* \
		/usr/lib/x86_64-linux-gnu/libcroco-0.6.so.* \
		/usr/lib/x86_64-linux-gnu/libpng16.so.* \
		/usr/lib/x86_64-linux-gnu/libpangoft2-1.0.so.* \
		/usr/lib/x86_64-linux-gnu/libthai.so.* \
		/usr/lib/x86_64-linux-gnu/libharfbuzz.so.* \
		/usr/lib/x86_64-linux-gnu/libdatrie.so.* \
		/usr/lib/x86_64-linux-gnu/libgraphite2.so.* \
		/usr/lib/x86_64-linux-gnu/libgobject-2.0.so.* \
		/lib/x86_64-linux-gnu/libglib-2.0.so.* \
		/lib/x86_64-linux-gnu/libpcre.so.* \
		./build/lib
	patchelf --set-rpath '$$ORIGIN' ./build/lib/*.so.*

lambda-environment-build:
	./build/main --help || (ldd ./build/main && exit 1)

vendor: ../Gopkg.lock ../Gopkg.toml
	dep ensure

test: vendor build-environment-build
	cd .. && CGO_CFLAGS='$(shell pkg-config --cflags librsvg-2.0)' CGO_LDFLAGS='$(shell pkg-config --libs librsvg-2.0) -Wl,-rpath,$$ORIGIN/lib' go test -cover -race -v ./...

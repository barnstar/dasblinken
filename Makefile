.PHONY: clean
clean:
	rm -f piled

# We need to link in a cGo library and building on docker is the way
.PHONY: builder-image
builder-image:
	docker buildx build --platform linux/arm/v7 --tag ws2811-builder --file docker/app-builder/Dockerfile .

# Create cache directories on host if they don't exist
.PHONY: cache-dirs
cache-dirs:
	@mkdir -p $(HOME)/.cache/go-build
	@mkdir -p $(HOME)/go/pkg/mod

piled: cache-dirs
	docker run --rm -v "$(PWD)":/usr/src/piled --platform linux/arm/v7 \
		-v "$(HOME)/.cache/go-build:/root/.cache/go-build" \
		-v "$(HOME)/go/pkg/mod:/go/pkg/mod" \
		-w /usr/src/piled/example ws2811-builder:latest go build -o "../piled" -v -ldflags "-linkmode external -extldflags -static"

all: clean builder-image piled
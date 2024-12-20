.PHONY: clean
clean:
	rm -f piled-armv7

# Copies the piled-armv7 binary over to your pi
# Your device may not be called minipi but if it is, lucky you!  If not, change it.
# I'm running Tailscale and Tailscale ssh on this nugget so dasnetwork alwasy dasworks.
.PHONY: deploy
deploy: piled-armv7
	rsync -ave ssh ./piled-armv7 admin@minipi:/home/admin/piled/piled

# We need to link in a cGo library and building on docker is the way
.PHONY: builder-image
builder-image:
	docker buildx build --platform linux/arm/v7 --tag ws2811-builder --file docker/app-builder/Dockerfile .

piled-armv7:
	docker run --rm -v "$(PWD)":/usr/src/piled --platform linux/arm/v7 \
 		-w /usr/src/piled ws2811-builder:latest go build -o "piled-armv7" -v -ldflags "-linkmode external -extldflags -static"

all: clean piled-armv7
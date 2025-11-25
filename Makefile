# Your device may not be called minipi but if it is, lucky you!  If not, change it.
# I'm running Tailscale and Tailscale ssh on this nugget so dasnetwork alwasy dasworks.ssh 
HOST=minipi
USER=admin
DEST_DIR=/home/admin/piled

.PHONY: clean
clean:
	rm -f piled-armv7

# Copies the piled-armv7 binary over to your pi
.PHONY: deploy
deploy: piled-armv7
	rsync -ave ssh ./piled $(USER)@$(HOST):$(DEST_DIR)/piled
	rsync -ave ssh ./default.json $(USER)@$(HOST):$(DEST_DIR)/default.json

.PHONY: run
run:
	ssh -n -f $(USER)@$(HOST) "sh -c 'sudo killall -9 piled > /dev/null 2>&1; cd $(DEST_DIR); nohup sudo ./piled -config=> /dev/null 2>&1 &'"

# We need to link in a cGo library and building on docker is the way
.PHONY: builder-image
builder-image:
	docker buildx build --platform linux/arm/v7 --tag ws2811-builder --file docker/app-builder/Dockerfile .

example:
	docker run --rm -v "$(PWD)":/usr/src/piled --platform linux/arm/v7 \
 		-w /usr/src/piled/example ws2811-builder:latest go build -o "../piled" -v -ldflags "-linkmode external -extldflags -static"

all: clean example
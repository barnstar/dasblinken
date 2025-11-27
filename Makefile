# Your device may not be called minipi but if it is, lucky you!  If not, change it.
# I'm running Tailscale and Tailscale ssh on this nugget so dasnetwork alwasy dasworks.ssh 
HOST=minipi
USER=admin
DEST_DIR=/home/admin/piled
CONFIG=config.json
TSNAME=dasblinken
AUTHKEY ?= $(shell echo $$DBL_AUTHKEY)


.PHONY: clean
clean:
	rm -f piled

# Copies the piled binary and the default.json config over to your pi
.PHONY: deploy
deploy: example
	rsync -ave ssh ./piled $(USER)@$(HOST):$(DEST_DIR)/piled
	rsync -ave ssh ./$(CONFIG) $(USER)@$(HOST):$(DEST_DIR)/$(CONFIG) 
	rsync -ave ssh ./effects/effects.json $(USER)@$(HOST):$(DEST_DIR)/effects.json 

.PHONY: run
run:
	ssh -n -f $(USER)@$(HOST) "sh -c 'sudo killall -9 piled > /dev/null 2>&1; cd $(DEST_DIR); nohup sudo ./piled --config=$(CONFIG) --authkey=$(AUTHKEY) --tsname=$(TSNAME) > /dev/null 2>&1 &'"

# We need to link in a cGo library and building on docker is the way
.PHONY: builder-image
builder-image:
	docker buildx build --platform linux/arm/v7 --tag ws2811-builder --file docker/app-builder/Dockerfile .

piled:
	docker run --rm -v "$(PWD)":/usr/src/piled --platform linux/arm/v7 \
 		-w /usr/src/piled/example ws2811-builder:latest go build -o "../piled" -v -ldflags "-linkmode external -extldflags -static"

all: clean piled
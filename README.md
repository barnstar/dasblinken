# DASBLINKEN

For making dasblinkenlights with a Pi in go.

This does all the fancy work of cross compiling for armv7 and linking in ws281x support bindings
so you can focus on dasblinken and not dasboilerplate.

With many thanks to https://pkg.go.dev/github.com/rpi-ws281x/rpi-ws281x-go

## Deploying

The Makefile assumes your pi is reachable as 'minipi'.. Yours might not be.  
Put Tailscale on your pi and set it's MagicDNS name as minipi and you're
off to the races (and best of all, you can change dasblinkeneffect from 
dasanywhere)

The only support arch is 32 bit armV7.  A pi zero 2w running bare bones
Rasbian Bookwork (Debian) will get it done.  

We use GPIO29 (last top pin) by default.  You will need to wire up the proper
level shifter.  We tested all of this with ws2812 strips.

We assume no responsibility for dasMagicSmoke...

```
% make builder-image
% make clean deploy
```

You can then http://minipi:8080

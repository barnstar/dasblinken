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

We use GPIO29 (last top pin) by default.  Often, no level shifter is needed.
YMMV, but if you see instability, a 3.3 -> 5V level shifter might be required

You probably want to modify the constants in the makefile.  The auth key is pulled
out of the local env.  Everything else is hardcoded to values that probably aren't 
what you want.

Modify config.json to match your LED strip.  The effects package includes a number
of effects that will work on strips and matrices.

The piled app runs as a tsnet client on your tailnet too. 



We assume no responsibility for dasMagicSmoke...

```bash
% export AUTHKEY=<an auth key for your tailnet>
% make builder-image
% make all
# make deploy
% make run
```

You can then http://dasblinken

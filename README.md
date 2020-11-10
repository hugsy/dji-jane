## DJI-Jane ##

`DJI-Jane` is the HTTP REST server collecting the information sent by the
[`DJI-Joe`](https://github.com/hugsy/dji-joe) probes. It will
collect all events submitted by drones and send decisions to them.

### Compilation (non-required)

Requires GoLang 1.8+

```
$ go get github.com/apsdehal/go-logger
$ go build -o bin/dji-jane-$(arch) src/main/main.go
```

### Runtime

If not compiled,
```
$ go run src/main/main.go
```

Else
```
$ bin/dji-jane
```

Which will by default listen to `0.0.0.0:8000` for HTTP events (can be changed
via the `--listen` command-line switch), which should look like this:

```
$ bin/dji-jane -listen 0.0.0.0:8000
#1 2017-06-16 00:56:08 main.go:16 ▶ INF Starting DJI-Jane [0.1]
#2 2017-06-16 00:56:08 djijane.go:24 ▶ DEB Preparing server
#3 2017-06-16 00:56:08 djijane.go:33 ▶ DEB Placing sighandlers
#4 2017-06-16 00:56:08 djijane.go:42 ▶ INF Starting REST server listening on '0.0.0.0:8000'
#5 2017-06-16 00:56:12 handlers.go:100 ▶ NOT Received WAKEUP from 'rpi-probe-2' at 2017-06-16 00:56:12.842589056 +0000 UTC (lat=0.000 long=0.000)
#6 2017-06-16 00:56:12 handlers.go:102 ▶ DEB Added '<Probe name='rpi-probe-2'>' to `probes[]`
#7 2017-06-16 00:56:41 handlers.go:155 ▶ NOT Received NEWDRONEINFO ProbeRequest by 'rpi-probe-2' at 2017-06-16 00:56:41.981361339 +0000 UTC from 60:60:1f:42:11:b8 (DJI) with strength: -8 dBm
#8 2017-06-16 00:56:41 handlers.go:155 ▶ NOT Received NEWDRONEINFO ProbeRequest by 'rpi-probe-2' at 2017-06-16 00:56:41.991112345 +0000 UTC from 60:60:1f:42:11:b8 (DJI) with strength: -12 dBm
#9 2017-06-16 00:56:42 handlers.go:155 ▶ NOT Received NEWDRONEINFO ProbeRequest by 'rpi-probe-2' at 2017-06-16 00:56:42.051174516 +0000 UTC from 60:60:1f:42:11:b8 (DJI) with strength: -10 dBm
[...]
#45 2017-06-16 00:56:57 handlers.go:155 ▶ NOT Received NEWDRONEINFO ProbeRequest by 'rpi-probe-2' at 2017-06-16 00:56:57.391200818 +0000 UTC from 60:60:1f:42:11:b8 (DJI) with strength: -34 dBm
#46 2017-06-16 00:56:57 handlers.go:155 ▶ NOT Received NEWDRONEINFO ProbeRequest by 'rpi-probe-2' at 2017-06-16 00:56:57.395023721 +0000 UTC from 60:60:1f:42:11:b8 (DJI) with strength: -34 dBm
#47 2017-06-16 00:57:00 handlers.go:124 ▶ NOT Received SHUTDOWN from 'rpi-probe-2' at 2017-06-16 00:57:00.701122946 +0000 UTC: nb_beacons: 0, nb_probes: 40
#48 2017-06-16 00:57:00 handlers.go:132 ▶ DEB Removed '<Probe name='rpi-probe-2'>' from `probes[]`
^C
#49 2017-06-16 00:57:09 djijane.go:38 ▶ INF Got 'interrupt': stopping cleanly
#50 2017-06-16 00:57:09 main.go:18 ▶ INF Leaving DJI-Jane
```

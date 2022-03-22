# raza : keep history forever

![](https://img.shields.io/github/license/wesen/raza)
![](https://img.shields.io/github/workflow/status/wesen/raza/golang-pipeline)

raza is a tool (client, server) meant to keep your shell history forever,
across devices. Every command you enter is stored in a database, and available
for query. It automatically synchronizes across devices, and works even if you are offline.

See [hacking notes](hacking-notes.md) for a free-flowing brainstorm of what this app should be.


## Configuration and logging

Per default, log to stderr. If stderr is a tty, use colors.

Adds the following options in the `~/.config/raza.yaml` config file (if present):

```yaml
root:
  address: localhost # used for grpc later
  debug: true|false # set logging level to debug
  log-error-stacktrace: true|false # log error stacktraces when using Error() level
  log-file: foobar # log to file
  log-line: true|false # (default true) log linenumber
```


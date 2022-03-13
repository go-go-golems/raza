## Todos

* add zerolog structured logging
* add viper configuration
* start daemon when needed
* generate sql stubs using xo

## Ideas

* add command to start/end transaction
* add metadata/tags/name to sessions
* hook on shell closing (if possible)
* resort to file backing when the daemon can't be started
* daemonless operation (straight to sqlite)
* multiple SQL backends
* add custom hook for programs to add their own metadata
* should we have tags besides just metadata? as an append set
* blacklist when an env variable is set / when a metadata value is set / when a tag is present
* register programs / API to be run to add metadata before a command is logged
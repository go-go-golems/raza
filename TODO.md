## Todos

* [x] add zerolog structured logging RAZA-1
* [x] add viper configuration RAZA-2
* [x] add command to start/end transaction (transactions are called sessions now)
* start daemon when needed

## Ideas

* add metadata/tags/name to sessions
* hook on shell closing (if possible)
* resort to file backing when the daemon can't be started
* daemonless operation (straight to sqlite)
* multiple SQL backends
* add custom hook for programs to add their own metadata
* should we have tags besides just metadata? as an append set
* blacklist when an env variable is set / when a metadata value is set / when a tag is present
* register programs / API to be run to add metadata before a command is logged
* should EndCommand/EndSession be allowed to add tags?
* should we be able to add tags to a session as it is running?
* what about categorizing sessions and commands after the fact?
* Potentially flagging / bookmarking / making curated selections of commands in the UI
* Ignore failed commands

* can we record the commands executed by a shell script, for example logging the execution of a shell script to our history?

## Deprecated ideas

* generate sql stubs using xo ??


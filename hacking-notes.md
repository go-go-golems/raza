# raza : keep history forever

## Overview

I want to build an app that:
* synchronizes shell history and keeps it forever
* across multiple devices
* adds metadata to each logged line
    * date
    * directory
    * git repository
    * git hash
    * hostname
    * OS
    * terminal session (when started, when ended)
    * execution time (can we get that?)
    * arbitrary key values / environment values
* query history quickly from the local device
* works offline if wanted, easy setup of daemon
    * can gather history until it is back online again and synchronizes
* make a little webapp to search history, say when you are at another person's computer

There are different aspects to design:
* data schema
    * seems pretty simple, one line per history entry
    * keep track of history sessions
    * we can probably store all the git directory metadata with each line
* how do we gather the data in the first place?
    * are there hooks to get the directory and co in bash/zsh
    * do we want to run the git command each time? How does starship do it?
* local storage until remote synchronization
    * this could actually be the same as just doing the local storage, and we could potentially use the same scheme with watermill and sql readers, storing to a local sqlite
    * how would we sync remote state back to a local search engine for quick access?
    * maybe something like segments, keeping multiple local sqlite that can be searched quickly? What about full text search in sqlite?
* remote synchronization
    * best solved with a local daemon
    * potentially use a system of plugins to make the same daemon architecture work when deployed in the cloud
* search query API / analysis tools

Is the protocol between local daemon and local client  the same as between local daemon  and cloud service to send data?

## Extracting the history out of zsh
Let’s look at this to see how they hook into zsh: [GitHub - larkery/zsh-histdb: A slightly better history for zsh](https://github.com/larkery/zsh-histdb)

It seems the magic is done through the `zsh-add-hook zshaddhistory`
* [mastering-zsh/hooks.md at master · rothgar/mastering-zsh · GitHub](https://github.com/rothgar/mastering-zsh/blob/master/docs/config/hooks.md)

## Configuration options
* filter out certain comments
* filter out certain directories
* enable/disable git integration
* which environment variables should be logged

## Extensibility
* make it possible to add your own key/values to history entries
    * say, over environment variables
* potentially run additional commands when adding history / plugin interface
* API for the service that receives the log lines
    * get notified when a new line is added
    * easily query the DB over an API
    * push to the service

## Command line tool
* query history
    * globally
    * per host
    * per day
    * per directory
    * per git repo
    * per git branch
* start session "transaction"
* restrict where to look for history (per directory with a file? With some environment variables?)
* can we recognize in which session the client runs, and then store the config that in the backend server so that we don’t’ have to set environment variables and the like?
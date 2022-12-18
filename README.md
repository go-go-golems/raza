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


```
there exists some code that transforms:

flags = --day 2022-03-25 --select "*aws*" --tally

parsedFlags = parseFlags(flags)

parsedFlags = {
   day: 2022-03-25,
   select: '*aws*',
   tally: true
}

abstractQuery = buildAbstractQuery(parsedFlags)

abstractQuery = 
{ 
   day: 2022-03-25,
   query: {
      type:aggregate,
      aggregation: count,
      filter: {
         column: 'name',
         filter: '*aws*'
      }
   },
   from: tables.commands
}

selectStatement = buildSelectStatement(abstractQuery)
into:
SELECT COUNT(*) FROM commands WHERE command LIKE '%aws%' AND day='2022-03-25'

explainedSelectAggregateStatement = buildSelectStatement(abstractQuery, explain=true)
into:
EXPLAIN SELECT COUNT(*) FROM commands WHERE command LIKE '%aws%' AND day='2022-03-25'

rawDataSelectStatement = buildSelectStatement(abstractQuery, explain=true)
into:
SELECT * FROM commands WHERE command LIKE '%aws%' AND day='2022-03-25'

---
flags = --tally \
     --from 2022-03 --to 2022-05 \
     --filter repo=my-big-aws-project \
     --select "*deploy*" \
     --sparklines

parsedFlags = {
   dateRange: {
      from: 2022-03,
      to: 2022-05
   },
   select: '*aws*',
   filter: 
   tally: true
}

abstractQuery = buildAbstractQuery(parsedFlags)

abstractQuery = 
{ 
   day: 2022-03-25,
   query: {
      type:aggregate,
      aggregation: count,
      filter: {
         column: 'name',
         filter: '*aws*'
      },
      groupBy: {
         time: day
      }
   },
   from: tables.commands
} 
```
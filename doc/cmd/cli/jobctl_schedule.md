## jobctl schedule

Schedule a job

### Synopsis

Schedule a job to run with the specified kind and arguments.

```
jobctl schedule kind [args...] [flags]
```

### Examples

```
job schedule kind arg1 arg2
```

### Options

```
  -h, --help   help for schedule
```

### Options inherited from parent commands

```
      --host string   gRPC host or address for a go-job instance (default "localhost")
      --port int      gRPC port number for a go-job instance (default 59051)
```

### SEE ALSO

* [jobctl](jobctl.md)	 - Job Control CLI


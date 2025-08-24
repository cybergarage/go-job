## jobctl cancel instances

Cancel job instances

### Synopsis

Cancel job instances by the specified query.

```
jobctl cancel instances [flags]
```

### Options

```
  -h, --help          help for instances
  -k, --kind string   Kind of the instances to cancel
  -u, --uuid string   UUID of the instances to cancel
```

### Options inherited from parent commands

```
      --host string   gRPC host or address for a go-job instance (default "localhost")
      --port int      gRPC port number for a go-job instance (default 59051)
```

### SEE ALSO

* [jobctl cancel](jobctl_cancel.md)	 - cancel the specified resource


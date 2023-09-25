## To run
```
cd to project root
go run . serve [--flag value]
```

## Configuration priority
1. Flag variables
2. ENV variables
3. yaml variables

### Env variables
If used as environment variable it must be prefixed with SL_, and be in all caps. Use _ instead of -.

For yaml: 
```
hostname: testname
port: 30333
metrics-port: 32111
sender-type: msGraph
debug: true
ms-graph-config:
  tenant-id: 
  client-id: 
  client-secret: 
  sender-user-id: 
basic-auth-config:
  enabled: true 
  username:
  password
```

The same variables can be used as flags, ie: `./skyline serve --port 123`

## Sending emails
To send emails you need to use an SMTP client like Thunderbird. You don't need to configure incoming settings.
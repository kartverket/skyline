# skyline

A tool for bridging good old SMTP and OAuth2 protected email services. Can be run as a sidecar (recommended) or as a standalone service.

> **Note:** This project is still in an early state and may be unsuitable for production workloads. PR's are most welcome 🙌

## Features

- Expose local SMTP server with optional basic authentication
- Send email via Microsoft Graph (Office 365)
- Prometheus metrics (default `:5353/metrics`)


## Configuration

The project can be configured using either command line flags, using a configuration file or with environment variables.

### Configuration file
The default location is `~/.skyline.yaml`
```yaml
debug: false
hostname: an-overridden-hostname                 # autodetected by default
port: 30333
metrics-port: 32111
sender-type: msGraph                             # currently the only implementation
ms-graph-config:
  tenant-id: <entra-id-tenant-id>
  client-id: <id>                                # app registration id
  client-secret: <secret>                        # app registration client secret
  sender-user-id: <user-object-id>               # found by viewing details of an user in Entra ID 
basic-auth-config:
  enabled: true 
  username: foo
  password: bar
```

### Environment variables

All configuration properties can be specified as environment variables by replacing structure indentation, hyphens and spaces with `_`. All environment variables must be prefixed with `SL_`. 

## Run

Use `go run`, `go build && ./skyline` or use one of the [prebuilt container images](https://github.com/kartverket/skyline/pkgs/container/skyline).
You can also create your own local image with `goreleaser release --snapshot --clean`

## Testing with Grafana Oncall
``` 
1. create file with envs local_testing/.env.skyline
    SL_BASIC_AUTH_CONFIG_ENABLED=TRUE
    SL_BASIC_AUTH_CONFIG_USERNAME=test
    SL_BASIC_AUTH_CONFIG_PASSWORD=test
    SL_MS_GRAPH_CONFIG_TENANT_ID=
    SL_MS_GRAPH_CONFIG_CLIENT_ID=
    SL_MS_GRAPH_CONFIG_CLIENT_SECRET=
    SL_MS_GRAPH_CONFIG_SENDER_USER_ID=
2. docker compose up -d
3. go to grafana at http://localhost:3000, user:pass admin:admin
4. enable oncall plugin: Administration -> Plugins -> Search for oncall -> Oncall backend url: http://engine:8080
5. edit oncall settings to use skyline. Host: Skyline, TLS: Off
```

### rebuild skyline container
`docker compose up --build skyline`
# MongoDb Exporter

| Status                   |                       |
| ------------------------ | --------------------- |
| Stability                | [alpha]               |
| Supported pipeline types | traces, logs          |
| Distributions            | [core], [contrib]     |

Exporter supports the following features：

+ Support for writing pipeline data to a file.

+ Support for rotation of telemetry files.

+ Support for compressing the telemetry data before exporting.


Please note that there is no guarantee that exact field names will remain stable.
This intended for primarily for debugging Collector without setting up backends.

The official [opentelemetry-collector-contrib container](https://hub.docker.com/r/otel/opentelemetry-collector-contrib/tags#!) does not have a writable filesystem by default since it's built using the special `from scratch` layer. As such, you will need to create a writable directory for the path, potentially by creating writable volumes or creating a custom image.

## Getting Started

The following settings are required:

- `conn_uri` [no default]: mongodb connection uri.

The following settings are optional:

- `logs_collection`: [default: Logs] settings to create a collection for storing logs data from otel agent.
- `traces_collection`: [default: Requests] settings to create a collection for storing traces data from otel agent.
- `db`: [default: OtelDB] settings to create a db for storing traces/logs data from otel agent.

## Example:

```yaml
exporters:
  mongodb:
    conn_uri: mongodb://foo:bar@localhost:27017
    logs_collection: Logs_coll
    traces_collection: Requests_coll
    db: OtelData

```

[alpha]:https://github.com/open-telemetry/opentelemetry-collector#alpha
[contrib]:https://github.com/open-telemetry/opentelemetry-collector-releases/tree/main/distributions/otelcol-contrib
[core]:https://github.com/open-telemetry/opentelemetry-collector-releases/tree/main/distributions/otelcol

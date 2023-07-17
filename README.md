# MongoDb Exporter

| Status                   |                       |
| ------------------------ | --------------------- |
| Stability                | [alpha]               |
| Supported pipeline types | traces, logs          |

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

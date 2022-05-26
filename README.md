# spanner-createdb

Create spanner instances/databases (standalone Go script)

## Install

```bash
go install github.com/lukasjapan/spanner-createdb@latest
```

## Usage

```
Usage:
  spanner-createdb {databaseId}
  spanner-createdb databases/{databaseId}
  spanner-createdb {instanceId}/databases/{databaseId}
  spanner-createdb instances/{instanceId}/databases/{databaseId}
  spanner-createdb {projectId}/instances/{instanceId}/databases/{databaseId}
  spanner-createdb projects/{projectId}/instances/{instanceId}/databases/{databaseId}

You can also pass the ids via environment variables:
  SPANNER_PROJECT_ID
  SPANNER_INSTANCE_ID
  SPANNER_DATABASE_ID
```

Further configuration can not be provided.
This script is intended to be mainly used with the spanner emulator.
Do not forget to set the `SPANNER_EMULATOR_HOST` environment variable.

### Example

```bash
export SPANNER_EMULATOR_HOST=localhost:9010
export SPANNER_PROJECT_ID=my-dev-project
export SPANNER_INSTANCE_ID=my-dev-instance
spanner-createdb my-dev-db
```

## Motivation

To provide a standalone script create new spanner instances/databases without the Google Cloud SDK.

Ideal for including in docker containers/images or if you do not want to install the cloud sdk.

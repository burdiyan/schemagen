Confluent Kafka Schema Registry CLI and client
==============================================

This repository contains a Command Line Interface (CLI) and a Go client for the REST API of Confluent's Kafka Schema Registry.

[![Build Status](https://travis-ci.org/datamountaineer/schema-registry.svg?branch=master)](https://travis-ci.org/datamountaineer/schema-registry)
[![GoDoc](https://godoc.org/github.com/datamountaineer/schema-registry?status.svg)](https://godoc.org/github.com/datamountaineer/schema-registry)

CLI
---

To install the CLI, assuming a properly setup Go installation, do:

`go get github.com/datamountaineer/schema-registry/schema-registry-cli`

After that, the CLI is found in `$GOPATH/bin/schema-registry-cli`. Running `schema-registry-cli` without arguments gives:

```
A command line interface for the Confluent schema registry

Usage:
  schema-registry-cli [command]

Available Commands:
  add         registers the schema provided through stdin
  exists      checks if the schema provided through stdin exists for the subject
  get         retrieves a schema specified by id or subject
  subjects    lists all registered subjects
  versions    lists all available versions

Flags:
  -h, --help         help for schema-registry-cli
  -e, --url string   schema registry url, overrides SCHEMA_REGISTRY_URL (default "http://localhost:8081")
  -v, --verbose      be verbose

Use "schema-registry-cli [command] --help" for more information about a command.
```

The schema registry url can be configured through the `SCHEMA_REGISTRY_URL` environment variable, and overridden through `--url`. When none is provided, `http://localhost:8081` is used as default.

Client
------

The client package provides a client to deal with the registry from code. It is used by the CLI internally. Usage looks like:

```go
import "github.com/datamountaineer/schema-registry"

client, _ := schemaregistry.NewClient(schemaregistry.DefaultUrl)
client.Subjects()
```

The documentation of the package can be found here: [![GoDoc](https://godoc.org/github.com/datamountaineer/schema-registry?status.svg)](https://godoc.org/github.com/datamountaineer/schema-registry)

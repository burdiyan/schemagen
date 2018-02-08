# Schemagen [![GoDoc](https://godoc.org/github.com/burdiyan/schemagen?status.svg)](https://godoc.org/github.com/burdiyan/schemagen)

This is a tool that fetches Avro schemas from [Confluent Schema Registry](https://github.com/confluentinc/schema-registry) and compiles them to Go code.

Code generation is entirely based on [gogen-avro](https://github.com/alanctgardner/gogen-avro).

NOTICE: Directory `gogen-avro` holds a fork of original gogen-avro that allows to generate Goka codecs as well. It is only for storage, so the same package is also vendored by `dep`. That sucks, but there were no easy way to extend gogen-avro, and I don't want to maintain a separate fork.

## Installation

Right now the only way to install `schemagen` is to build it from source:

```
go install github.com/burdiyan/schemagen/cmd/...
```

## Getting Started

1. Create a file named `.schemagen.yaml` in the root of your project.
2. Specify Schema Registry URL, subjects and versions of the schema you want to download and compile.
3. Run `schemagen` to download the schemas from Schema Registry and compile them.

### Config Example

```
kind: Avro
registry: http://confluent-schema-registry.default.svc.cluster.local:8081
schemas:
  - subject: my-topic-value
    version: latest
    package: country # This is the name of the Go package that will be generated.
  - subject: another-topic-value
    version: "2"
    package: anothertopic
compile: true
outputDir: ./foo
```

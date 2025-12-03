# Examples

This directory contains standalone example programs demonstrating how to use the `go-wxr` package.

## Basic Usage

The `basic` example shows the simplest way to parse a WXR file:

```bash
go run examples/basic/main.go export.xml
```

## Custom Logger

The `custom-logger` example demonstrates how to use a custom logger for debugging:

```bash
go run examples/custom-logger/main.go export.xml
```

## Building Examples

To build an example as a standalone binary:

```bash
go build -o wxr-parser ./examples/basic
./wxr-parser export.xml
```


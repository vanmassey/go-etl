# go-etl

A type-safe, concurrent ETL streaming data engine built in Go. It processes high-volume datasets using zero external dependencies, leveraging channels, generics, and goroutines to process data line-by-line under a constant memory footprint.

## Features
- **Concurrent Processing:** Uses Go's native goroutines to multiplex filtering logic across multiple CPU worker threads simultaneously.
- **Low Memory Footprint:** Streams datasets line-by-line rather than loading entire files into system RAM, preventing memory spikes or out-of-memory crashes on large files.
- **Type-Safe Generics:** Built entirely using Go generics, allowing the pipeline engine to ingest and stream any custom data structure seamlessly.
- **Zero Dependencies:** Relies purely on the Go standard library for network, file, and data processing.

## Project Structure
- `main.go`: The core execution logic and concurrent pipeline architecture.
- `data.csv`: A mock dataset containing system logs for worker pipeline validation.
- `Makefile`: Project automation shortcuts for local building, running, and formatting.
- `.github/workflows/go.yml`: Automated continuous integration configuration to verify compilation on push events.

## Installation and Setup

Ensure you have a modern Go environment installed locally, then execute:

```bash
git clone https://github.com
cd go-etl
make run
```

To compile the application down into a single, high-performance binary executable machine-code file:

```bash
make build
```

## License
Distributed under the permissive [BSD 3-Clause License](LICENSE).

=# go-etl

A high-performance, type-safe concurrent ETL microservice engine built in Go. It ingests high-volume JSON data arrays via HTTP, streaming records line-by-line across multiple concurrent CPU worker threads under a constant, low memory footprint.

## Features
- **Concurrent API Processing:** Multiplexes inbound HTTP request payloads across a pool of concurrent worker goroutines.
- **Low Memory Streaming:** Uses Go channels to pipe incoming data dynamically, preventing memory spikes or out-of-memory crashes on massive datasets.
- **Type-Safe Generics:** Built using Go generics, allowing the internal streaming pipeline to process any structured schema seamlessly.
- **Continuous Integration:** Built-in automated GitHub Actions validation loops for syntax formatting, compilation, and unit tests.
- **Containerized Architecture:** Complete with a multi-stage Docker build, yielding an production-ready runtime footprint of under 20MB.

## Project Structure
- `main.go`: Core execution logic and concurrent HTTP microservice pipeline.
- `main_test.go`: Automated unit testing suite validating mock network API payloads.
- `Makefile`: Project automation shortcuts for local running, compiling, and formatting.
- `Dockerfile`: Multi-stage container instructions for lean production deployment.
- `.github/workflows/go.yml`: Cloud integration runner validating pull requests on push events.

## Getting Started

### Local Setup
Ensure you have a modern Go environment installed locally, then execute:

```bash
git clone https://github.com
cd go-etl
make run
```

### Docker Deployment
To build and run the microservice within an isolated, lightweight container environment:

```bash
docker build -t go-etl .
docker run -p 8080:8080 go-etl
```

## Interfacing with the API

Once the engine is running on port `8080`, transmit a `POST` request to the `/process` endpoint containing a JSON array of records to filter:

```bash
curl -X POST http://localhost:8080/process \
  -H "Content-Type: application/json" \
  -d '[
    {"id": 1, "name": "Alice", "age": 25, "state": "CA"},
    {"id": 2, "name": "Bob", "age": 19, "state": "NY"},
    {"id": 3, "name": "Charlie", "age": 30, "state": "CA"}
  ]'
```

## License
Distributed under the permissive [BSD 3-Clause License](LICENSE).

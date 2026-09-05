# Run the concurrent data pipeline locally
run:
	go run main.go

# Compile the Go code into a single, high-performance binary file
build:
	go build -o go-etl main.go

# Clean up compiled binary executables
clean:
	rm -f go-etl

# Run basic static code formatting checks to keep the repo clean
fmt:
	go fmt ./...

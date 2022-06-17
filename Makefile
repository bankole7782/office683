build:
	rm -rf bin
	mkdir -p bin

	go build -o bin/paelito_maker ./lgcp
	go build -o bin/sites ./sites

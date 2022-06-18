build:
	rm -rf bin
	mkdir -p bin

	go build -o bin/lgcp ./lgcp
	go build -o bin/sites ./sites
	go build -o bin/ssl ./ssl

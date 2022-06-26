build:
	rm -rf tmp/
	rm -rf bin
	mkdir -p bin
	rm -rf office683.tar.xz

	go build -o bin/o6sites ./sites
	go build -o bin/o6ssl ./ssl

	cp services/* bin/

	wget -qO- "https://getbin.io/suyashkumar/ssl-proxy" | tar xvz
	mv ssl-proxy-linux-amd64 bin/
	tar -cJf office683.tar.xz bin/*

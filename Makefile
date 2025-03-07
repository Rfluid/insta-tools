build:
	echo "Compiling for your OS"
	go build -o ./bin/insta-tools

compile:
	echo "Compiling for every OS and Platform"
	GOOS=linux GOARCH=arm go build -o ./bin/insta-tools
	GOOS=linux GOARCH=arm64 go build -o ./bin/insta-tools-arm64
	GOOS=windows GOARCH=386 go build -o ./bin/insta-tools-386
	GOOS=windows GOARCH=arm64 go build -o ./bin/insta-tools-arm64

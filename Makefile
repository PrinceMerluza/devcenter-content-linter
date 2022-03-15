build:
	go build -o bin/content_linter
	zip content-linter.zip ./rules/* ./bin/content_linter

run:
	go run .

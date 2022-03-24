build-blueprint-linter:
	go build -o bin/content_linter
	zip -j content-linter.zip ./bin/content_linter ./blueprint.rule.json

run:
	go run .

test:
	go test -v ./...

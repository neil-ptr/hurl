build:
	go build

run: build
	go run hurl.go

test: build
	go test ./...

# make release version=v1.0.1
.PHONY: release
release:
	@echo "Checking out main branch..."
	git checkout main
	@echo "Tagging with version $(version)..."
	git tag $(version)
	@echo "Pushing tag $(version) to origin..."
	git push origin $(version)
	@echo "Releasing with GoReleaser..."
	goreleaser release --rm-dist

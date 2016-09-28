
GO = go
SOURCE_FILES = $(shell find . -name "*.go" -and -not -name ".git")

build: release/jwt

release/jwt: $(SOURCE_FILES)
	$(GO) build -o release/jwt .

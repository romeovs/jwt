
GO = go
SOURCE_FILES = $(shell ls ./**/*.go *.go)

build: release/jwt

release/jwt: $(SOURCE_FILES)
	$(GO) build -o release/jwt .

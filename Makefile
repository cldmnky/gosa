package = github.com/cldmnky/gosa

.PHONY: release

release:
	mkdir -p release
	GOOS=linux GOARCH=amd64 go build -o release/gosa-cli $(package)

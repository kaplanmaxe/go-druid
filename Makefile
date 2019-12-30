default: test

test:
	@(env GO111MODULE=on go test -cover ./...)

test-ci:
	(env GO111MODULE=on go test -coverprofile=coverage.txt ./...

.PHONY: test
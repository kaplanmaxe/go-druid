default: test

test:
	@(env GO111MODULE=on go test -v -cover ./...)

test-ci:
	(env GO111MODULE=on go test -v -coverprofile=coverage.txt ./...)

.PHONY: test test-ci
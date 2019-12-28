default: test

test:
	@(env GO111MODULE=on go test -cover ./...)

.PHONY: test
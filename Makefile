default: test

test:
	@(go test -mod=vendor -cover ./...)
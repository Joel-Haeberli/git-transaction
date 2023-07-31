VERSION=v0.1.0

tag:
	git tag $(VERSION)
	git push origin $(VERSION)

release: tag
	GOPROXY=proxy.golang.org go list -m github.com/Joel-Haeberli/git-transaction@$(VERSION)
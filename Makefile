# publish package to golang repository
publish:
	@echo "publishing package to golang repository"
	git tag v0.10.0
	git push origin v0.10.0
	@GOPROXY=proxy.golang.org go list -m github.com/merkle3/merkle-sdk-go@v0.10.0
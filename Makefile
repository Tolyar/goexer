PKGS := github.com/Tolyar/goexer

GO := go
GOLINT := golangci-lint

lint:
	 $(GOLINT) run

test: 
	$(GO) test -v -coverprofile cover.out $(PKGS)

cover: | test
	go tool cover -html cover.out
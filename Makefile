update-pkg-cache:
    GOPROXY=https://proxy.golang.org GO111MODULE=on \
    go get github.com/arkemishub/arkectl@v$(VERSION)
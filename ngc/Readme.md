#make help
IMPORTANT: understand make options by this help command

#make build

#make oam
hello-world level sample code right now
generated bins will put into dist folder

#make af
sample code

#make test-unit
run unit test
need to install ginkgo as below:
$ go get github.com/onsi/ginkgo/ginkgo
$ go get github.com/onsi/gomega/...

#make lint
need to install golangci-lint
curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh| sh -s -- -b $(go env GOPATH)/bin v1.21.0

browserd:
	time nice -n 19 go build -p 2 -o dist/$(shell go env GOOS)/browserd github.com/fasionchan/goutils/libs/browser/cmd/browserd

browserd-linux-amd64:
	GOOS=linux GOARCH=amd64 time nice -n 19 go build -p 2 -o dist/linux/amd64/browserd github.com/fasionchan/goutils/libs/browser/cmd/browserd

browserd-linux-arm64:
	GOOS=linux GOARCH=arm64 time nice -n 19 go build -p 2 -o dist/linux/arm64/browserd github.com/fasionchan/goutils/libs/browser/cmd/browserd
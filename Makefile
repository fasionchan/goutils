browserd:
	time nice -n 19 go build -p 2 -o dist/$(shell go env GOOS)/browserd github.com/fasionchan/goutils/libs/browser/cmd/browserd

browserd-linux:
	GOOS=linux time nice -n 19 go build -p 2 -o dist/linux/browserd github.com/fasionchan/goutils/libs/browser/cmd/browserd

browserd-docker: browserd-linux
	docker build -t docker-hub.fasionchan.com/fasionchan/browserd:latest -f libs/browser/cmd/browserd/Dockerfile .

run-browserd: browserd-docker
	docker run -it --rm -p 8080:8080 --shm-size=1g -e BROWSER_headless=new docker-hub.fasionchan.com/fasionchan/browserd:latest pool
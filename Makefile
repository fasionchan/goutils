IMAGE_REGISTRY ?= docker-hub.fasionchan.com/my

DEVTEAM_IMAGE_TAG ?= v0.1

browserd:
	time nice -n 19 go build -p 2 -o dist/$(shell go env GOOS)/browserd github.com/fasionchan/goutils/libs/browser/cmd/browserd

browserd-linux:
	GOOS=linux time nice -n 19 go build -p 2 -o dist/linux/browserd github.com/fasionchan/goutils/libs/browser/cmd/browserd

browserd-docker: browserd-linux
	docker build -t docker-hub.fasionchan.com/fasionchan/browserd:latest -f libs/browser/cmd/browserd/Dockerfile .

run-browserd: browserd-docker
	docker run -it --rm -p 8080:8080 --shm-size=1g -e BROWSER_headless=new docker-hub.fasionchan.com/fasionchan/browserd:latest pool

pi-docker:
	docker build -t docker-hub.fasionchan.com/fasionchan/pi:v0.82.1 -f docker/pi/Dockerfile .

devteam-image-amd64:
	time docker build \
	--platform linux/amd64 \
	-t $(IMAGE_REGISTRY)/devteam:$(DEVTEAM_IMAGE_TAG)-amd64 \
	-f docker/devteam/Dockerfile .

devteam-image-arm64:
	time docker build \
	--platform linux/arm64 \
	-t $(IMAGE_REGISTRY)/devteam:$(DEVTEAM_IMAGE_TAG)-arm64 \
	-f docker/devteam/Dockerfile .

devteam-images: devteam-image-amd64 devteam-image-arm64

devteam-manifest: devteam-images
	docker manifest create $(IMAGE_REGISTRY)/devteam:$(DEVTEAM_IMAGE_TAG) \
	$(IMAGE_REGISTRY)/devteam:$(DEVTEAM_IMAGE_TAG)-amd64 \
	$(IMAGE_REGISTRY)/devteam:$(DEVTEAM_IMAGE_TAG)-arm64

run-devteam:
	docker run -it --rm --env-file .env docker-hub.fasionchan.com/fasionchan/devteam:v0.1

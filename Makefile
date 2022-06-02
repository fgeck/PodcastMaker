
build:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /podcastMaker

docker-build:
	docker build -t floge77/podcastmaker --platform linux/amd64 .

docker-push:
	docker push floge77/podcastmaker

all: build-appservice


build-appservice:
	docker build -t harbor.ym/devops/devops-appservice:v0.0.1 -f docker/Dockerfile.appservice .
	docker push harbor.ym/devops/devops-appservice:v0.0.1


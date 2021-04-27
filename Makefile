all: build-appservice build-artifactory


build-appservice:
	docker build -t harbor.ym/devops/devops-appservice:v0.0.1 -f docker/Dockerfile.appservice .
	docker push harbor.ym/devops/devops-appservice:v0.0.1

build-artifactory:
	docker build -t harbor.ym/devops/devops-artifactory:v0.0.1 -f docker/Dockerfile.artifactory .
	docker push harbor.ym/devops/devops-artifactory:v0.0.1


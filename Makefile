all: build-appservice build-artifactory build-globalconfig build-watch-controller build-workorder build-appservice-controller build-base


build-appservice:
	docker build -t harbor.ym/devops/devops-appservice:v0.0.1 -f docker/Dockerfile.appservice .
	docker push harbor.ym/devops/devops-appservice:v0.0.1

build-artifactory:
	docker build -t harbor.ym/devops/devops-artifactory:v0.0.1 -f docker/Dockerfile.artifactory .
	docker push harbor.ym/devops/devops-artifactory:v0.0.1

build-globalconfig:
	docker build -t harbor.ym/devops/devops-globalconfig:v0.0.1 -f docker/Dockerfile.globalconfig .
	docker push harbor.ym/devops/devops-globalconfig:v0.0.1

build-watch-controller:
	docker build -t harbor.ym/devops/devops-watch-controller:v0.0.1 -f docker/Dockerfile.watch-controller .
	docker push harbor.ym/devops/devops-watch-controller:v0.0.1

build-workorder:
	docker build -t harbor.ym/devops/devops-workorder:v0.0.1 -f docker/Dockerfile.workorder .
	docker push harbor.ym/devops/devops-workorder:v0.0.1


build-appservice-controller:
	docker build -t harbor.ym/devops/devops-appservice-controller:v0.0.1 -f docker/Dockerfile.appservice-controller .
	docker push harbor.ym/devops/devops-appservice-controller:v0.0.1

build-base:
	docker build -t harbor.ym/devops/devops-base:v0.0.1 -f docker/Dockerfile.base .
	docker push harbor.ym/devops/devops-base:v0.0.1
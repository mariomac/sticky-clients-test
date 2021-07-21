images:
	docker build -f deployments/Dockerfile-client --tag="quay.io/mmaciasl/sticky-test-client:latest" .
	docker build -f deployments/Dockerfile-server --tag="quay.io/mmaciasl/sticky-test-server:latest" .
	docker push quay.io/mmaciasl/sticky-test-client:latest
	docker push quay.io/mmaciasl/sticky-test-server:latest

cluster:
	kind create cluster --config deployments/kind-cluster.yml


PHONY: images
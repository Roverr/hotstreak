coverage: ## Deploys coverage to codecov
	go test ./... -coverprofile=coverage.txt -covermode=atomic && \
	curl -s https://codecov.io/bash > deploy.sh && \
	chmod +x ./deploy.sh && \
	./deploy.sh && rm ./deploy.sh
.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.DEFAULT_GOAL := help

.PHONY: build-server
build-server: ## Push the bundle image.
	$(MAKE) -C server docker-build docker-push


.PHONY: integration-tests

integration-tests:
	@set -e; \
	trap 'docker compose -f test/docker-compose.yaml --project-directory test down' EXIT; \
	docker compose -f test/docker-compose.yaml --project-directory test up -d --build; \
	go test -v ./test/...


docker.start:
	docker compose --file ./tests/config/docker-compose.yml  up -d

docker.stop:
	docker compose --file ./tests/config/docker-compose.yml  down

test.unit:
		go test -v ./...

test.integration: docker.start
		go test -count=1 -tags integration -v ./tests/integration/short_url/

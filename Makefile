docker.test.start:
	docker compose --file ./tests/config/docker-compose.yml  up -d

docker.test.stop:
	docker compose --file ./tests/config/docker-compose.yml  down

test.unit:
		go test -v ./...

test.integration: docker.test.start
		go test -count=1 -tags integration -v ./tests/integration/short_url/


docker.start:
	docker compose --file ./docker-compose.yml  up -d

docker.stop:
	docker compose --file ./docker-compose.yml  down

run: docker.start
		go run main.go

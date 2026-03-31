.PHONY: db api worker build clean

db:
	docker-compose up -d postgres

api:
	go run ./cmd/api/

worker:
	go run ./cmd/worker/

build:
	go build -o bin/api ./cmd/api/
	go build -o bin/worker ./cmd/worker/

clean:
	rm -rf bin/

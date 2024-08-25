run:
	go run ./cmd/app/main.go

du:
	docker compose up --build

dd:
	docker compose down  --remove-orphans

test:
	go test ./...
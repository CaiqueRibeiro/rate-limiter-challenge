# run: run application
run:
	@go run src/cmd/main.go

# save-token: saves a token with max request in tyle: `make save-token token=XYZ maxreq=50`
save-token:
	@go run src/cli/main.go --token=$(token) --maxreq=$(maxreq)

docker-up:
	@docker compose up -d

docker-down:
	@docker compose down

test:
	@go test ./...
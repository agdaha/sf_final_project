.PHONY: api_gateway censor-service comments-service news-service all 

api_gateway:
	go build -o ./builds/api_gateway ./api_gateway/cmd/main/main.go

censor-service:
	go build -o ./builds/censor-service ./censor_service/cmd/main/main.go 

comments-service:
	go build -o ./builds/comments-service ./comments_service/cmd/main/main.go 

news-service:
	go build -o ./builds/news-service ./news_service/cmd/main/main.go 

all: api_gateway censor-service comments-service news-service

.PHONY: prune down start_api_gateway rebuild_and_start_api_gateway work testall
start_api_gateway: 
	docker compose up

rebuild_and_start_api_gateway:
	docker compose up --build

prune:
	docker image prune

down:
	docker compose down

work:
	go work init ./api_gateway/ ./censor_service/ ./comments_service/ ./news_service/

testall:
	go test ./api_gateway/... ./censor_service/... ./comments_service/... ./news_service/...


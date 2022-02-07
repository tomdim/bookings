# include variables from the .env file
#!make

help:
	@echo "Please use 'make <target>' where <target> is one of the following:"
	@echo "  up                                 to run compose."
	@echo "  build                              to build and run docker-compose."
	@echo "  down                               to stop docker-compose."
	@echo "  down-volumes                       to stop docker-compose cleaning up volumes."
	@echo "  restart                            to restart docker-compose."
	@echo "  full-restart                       to restart docker-compose cleaning up volumes and rebuilding images."
	@echo "  logs                               to follow logs."
	@echo "  run                                to build the go app locally."
	@echo "  local-build                        to build the go project locally."
	@echo "  test                               to run the tests locally and print coverage."
	@echo "  coverage                           to export coverage report in html format."


up:
	docker-compose up -d

build:
	docker-compose up -d --build

down:
	docker-compose down

down-volumes:
	docker-compose down --volumes --remove-orphans

restart: down up

full-restart: down-volumes build

logs:
	docker-compose logs -f -t

local-build:
	go build -o bookings cmd/web/*.go

run: local-build
	./bookings

test:
	go test -v ./... -cover

coverage:
	go test ./... -coverprofile=coverage.out && go tool cover -html=coverage.out


VERSION := 0.1.0

ROOT_DIR:=$(shell dirname $(realpath $(firstword $(MAKEFILE_LIST))))

ifneq (,${ROOT_DIR}/.version)
    include ${ROOT_DIR}/.version
    export
endif

.PHONY: build clean patch minor major deploy fmt

build: clean
	cd frontend; npm run build
	go build

clean:
	-rm -f ${APPNAME}

patch:
	git aftermerge patch || exit 1

minor:
	git aftermerge minor || exit 1

major:
	git aftermerge major || exit 1

deploy:
	bash deploy/deploy.sh

fmt:
	go fmt ./...

dev: build
	sleep 1 && ./sqlchk server

up:
	docker compose -f deploy/docker-compose.yml up -d --build

logs:
	docker compose -f deploy/docker-compose.yml logs -f
down:
	docker compose -f deploy/docker-compose.yml down
vol:
	docker volume rm deploy_postgres_local

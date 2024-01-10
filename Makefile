.PHONY: build

build:
	cd frontend; npm run build
	go build

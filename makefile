PWD = $(shell pwd)

IMAGE_NAME = expense-tracker-go-image

mock:
	docker run --rm -v "$(PWD)":/src -w /src vektra/mockery --all

clean:
	docker-compose down

remove-image:
	docker rmi $(IMAGE_NAME) || true

build:
	docker-compose build

start:
	docker-compose up -d

rebuild: clean remove-image build start
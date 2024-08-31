PWD = $(shell pwd)

FRONTEND_IMAGE_NAME = expense-tracker-frontend-image
BACKEND_IMAGE_NAME = expense-tracker-backend-image
DOCKER_COMPOSE_FILES = docker-compose.yml

# Function to add frontend compose file if needed
ifeq ($(with_frontend),true)
    DOCKER_COMPOSE_FILES = docker-compose.with.frontend.yml
endif

mock:
	docker run --rm -v "$(PWD)":/src -w /src vektra/mockery:v2.40.1 --all

clean:
	docker-compose -f $(DOCKER_COMPOSE_FILES) down

remove-image:
	docker rmi $(FRONTEND_IMAGE_NAME) $(BACKEND_IMAGE_NAME) || true

build:
	docker-compose -f $(DOCKER_COMPOSE_FILES) build --force-rm

start:
	docker-compose -f $(DOCKER_COMPOSE_FILES) up -d

run: clean remove-image build start

run-with-frontend:
	make run with_frontend=true

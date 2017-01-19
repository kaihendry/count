NAME=count
REPO=hendry/$(NAME)

.PHONY: start stop build sh

all: build

build:
	docker build -t $(REPO) --build-arg COMMIT=$(shell git describe --always) .

start:
	docker run -d --name $(NAME) -p 9000:9000 $(REPO)

stop:
	docker stop $(NAME)
	docker rm $(NAME)

sh:
	docker exec -it $(NAME) /bin/sh

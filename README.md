# Count !

[![Build Status](https://travis-ci.org/kaihendry/count.svg?branch=master)](https://travis-ci.org/kaihendry/count)

[![Go Report Card](https://goreportcard.com/badge/github.com/kaihendry/count)](https://goreportcard.com/report/github.com/kaihendry/count)

Count is a simple Web application to kick the wheels of various deployment
methodologies like that of Github -> Docker Hub -> CoreOS or dokku.

# Building an image with a git commit

	docker build -t hendry/count --build-arg COMMIT=$(git describe --always) .
	docker tag hendry/count hendry/count:purple
	docker push hendry/count:purple

# AWS ECS guide

Inspired by <https://github.com/aws/amazon-ecs-cli/issues/21#issuecomment-235480429>

Ensure <https://github.com/aws/amazon-ecs-cli> is installed:

	ecs-cli --version;wget https://s3.amazonaws.com/amazon-ecs-cli/ecs-cli-linux-amd64-latest -O /usr/local/bin/ecs-cli; ecs-cli --version

The tricky part is setting up the load balancer:

	ecs-cli up --capability-iam --keypair $(whoami)
	./create-task.sh
	./create-load-balancer.sh
	./create-service.sh

Scale service like so:

	ecs-cli scale --size 2 --capability-iam
	ecs-cli compose service scale 2

To roll out an update, change the compose file and:

	ecs-cli compose service up

Bring it all down:

	ecs-cli compose service down
	ecs-cli down --force

Helpers like `./deploy.sh` are to help developers new to Docker build, test & deploy via AWS ECR

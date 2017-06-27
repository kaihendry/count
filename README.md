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

Youtube video that runs through these steps: <https://youtu.be/5YF3qJM3uHA>

Inspired by <https://github.com/aws/amazon-ecs-cli/issues/21#issuecomment-235480429>

Ensure <https://github.com/aws/amazon-ecs-cli> is installed:

	ecs-cli --version;wget https://s3.amazonaws.com/amazon-ecs-cli/ecs-cli-linux-amd64-latest -O /usr/local/bin/ecs-cli; ecs-cli --version

The tricky part is setting up the load balancer:

	ecs-cli configure -r ap-southeast-1 --cluster count -p PROFILE_NAME

	./up.sh

	./create-task.sh

TODO: How do you figure out the VPC the cluster was created in?

	./create-load-balancer.sh

Be aware that I assume all your subnets are in one VPC. https://s.natalian.org/2017-01-19/ecs-setup.png
This can easily be the case if you have some residue.

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

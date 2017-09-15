#!/bin/bash -e

export COMPOSE_PROJECT_NAME=count

echo Please choose the cluster you want to deploy to:
select CLUSTER in $(aws ecs list-clusters | jq -r '.clusterArns[]')
do
	C=$(basename $CLUSTER)
	if ! grep -q $C ~/.ecs/config
	then
		echo Please update the cluster = $C in your ~/.ecs/config
		echo https://github.com/aws/amazon-ecs-cli
		echo Setup your ecs-cli config like so:
		cat <<EOF
[ecs]
cluster                     = $C
region                      = ap-southeast-1
aws_access_key_id           = YOUR_AWS_KEY
aws_secret_access_key       = YOUR_SECRET_KEY
compose-project-name-prefix = ecscompose-
compose-service-name-prefix = ecscompose-service-
cfn-stack-name-prefix       = amazon-ecs-cli-setup-
EOF
exit 1
	fi
	break
done

export COMMIT=$(git describe --always)

ask() {
    # http://djm.me/ask
    local prompt default REPLY

    while true; do

        if [ "${2:-}" = "Y" ]; then
            prompt="Y/n"
            default=Y
        elif [ "${2:-}" = "N" ]; then
            prompt="y/N"
            default=N
        else
            prompt="y/n"
            default=
        fi

        # Ask the question (not using "read -p" as it uses stderr not stdout)
        echo -n "$1 [$prompt] "

        # Read the answer (use /dev/tty in case stdin is redirected from somewhere else)
        read REPLY </dev/tty

        # Default?
        if [ -z "$REPLY" ]; then
            REPLY=$default
        fi

        # Check if the reply is valid
        case "$REPLY" in
            Y*|y*) return 0 ;;
            N*|n*) return 1 ;;
        esac

    done
}

if ask "Do you want to build the $COMPOSE_PROJECT_NAME image?"
then
	docker build -t $COMPOSE_PROJECT_NAME --build-arg COMMIT=$(git describe --always) .
fi

if ask "Do you want to run/test the $COMPOSE_PROJECT_NAME image?"
then
	docker run --name $COMPOSE_PROJECT_NAME --rm -it -p 9000:9000 $COMPOSE_PROJECT_NAME
fi

docker images $COMPOSE_PROJECT_NAME

echo Please choose the private AWS ECR repositry:
select repo in $(aws ecr describe-repositories | jq -r '.repositories[].repositoryUri')
do
	echo Selected $repo
	break
done

if ask "Do you want to deploy the current local $COMPOSE_PROJECT_NAME image COMMIT: ${COMMIT} to $repo?"
then
	docker tag $COMPOSE_PROJECT_NAME $repo:$COMMIT
	eval `aws ecr get-login --region ap-southeast-1`
	docker push $repo
else
echo Would you like to deploy an existing image instead?
unset COMMIT
select COMMIT in $(aws ecr list-images --repository-name $(basename $repo) | jq -r '.imageIds[].imageTag|select(. != null)')
do
	break
done

fi

test "$COMMIT" || exit
echo "You have chosen $COMMIT!"

cat << EOF > docker-compose.yml
version: '2'
services:
  web:
    image: $repo:$COMMIT
    ports:
     - "80:9000"
    memoryReservation: 1100297728
EOF

echo "You should be able to also monitor for deployment issues at"
echo https://ap-southeast-1.console.aws.amazon.com/ecs/home?region=ap-southeast-1#/clusters/$C/services/ecscompose-service-$COMPOSE_PROJECT_NAME/events
#https://ap-southeast-1.console.aws.amazon.com/ecs/home?region=ap-southeast-1#/clusters/$COMPOSE_PROJECT_NAME/services/ecscompose-service-spuul-$C/events"
ecs-cli compose service up

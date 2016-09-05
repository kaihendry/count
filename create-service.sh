export CLUSTER_NAME=default

aws --region ap-southeast-1 ecs create-service --service-name "ecscompose-service-count" \
	--cluster "$CLUSTER_NAME" \
	--task-definition "ecscompose-count" \
	--load-balancers "loadBalancerName=$CLUSTER_NAME,containerName=web,containerPort=9000" \
	--desired-count 1 --deployment-configuration "maximumPercent=100,minimumHealthyPercent=50" --role ecsServiceRole

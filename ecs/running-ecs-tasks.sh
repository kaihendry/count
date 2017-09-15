#!/bin/sh
aws ecs list-clusters | jq -r '.clusterArns[]' | while read clusterdfn
do
	C=$(basename $clusterdfn)
	echo $C
	aws ecs list-tasks --cluster $C | jq -r '.taskArns[]' | while read taskarn
do
	aws ecs describe-tasks --cluster $C --tasks $taskarn | jq -r '.tasks[].taskDefinitionArn' | while read definitionarn
do
	echo $definitionarn
	aws ecs describe-task-definition --task-definition $definitionarn | jq -r '.taskDefinition.containerDefinitions[].image'
done
done
done

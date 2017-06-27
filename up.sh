#!/bin/bash -e

export PROFILE=$(awk -F "=" '/aws_profile/ {print $2}' ~/.ecs/config | tr -d ' ')

select VpcId in $(aws --profile $PROFILE ec2 describe-security-groups | jq -r '.SecurityGroups[].VpcId' | grep -v null | uniq)
do
	break
done

subnets=$(aws --profile "$PROFILE" ec2 describe-subnets | jq -r ".Subnets[] | select(.VpcId==\"$VpcId\") | .SubnetId")

echo ecs-cli up --capability-iam --keypair $(whoami) --vpc $VpcId --subnets $(echo $subnets | tr ' ' ',')

set -e

export CLUSTER_NAME=default

subnet=$(aws ec2 describe-subnets | jq -r '.Subnets[].SubnetId')
aws ec2 describe-security-groups | jq '.SecurityGroups[]| .Description + " " + .GroupId'

select sggroup in $(aws ec2 describe-security-groups | jq -r '.SecurityGroups[].GroupId')
do
	break
done

echo sggroup $sggroup
echo subnet $subnet

aws elb create-load-balancer --load-balancer-name "$CLUSTER_NAME" \
	--listeners Protocol="HTTP,LoadBalancerPort=80,InstanceProtocol=HTTP,InstancePort=80" \
	--subnets $subnet --security-groups $sggroup
	#--subnets subnet-330ebe57 subnet-2b08aa5d

aws elb configure-health-check --load-balancer-name "$CLUSTER_NAME" \
	--health-check Target=HTTP:80/,Interval=10,UnhealthyThreshold=2,HealthyThreshold=2,Timeout=3

set -e

# TODO Get this from ~/.ecs/config ??
export CLUSTER_NAME=count

echo Make sure this is the same as your ECS cluster!
select VpcId in $(aws ec2 describe-security-groups | jq -r '.SecurityGroups[].VpcId' | grep -v null | uniq)
do
	break
done

echo Choose a security group for your ELB
select sggroup in $(aws ec2 describe-security-groups | jq -r ".SecurityGroups[] | select(.VpcId==\"$VpcId\") | .GroupId")
do
	break
done

subnets=$(aws ec2 describe-subnets | jq -r ".Subnets[] | select(.VpcId==\"$VpcId\") | .SubnetId")

echo sggroup $sggroup
echo subnet $subnets

aws elb create-load-balancer --load-balancer-name "$CLUSTER_NAME" \
	--listeners Protocol="HTTP,LoadBalancerPort=80,InstanceProtocol=HTTP,InstancePort=80" \
	--subnets "$subnets" --security-groups $sggroup

aws elb configure-health-check --load-balancer-name "$CLUSTER_NAME" \
	--health-check Target=HTTP:80/,Interval=10,UnhealthyThreshold=2,HealthyThreshold=2,Timeout=3

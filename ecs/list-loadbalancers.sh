export PROFILE=$(awk -F "=" '/aws_profile/ {print $2}' ~/.ecs/config | tr -d ' ')
aws --profile "$PROFILE" elb describe-load-balancers | jq -r '.LoadBalancerDescriptions[] | .CanonicalHostedZoneName + " " + .LoadBalancerName + " " + .Instances[].InstanceId'

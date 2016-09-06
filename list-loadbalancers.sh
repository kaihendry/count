aws elb describe-load-balancers | jq -r '.LoadBalancerDescriptions[] | .CanonicalHostedZoneName + " " + .LoadBalancerName + " " + .Instances[].InstanceId'

https://sam.goserverless.sg/

Have to dig out Cloudfront CNAME from:
https://ap-southeast-1.console.aws.amazon.com/apigateway/main/publish/domain-names?domain=sam.goserverless.sg&region=ap-southeast-1

The [ApiGatewayDomainName7a4a41c73c trick](https://github.com/kaihendry/sam-custom-domain-go/blob/master/template.yaml#L51) did not work.


Does not appear to support static/ files

Tricky / slow to do local development, unless https://github.com/kaihendry/aws-sam-gateway-example structure is used

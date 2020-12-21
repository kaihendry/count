STACK = count
PROFILE = gosls

.PHONY: build deploy validate destroy

# aws --region us-east-1 --profile gosls acm list-certificate
DOMAINNAME = sam.goserverless.sg
ACMCERTIFICATEARN = arn:aws:acm:us-east-1:862322258447:certificate/a7cf55d9-00cb-4672-a721-4b2c2700e23c

deploy: build
	AWS_PROFILE=$(PROFILE) sam deploy --stack-name $(STACK) --parameter-overrides DomainName=$(DOMAINNAME) ACMCertificateArn=$(ACMCERTIFICATEARN)

build:
	CGO_ENABLED=0 sam build

validate:
	aws cloudformation validate-template --template-body file://template.yaml

destroy:
	AWS_PROFILE=$(PROFILE) aws cloudformation delete-stack --stack-name $(STACK)

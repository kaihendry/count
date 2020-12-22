STACK = samcount
PROFILE = gosls

.PHONY: build deploy validate destroy

# aws --region us-east-1 --profile gosls acm list-certificate
DOMAINNAME = sam.goserverless.sg
ACMCERTIFICATEARN = arn:aws:acm:us-east-1:862322258447:certificate/a7cf55d9-00cb-4672-a721-4b2c2700e23c

deploy: build
	AWS_PROFILE=$(PROFILE) sam deploy --stack-name $(STACK) --parameter-overrides DomainName=$(DOMAINNAME) ACMCertificateArn=$(ACMCERTIFICATEARN) --no-confirm-changeset

build:
	CGO_ENABLED=0 sam build

validate:
	AWS_PROFILE=$(PROFILE) aws cloudformation validate-template --template-body file://template.yml

destroy:
	AWS_PROFILE=$(PROFILE) aws cloudformation delete-stack --stack-name $(STACK)

sam-tail-logs:
	# An error occurred (ResourceNotFoundException) if there are no logs yet
	AWS_PROFILE=$(PROFILE) sam logs --name CountFunction --tail

STACK = samcount
PROFILE = gosls
VERSION = $(shell git rev-parse --abbrev-ref HEAD)-$(shell git rev-parse --short HEAD)

.PHONY: build deploy validate destroy

# aws --region us-east-1 --profile gosls acm list-certificate
DOMAINNAME = sam.goserverless.sg
ACMCERTIFICATEARN = arn:aws:acm:us-east-1:862322258447:certificate/a7cf55d9-00cb-4672-a721-4b2c2700e23c

deploy:
	sam build
	AWS_PROFILE=$(PROFILE) sam deploy --stack-name $(STACK) --parameter-overrides DomainName=$(DOMAINNAME) ACMCertificateArn=$(ACMCERTIFICATEARN) --no-confirm-changeset --no-fail-on-empty-changeset

build-CountFunction:
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-X main.Version=${VERSION}" -o ${ARTIFACTS_DIR}/count

validate:
	AWS_PROFILE=$(PROFILE) aws cloudformation validate-template --template-body file://template.yml

destroy:
	AWS_PROFILE=$(PROFILE) aws cloudformation delete-stack --stack-name $(STACK)

sam-tail-logs:
	# An error occurred (ResourceNotFoundException) if there are no logs yet
	AWS_PROFILE=$(PROFILE) sam logs --name CountFunction --tail

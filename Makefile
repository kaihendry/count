STACK = sarmcount
PROFILE = gosls
VERSION = $(shell git rev-parse --abbrev-ref HEAD)-$(shell git rev-parse --short HEAD)

.PHONY: build deploy validate destroy

DOMAINNAME = sarm.goserverless.sg
ACMCERTIFICATEARN = arn:aws:acm:ap-southeast-1:862322258447:certificate/a9363cb2-7413-430c-89d3-b64f634480a7

deploy:
	sam build
	AWS_PROFILE=$(PROFILE) sam deploy --resolve-s3 --stack-name $(STACK) --parameter-overrides DomainName=$(DOMAINNAME) ACMCertificateArn=$(ACMCERTIFICATEARN) --no-confirm-changeset --no-fail-on-empty-changeset --capabilities CAPABILITY_IAM

build-CountFunction:
	CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags "-X main.Version=${VERSION}" -o ${ARTIFACTS_DIR}/bootstrap

validate:
	AWS_PROFILE=$(PROFILE) aws cloudformation validate-template --template-body file://template.yml

destroy:
	AWS_PROFILE=$(PROFILE) aws cloudformation delete-stack --stack-name $(STACK)

sam-tail-logs:
	AWS_PROFILE=$(PROFILE) sam logs --stack-name $(STACK) --tail

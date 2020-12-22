APP = counttesting
RESOURCE_GROUP = AzureFunctionsQuickstart-rg
SUBSCRIPTION_ID =40a590ad-f23f-45d1-86ad-cb952255b437


bin/local/main: main.go
	mkdir -p bin/local
	CGO_ENABLED=0 go build -o bin/local/main main.go

localdev:
	func start --verbose

deploy:
	func azure functionapp publish $(APP)

# https://github.com/Azure/actions-workflow-samples/blob/master/assets/create-secrets-for-GitHub-workflows.md
AZURE_RBAC_CREDENTIALS:
	az ad sp create-for-rbac --name $(APP) --role contributor \
                            --scopes /subscriptions/$(SUBSCRIPTION_ID)/resourceGroups/$(RESOURCE_GROUP) \
                            --sdk-auth

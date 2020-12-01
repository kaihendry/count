1. Create Azure Functions config files.

    ```bash
    func init --custom
    ```

1. Change `main.go` to listen to `FUNCTIONS_CUSTOMHANDLER_PORT`.

    ```go
    func main() { log.Fatal(http.ListenAndServe(":"+os.Getenv("FUNCTIONS_CUSTOMHANDLER_PORT"), routes())) }
    ```

1. Update `host.json` to start the app. We'll set it to the value for production. We'll override this for local development in the next step.

    ```json
    {
        "version": "2.0",
        "logging": {
            "applicationInsights": {
            "samplingSettings": {
                "isEnabled": true,
                "excludedTypes": "Request"
            }
            }
        },
        "extensionBundle": {
            "id": "Microsoft.Azure.Functions.ExtensionBundle",
            "version": "[1.*, 2.0.0)"
        },
        "customHandler": {
            "description": {
                "defaultExecutablePath": "bin/production/main"
            },
            "enableForwardingHttpRequest": true
        },
        "extensions": {
            "http": {
                "routePrefix": ""
            }
        }
    }
    ```

    Also note that we added a config to change the default route prefix from `/api` to `/`.

1. Update `local.settings.json` to override the executable to start when running locally. We do this because we're running a different OS or architecture locally than in production. If running locally on Windows, this might be `"bin\\local\\main.exe"` instead.

    ```json
    {
        "IsEncrypted": false,
        "Values": {
            "FUNCTIONS_WORKER_RUNTIME": "custom",
            "AzureWebJobsStorage": "",
            "AzureFunctionsJobHost__customHandler__description__defaultExecutablePath": "bin/local/main"
        }
    }
    ```

1. Add an HTTP function.

    ```bash
    func new --name httptrigger --template "HTTP Trigger"
    ```

1. Update the function to trigger on all routes and methods. Also change auth to anonymous.

    ```json
    {
    "bindings": [
        {
            "authLevel": "anonymous",
            "type": "httpTrigger",
            "direction": "in",
            "name": "req",
            "route": "{*path}"
        },
        {
            "type": "http",
            "direction": "out",
            "name": "res"
        }
    ]
    }
    ```

1. Build the app locally and run it.

    ```bash
    # linux/macOS
    go build -o bin/local/main main.go

    # Windows
    go build -o bin/local/main.exe main.go

    func start --verbose
    ```

1. (Optional) Add a `.funcignore` file so we don't deploy the local binary too.

    ```
    bin/local
    ```

1. Build the app for production.

    ```bash
    GOOS=linux GOARCH=amd64 go build -o bin/production/main main.go
    ```

1. Deploy the app.

    ```bash
    az login
    az group create --name AzureFunctionsQuickstart-rg --location westeurope
    az storage account create --name <STORAGE_NAME> --location westeurope --resource-group AzureFunctionsQuickstart-rg --sku Standard_LRS
    az functionapp create --resource-group AzureFunctionsQuickstart-rg --consumption-plan-location westeurope --runtime custom --functions-version 3 --os-type Linux --name <APP_NAME> --storage-account <STORAGE_NAME>
    func azure functionapp publish <APP_NAME>
    ```
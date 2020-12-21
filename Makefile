PROJECTID = idiotbox
SERVICE = count

# Takes the most time
build:
	gcloud builds submit --tag gcr.io/$(PROJECTID)/$(SERVICE)

# https://cloud.google.com/sdk/gcloud/reference/run/deploy
deploy:
	gcloud run deploy $(SERVICE) --image gcr.io/$(PROJECTID)/$(SERVICE) --platform managed --region asia-southeast1 --allow-unauthenticated

localbuild:
	docker build -t count .

localrun:
	docker run -p 3000:8080 count

service.yaml:
	gcloud run services describe count --format yaml --platform managed > service.yaml

# https://cloud.google.com/run/docs/mapping-custom-domains#command-line
domainsetup:
	gcloud beta run domain-mappings create --service $(SERVICE) --domain $(SERVICE).dabase.com --platform managed

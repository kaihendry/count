REGION=asia-southeast1
RUNTIME=go119

deploy:
	gcloud functions deploy Countpage --runtime $(RUNTIME) --trigger-http --region $(REGION) --allow-unauthenticated

logs:
	gcloud functions logs read --region $(REGION)

describe:
	gcloud functions describe Countpage --region $(REGION)

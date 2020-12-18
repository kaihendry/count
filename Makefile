PROJECTID = idiotbox
SERVICE = count

build:
	gcloud builds submit --tag gcr.io/$(PROJECTID)/$(SERVICE)

deploy:
	gcloud run deploy --image gcr.io/$(PROJECTID)/$(SERVICE) --platform managed --region asia-southeast1 --allow-unauthenticated

localbuild:
	docker build -t count .

localrun:
	docker run -p 3000:8080 count

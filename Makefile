deploy:
	time gcloud functions deploy api --entry-point Countpage --runtime go111 --trigger-http

logs:
	gcloud functions logs read

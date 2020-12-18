deploy:
	# Not available in Singapore :( https://cloud.google.com/functions/docs/locations
	gcloud functions deploy api --entry-point Countpage --runtime go113 --trigger-http --region asia-east2

logs:
	gcloud functions logs read

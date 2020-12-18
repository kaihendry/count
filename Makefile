deploy:
	# Not available in Singapore :( https://cloud.google.com/functions/docs/locations
	gcloud functions deploy Countpage --runtime go113 --trigger-http --region asia-east2 --allow-unauthenticated

logs:
	gcloud functions logs read --region asia-east2

describe:
	gcloud functions describe Countpage --region asia-east2

# Count !

[![Build Status](https://travis-ci.org/kaihendry/count.svg?branch=master)](https://travis-ci.org/kaihendry/count)

[![Go Report Card](https://goreportcard.com/badge/github.com/kaihendry/count)](https://goreportcard.com/report/github.com/kaihendry/count)

Count is a simple Web application to kick the wheels of various deployment
methodologies like that of Github -> Docker Hub -> CoreOS or dokku.

* [How to deploy this without any downtime to my users with AWS?](https://www.youtube.com/watch?v=onTnyvrHggo)

# Building an image with a git commit

	docker build -t hendry/count --build-arg COMMIT=$(git describe --always) .
	docker tag hendry/count hendry/count:purple
	docker push hendry/count:purple

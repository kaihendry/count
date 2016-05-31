	docker build -t count github.com/kaihendry/count
	docker run -it -p 9000:9000 count

Please refer to [count.service](count.service)

# count.service

Remove `User=core` if you don't use CoreOS.

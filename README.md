	docker build -t count github.com/kaihendry/count
	docker run -it -p 9000:9000 count

# count.service

Remove `User=core` if you don't use CoreOS.

# Restart issue

I noticed without `ExecStop=/usr/bin/docker stop count` it would take 1m30s to restart!

With `docker stop` it takes 10s, but even that seems very slow. What am I missing?

	May 31 15:26:23 teefour systemd[1]: Stopping Count...
	May 31 15:26:33 teefour docker[492]: time="2016-05-31T15:26:33.137235364+08:00" level=info msg="Container d8b8df91f16988469a59c2590b41653c9436422da95b3a40968b59f1d3f74126 failed to exit within 10 seconds of signal 15 - using the force"

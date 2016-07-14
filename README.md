# Count !

[![Go Report Card](https://goreportcard.com/badge/github.com/kaihendry/count)](https://goreportcard.com/report/github.com/kaihendry/count)

Count is a simple Web application to kick the wheels of various deployment
methodologies like that of Github -> Docker Hub -> CoreOS or dokku.

* [How to deploy this without any downtime to my users?](https://youtu.be/04np_kwmv_g)
* How to roll back?
* How to add SSL? A: Via [caddy](https://caddyserver.com/)
* How do I monitor this? A: `journalctl -u count -f`

# [count.service](count.service)

Remove `User=core` if you don't use CoreOS.

Notice I rely on <https://caddyserver.com/> to proxy the site.

# Example caddy configuration

	count.dabase.com {
		tls foo@example.com
		proxy / count:9000
	}

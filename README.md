# Luzifer / restis

`restis`, composed from `REST` and `Redis`, is a simple HTTP API to get / set / delete Redis keys.

## Q&A

- **Why?**  
  Previously I used a kyototycoon to store simple key-value-data through a REST interface but needed something not using a proprietary database format below it to have it running statelessly in a cluster. Searching for some alternative I gave up and invested 20 minutes in putting together this.
- **Is there any security built in?**  
  No. And it never will be. Don't expose it without securing it for example using [nginx-sso](https://github.com/Luzifer/nginx-sso) or Gatekeeper or whatever. This is just a simple API, you're responsible for the rest.

## Usage

```console
# restis --help
Usage of ./restis:
      --disable-cors               Disable setting CORS headers for all requests
      --listen string              Port/IP to listen on (default ":3000")
      --log-level string           Log level (debug, info, warn, error, fatal) (default "info")
      --redis-conn-string string   Connection string for redis (default "redis://localhost:6379/0")
      --redis-key-prefix string    Prefix to prepend to keys (will be prepended without delimiter!)
      --version                    Prints current version and exits
```

```console
# curl -fX GET localhost:3000/mykey
curl: (22) The requested URL returned error: 404

# echo "mycontent" | curl -fX PUT --data-binary @- localhost:3000/mykey

# curl -fX GET localhost:3000/mykey
mycontent

# curl -fX DELETE localhost:3000/mykey

# curl -fX GET localhost:3000/mykey
curl: (22) The requested URL returned error: 404

# echo "mycontent" | curl -fX PUT --data-binary @- localhost:3000/mykey\?expire=30s

# curl -fX GET localhost:3000/mykey
mycontent

# sleep 30; curl -fX GET localhost:3000/mykey
curl: (22) The requested URL returned error: 404
```

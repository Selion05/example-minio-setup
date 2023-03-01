# Example minio setup for local dev environment as s3 replacement

I needed a local s3 api for development.

## Requirements
- s3 access from host, container and browser
- pre-signed urls (cors, mixed content)
- working minio frontend 

Note that I use a self-signed certificate `https://*.docker.localhost` in my traefik setup, 
so I needed https otherwise I would  get [mixed content](https://developer.mozilla.org/en-US/docs/Web/Security/Mixed_content) errors.
If you don't just add `minio` and `minioui` to your `/etc/hosts` file on your host and adapt urls accordingly


## Setup

`docker compose up`

[login](https://minioui.docker.localhost/login) 
user: root 
password: password

[create bucket](https://minioui.docker.localhost/buckets/add-bucket) `test`

[create access key](https://minioui.docker.localhost/access-keys/new-account) and use as credentials in `main.go`

`go mod tidy`


`go run main.go` Example how to upload and download a file with pre-signed urls



# ip-check

A simple website for checking IPv4 and IPv6 of client.

Demo: [https://ipcheck.thomasdang.ca](https://ipcheck.thomasdang.ca)

## Features

Shows the IPv4 or IPv6 of a client. Can be curl'd (or other non-interactive request method) to receive specific information. Commands are displayed on the web page. 

Has a simple javascript checker that determines if IPv4 and IPv6 are each available and displays the respective IP. 

## Usage

You will want to update the template in `golang-webserver/src/handlers/templates` with your URLs.

Compile using Go 1.23 using `go build main.go`

The server automatically starts on port 8080, you will probably want a systemd file to automatically restart this and do any necessary logging, here is a minimal working file:

```
[Unit]
Description=IP Check Web Server
After=network.target
Wants=network.target

[Service]
ExecStart=/usr/local/bin/ip-check
Restart=always
RestartSec=5
WorkingDirectory=/usr/local/bin
StandardOutput=append:/var/log/ip-check.log
StandardError=append:/var/log/ip-check.log

[Install]
WantedBy=multi-user.target
```

I run an nginx instance in a container hosted with the file. There are many other ways to deploy this more securely or efficiently.

You'll want to configure a subdomain with just A records and one with just AAAA records, and one with both. This will allow you to test the connectivity of both and provide the option to resolve one or the other.

## To-do

- Docker Compose deployment
- Use environment variables for the URLs
- Add whois data for ASN tracking etc.
- Add a "share" style link for sharing the basic information to administrators from users

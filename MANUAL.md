## Installing (from source)

Clone or otherwise download the Popple source code, then:

```console
$ go test ./...
$ go build ./cmd/discord/popple
$ go install ./cmd/discord/popple
```

Alternatively:

```console
$ go install github.com/connorkuehl/popple@latest
```

Or, instead of `latest`, any release tag may be used.

## Installing (from package manager)

Popple is not currently packaged for any software distribution. The only
way to build or install this software is to do so from source or via
a container image.

## Deploying (systemd)

Create a new user, `popple`, and then create a systemd unit for it:

Example systemd unit:

```
/etc/systemd/system/popple.service
```

```systemd
[Unit]
Description=Popple, a karma bot for Discord
After=network-online.target
Wants=network-online.target

[Service]
Type=simple
User=popple
Group=popple
ExecStart=/usr/local/bin/popple -token-file /srv/popple/popple.token -database /srv/popple/popple.sqlite
Restart=on-failure
RestartSec=5

[Install]
WantedBy=multi-user.target
```

```console
$ systemctl enable popple
$ systemctl start popple
```

## Deploying (container/podman/docker)

```console
$ docker pull docker.io/conkue/popple:latest
$ docker run -d --restart always \
	--name popple \
	-v /srv:/data \
	popple:latest -token-file popple.token -database popple.sqlite
```

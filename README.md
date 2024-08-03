# Rcee Isolate

Remote Code Execution Engine

## Deploy

Tested in Ubuntu 24.04

Install isolate from github release <https://github.com/soufrabi/isolate/releases/tag/unstable/>

Build go binary
```sh
go build
```

Place the `rcee-isolate` in a directory in `PATH`

Place the [systemd service](./systemd/rcee-isolate.service) in `/etc/systemd/system` directory

Instruct systemd to re-read its configuration files and register `rcee-isolate.service`
```sh
sudo systemctl daemon-reload
```

Enable and start the systemd service
```sh
sudo systemctl enable --now rcee-isolate.service
```

The app should be running at PORT `8080`

Check if app is running using
```sh
curl http://localhost:8080/ping
```

Check systemd logs using
```sh
sudo journalctl -xeu rcee-isolate.service
```

## Environment variables

- GO_ENV
: Set to `production` in production environment
- PORT
: Change PORT of application

## build

```bash
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -tags=windows -ldflags "-s -w -H=windowsgui" -o cron_bilitools.exe main.go cmd_windows.go
```

## run

windows：

```cmd
cron_bilitools.exe -config=./config/config.json -time=08:08:08 -start=false
```

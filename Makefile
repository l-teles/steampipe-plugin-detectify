install:
	go build -o ~/.steampipe/plugins/hub.steampipe.io/plugins/l-teles/detectify@latest/steampipe-plugin-detectify.plugin *.go
	
# LOCAL DEVELOPMENT
# install:
# 	go build -o ~/.steampipe/plugins/local/detectify/detectify.plugin *.go

# Linux ARM64
# install:
# 	GOOS=linux GOARCH=arm64 go build -o ~/.steampipe/plugins/local/detectify/detectify-linux.plugin *.go
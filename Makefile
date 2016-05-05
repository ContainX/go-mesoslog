VERSION = 0.2

GO_FMT = gofmt -s -w -l .
GO_XC = goxc -os="linux darwin windows" -tasks-="rmbin"

GOXC_FILE = .goxc.local.json

all: deps compile

compile: goxc

goxc:
	$(shell echo '{\n "ConfigVersion": "0.9",\n "PackageVersion": "$(VERSION)",' > $(GOXC_FILE))
	$(shell echo ' "TaskSettings": {' >> $(GOXC_FILE))
	$(shell echo '  "bintray": {\n   "apikey": "$(BINTRAY_APIKEY)"' >> $(GOXC_FILE))
	$(shell echo '  },' >> $(GOXC_FILE))
	$(shell echo '  "publish-github": {' >> $(GOXC_FILE))
	$(shell echo '     "apikey": "$(GITHUB_APIKEY)",' >> $(GOXC_FILE))
	$(shell echo '     "body": "",' >> $(GOXC_FILE))
	$(shell echo '     "include": "*.zip,*.tar.gz,*.deb,go-mesoslog_$(VERSION)_linux_amd64-bin"' >> $(GOXC_FILE))
	$(shell echo '  }\n } \n}' >> $(GOXC_FILE))
	$(GO_XC)

deps:
	go get

format:
	$(GO_FMT)

bintray:
	$(GO_XC) bintray

github:
	$(GO_XC) publish-github

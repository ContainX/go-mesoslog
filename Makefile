VERSION = 0.1

GO_FMT = gofmt -w
GO_XC = goxc -os="linux darwin windows freebsd openbsd"

GOXC_FILE = .goxc.local.json

all: deps

compile: goxc

goxc:
	$(shell echo '{\n "ArtifactsDest": "build",\n "ConfigVersion": "0.9",' > $(GOXC_FILE))
	$(shell echo ' "PackageVersion": "$(VERSION)",\n "TaskSettings": {' >> $(GOXC_FILE))
	$(shell echo '  "bintray": {\n   "user": "gondor",\n   "apikey": "$BINTRAY_APIKEY",\n   "package": "go-mesoslog",' >> $(GOXC_FILE))
	$(shell echo '   "repository": "utils",\n   "subject": "pacesys"' >> $(GOXC_FILE))
	$(shell echo '  }\n }\n}' >> $(GOXC_FILE))
	$(GO_XC) 

deps:
	go get

bintray:
	$(GO_XC) bintray
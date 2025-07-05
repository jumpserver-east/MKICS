NAME=evobot
BUILDDIR=build

BASEPATH := $(shell pwd)
BRANCH := $(shell git symbolic-ref HEAD 2>/dev/null | cut -d"/" -f 3)
BUILD := $(shell git rev-parse --short HEAD)

VERSION ?= $(BRANCH)-$(BUILD)
BuildTime:= $(shell date -u '+%Y-%m-%d %I:%M:%S%p')
VERSION ?= $(BRANCH)-$(BUILD)
TARGETARCH ?= amd64

UIDIR=frontend

.PHONY: docker
docker:
	@echo "build docker images"
	docker build --platform=linux/amd64 -f Dockerfile -t jumpserver-east/evobot:$(VERSION) .

.PHONY: clean
clean:
	rm -rf $(BUILDDIR)
	rm -rf $(UIDIR)/dist

.PHONY: package
package: docker
	@echo "▶ 开始打包版本: $(VERSION)"
	@rm -rf $(BUILDDIR)
	@mkdir -p $(BUILDDIR)/evobot/
	
	cp installer/docker-compose.yml $(BUILDDIR)/evobot/
	cp cmd/conf/config-example.yaml $(BUILDDIR)/evobot/
	
	@echo "▶ 版本号为: $(VERSION)"
	sed -i '' 's|^\([[:space:]]*image:[[:space:]]*evobot:\).*|\1$(VERSION)|' build/evobot/docker-compose.yml

	docker save -o $(BUILDDIR)/evobot/evobot-$(VERSION)-docker-image.tar jumpserver-east/evobot:$(VERSION)
	
	tar -czvf $(BUILDDIR)/evobot-installer-$(VERSION)-package.tar.gz -C $(BUILDDIR)/evobot .
	
	@echo "✓ 打包完成: $(BUILDDIR)/evobot-installer-$(VERSION)-package.tar.gz"

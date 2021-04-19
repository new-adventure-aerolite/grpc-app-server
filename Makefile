OUT_DIR=.build
APP=app-server
BRANCH=$(shell git branch | grep \* | cut -d ' ' -f2)
COMMIT=$(shell git rev-parse --short HEAD)
DATE=$(shell date +'%Y%m%d%H%M%S')

all: clean directories $(TARGET)
.PHONY: all

$(OUT_DIR)/$(APP):
	go build -ldflags \
	"-X main.buildTime=`date -u '+%Y-%m-%d_%I:%M:%S%p'` \
	-X main.branch=$(BRANCH) \
	-X main.commit=$(COMMIT)" -o $@ -mod=vendor ./cmd/$(APP)

run:
	go run ./cmd/$(APP)/main.go
.PHONY: run

docker:
	docker build -t app-server:latest .
.PHONY: docker

test:
	go test -v ./... -cover -count=1
.PHONY: test

all: clean directories $(OUT_DIR)/$(APP)
.PHONY: all

directories:
	mkdir -p $(OUT_DIR)
.PHONY: directories

clean:
	rm -rf $(OUT_DIR)
.PHONY: clean

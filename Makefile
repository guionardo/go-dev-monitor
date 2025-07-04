.DEFAULT_TARGET: help

.PHONY: help
help: ## Display this help
    @grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(firstword $(MAKEFILE_LIST)) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: json
json: ## Build json marshalling functions - execute after update the models
	go get github.com/mailru/easyjson
	go install github.com/mailru/easyjson/...@latest
	easyjson -all internal/repository/repo.go
	easyjson -all internal/api/request_response.go

.PHONY: frontend
frontend: ## Build distribution files for frontend into server static dir
	cd frontend ; bun run lint ; bun run format ; bun run build

build: frontend
	go build cmd/server

gdm:
	go build -ldflags="-X internal.Version=1.0.0 -X internal.Build=$(date +%Y%m%d%H%M%S)" -o gdm cmd/gdm/main.go

install_service: gdm
	mv gdm ~/go/bin/gdm
	@echo [Unit] > go-dev-monitor.service
	@echo Description=go-dev-monitor service >> go-dev-monitor.service
	@echo After=network-online.target >> go-dev-monitor.service
	@echo [Service] >> go-dev-monitor.service
	@echo Environment=HOME=$(HOME) >> go-dev-monitor.service
	@echo ExecStart=$(HOME)/go/bin/gdm serve  >> go-dev-monitor.service
	@echo [Install] >> go-dev-monitor.service
	@echo WantedBy=multi-user.target >> go-dev-monitor.service
	@sudo cp go-dev-monitor.service /etc/systemd/system/go-dev-monitor.service
	@sudo systemctl daemon-reload
	@echo Run service with systemctl start go-dev-monitor

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
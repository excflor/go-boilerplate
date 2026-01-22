.PHONY: help
help:
	@echo 'Usage:'
	@grep -hE '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2}'

.PHONY: tidy
tidy: ## format code and tidy modfiles
	@echo 'Tidying and formatting Go mod files...'
	go fmt ./...
	go mod tidy

.PHONY: watch
watch: ## run the application with live reloading (Air)
	@if command -v air > /dev/null; then \
	    air; \
	else \
	    read -p "Air is not installed. Do you want to install it? [y/N] " choice; \
	    if [ "$$choice" = "y" ] || [ "$$choice" = "Y" ]; then \
	        go install github.com/air-verse/air@latest; \
	        air; \
	    fi; \
	fi
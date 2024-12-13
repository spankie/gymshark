include .env

IMAGE_NAME=gymshark-api
IMAGE_VERSION=latest
GCP_PROJECT_ID=gymshark-interview
GCR_IMAGE=gcr.io/$(GCP_PROJECT_ID)/$(IMAGE_NAME):$(IMAGE_VERSION)

# Build the application
all: build

build:
	@echo "Building..."
	@go build -o main cmd/api/main.go

build-image:
	@docker build -t $(IMAGE_NAME):$(IMAGE_VERSION) .
	@docker tag $(IMAGE_NAME):$(IMAGE_VERSION) $(GCR_IMAGE)
	@docker push $(GCR_IMAGE)

deploy:
	@gcloud run deploy gymshark-api --image $(GCR_IMAGE)

# Run the application
run:
	@go run cmd/api/main.go

# Create DB container
docker-run:
	@if docker compose up 2>/dev/null; then \
		: ; \
	else \
		echo "Falling back to Docker Compose V1"; \
		docker-compose up; \
	fi

# Shutdown DB container
docker-down:
	@if docker compose down 2>/dev/null; then \
		: ; \
	else \
		echo "Falling back to Docker Compose V1"; \
		docker-compose down; \
	fi

lint:
	@echo "Linting..."
	@golangci-lint run

# Test the application
test:
	@echo "Testing..."
	@go test ./... -v


# Integrations Tests for the application
itest:
	@echo "Running integration tests..."
	@go test -tags="integration" -v ./...


# Clean the binary
clean:
	@echo "Cleaning..."
	@rm -f main

# Live Reload

watch:
	@if command -v air > /dev/null; then \
            air; \
            echo "Watching...";\
        else \
            read -p "Go's 'air' is not installed on your machine. Do you want to install it? [Y/n] " choice; \
            if [ "$$choice" != "n" ] && [ "$$choice" != "N" ]; then \
                go install github.com/air-verse/air@latest; \
                air; \
                echo "Watching...";\
            else \
                echo "You chose not to install air. Exiting..."; \
                exit 1; \
            fi; \
        fi


.PHONY: all build run test clean watch

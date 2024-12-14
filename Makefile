include .env

BACKEND_IMAGE_NAME=gymshark-api
BACKEND_IMAGE_VERSION=latest
FRONTEND_IMAGE_NAME=gymshark-frontend
FRONTEND_IMAGE_VERSION=latest
REGION=us-east-1
FRONTEND_IMAGE=$(ACCOUNT_ID).dkr.ecr.$(REGION).amazonaws.com/$(FRONTEND_IMAGE_NAME):$(FRONTEND_IMAGE_VERSION)
BACKEND_IMAGE=$(ACCOUNT_ID).dkr.ecr.$(REGION).amazonaws.com/$(BACKEND_IMAGE_NAME):$(BACKEND_IMAGE_VERSION)
FRONTEND_CLUSTER=gymshark-frontend
BACKEND_CLUSTER=gymshark-backend
BACKEND_URL=http://gymshark-lb-975620203.us-east-1.elb.amazonaws.com

# Build the application
all: build

run-frontend: build-frontend
	@docker run -e VITE_API_BASE_URL=http://localhost:8080 \
		-p 5173:80 $(FRONTEND_IMAGE)

build-frontend:
	@docker buildx build --platform=linux/amd64 --build-arg VITE_API_BASE_URL=$(BACKEND_URL) -t $(FRONTEND_IMAGE) ./frontend
	@docker push $(FRONTEND_IMAGE)

build:
	@echo "Building..."
	@go build -o main cmd/api/main.go

# Run the application
run:
	@go run cmd/api/main.go

build-image:
	@docker buildx build --platform=linux/amd64 -t $(BACKEND_IMAGE) .
	@docker push $(BACKEND_IMAGE)

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

# Clean the binary
clean:
	@echo "Cleaning..."
	@rm -f main
	@rm -f cmd/lambda/bootstrap

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


# create ecr repository
create-repo:
	@aws ecr create-repository --repository-name $(BACKEND_IMAGE_NAME)
	@aws ecr create-repository --repository-name $(FRONTEND_IMAGE_NAME)

auth-docker:
	@aws ecr get-login-password --region $(REGION) | docker login --username AWS --password-stdin $(ACCOUNT_ID).dkr.ecr.$(REGION).amazonaws.com

build-lambda:
	@GOOS=linux GOARCH=amd64 go build -tags lambda.norpc -o cmd/lambda/bootstrap cmd/lambda/main.go
	@zip -j cmd/lambda/ShippingPacks.zip cmd/lambda/bootstrap

deploy-lambda: build-lambda
	@aws lambda create-function --function-name ShippingPacks \
		--runtime provided.al2023 --handler bootstrap \
		--architectures x86_64 \
		--role $(AWS_LAMBDA_ARN) \
		--zip-file fileb://cmd/lambda/ShippingPacks.zip

frontend-cluster:
	@aws ecs create-cluster --cluster-name $(FRONTEND_CLUSTER)

.PHONY: all build run test clean watch

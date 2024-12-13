# Non Technical Overview
This is a web server that handles calculation of shipping packaging based on the order provided and available packs

## 1. Creating an Order:
 
 When someone sends a request to create an order, the server first checks if the orders items is valid (greater than zero).
 If the order is empty or there's an error in processing it, the server responds with an error.
 If everything is fine, the server calculates the optimal number of packs necessary to ship the order and saves
 the order alongside the shipping details to a database and confirms the successful creation
 by returning the id of the order. This ensures that only valid orders are stored and helps maintain the integrity of the data.


## 2. Retrieving an Order:

When someone requests a specific order by its ID, the server retrieves the order from the database.
If the order is not found (e.g., the ID does not exist), the server logs a debug message and responds with an error that says not found.
If the order is successfully retrieved, the server responds with the order content.
This ensures that only valid order are stored and retrieved, maintaining the integrity of the data
and providing appropriate responses for different scenarios.

---

# How to Run the Code

This project uses a Makefile to simplify various tasks such as building,
running, and deploying the application to AWS. Below are the steps to run the code:

1. **Build the Application**:
   - To compile the Go application, run:
     ```sh
     make
     ```
   - This command will build the application and create an executable named `main`.

2. **Run the Application**:
   - Prerequisites:
     - Ensure you have a PostgreSQL database running on your machine or use the provided Docker Compose file to start a database container.
     - Set the environment variables to run the application and connect to the database.
      ```env
        export GYMSHARK_PORT=8000
        export GYMSHARK_APP_ENV=local
        export GYMSHARK_DB_HOST=localhost
        export GYMSHARK_DB_PORT=5432
        export GYMSHARK_DB_USERNAME=changeusername
        export GYMSHARK_DB_ROOT_PASSWORD=changepassword
        export GYMSHARK_DB_NAME=changedatabasename
        export GYMSHARK_LOG_LEVEL=debug
      ```
   - If you don't have a postgres instance running on your machine,
      you can use the provided docker-compose file to start a postgres container.
      To start the database container, use:
     ```sh
     make docker-run
     ```
     This command will attempt to start the database container using Docker Compose. If Docker Compose V2 is not available, it will fall back to Docker Compose V1.

   - Once the database container is running, you can start the application using:
     ```sh
     make run
     ```
   - This command will execute the Go application located at [`cmd/api/main.go`](./cmd/api/main.go).

3. **Build and Push Docker Image**:
   - To build a Docker image and push it to a container registry, use:
     ```sh
     make build-image
     ```
   - Ensure that the `IMAGE_NAME`, `IMAGE_VERSION` and `GCR_IMAGE` variables are set appropriately in your environment.

4. **Deploy the Application**:
   - To deploy the application to Google Cloud Run, use:
     ```sh
     make deploy
     ```
   - This command will deploy the Docker image to Google Cloud Run using the image specified by `GCR_IMAGE`.

5. **Shutdown Database Container**:
   - To stop the database container, use:
     ```sh
     make docker-down
     ```
   - This command will attempt to stop the database container using Docker Compose. If Docker Compose V2 is not available, it will fall back to Docker Compose V1.

6. **Lint the Code**:
   - To lint the code, use:
     ```sh
     make lint
     ```
   - This command will run the linting process to check for code quality issues.

By using these commands, you can easily manage the build, run, and deployment processes for the application.

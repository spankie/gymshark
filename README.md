# Non Technical Overview
This is a web application that handles calculation of shipping packaging based on the order provided and available packs

## Architecture

The application consists of three main components: frontend, backend, and database.

- **Backend**: Built with **Golang**, deployed on **AWS ECS Fargate** (serverless).
- **Frontend**: Developed using **React** with **TypeScript**, also deployed on **AWS ECS Fargate**.
- **Database**: A **PostgreSQL** instance hosted on **AWS RDS**.

Both the frontend and backend are behind a **Load Balancer**, which routes traffic to the ECS instances running the application.

# Features

## 1. Creating an Order

- On the **frontend**, users input the order quantity and click **Add Order**.  
- A **request** is sent to the server, which validates the order:  
  - Ensures the order quantity is greater than zero.  
  - If the order is invalid (empty or erroneous), the server responds with an **error**.  
- For valid orders:  
  - The server calculates the **optimal number of packs** required to fulfill the order.  
  - The order, along with shipping details, is saved to the **database**.  
  - A **confirmation** with order details is returned to the frontend.  

This process ensures that only valid orders are stored, maintaining data integrity.

## 2. Retrieving Orders

- The **main page** displays a table with a list of all orders.  
- Each order includes **packaging details** calculated based on predefined criteria for minimum items and optimal packaging.  

---

# How to Run the Code

This project uses a Makefile to simplify various tasks such as building,
running. Below are the steps to run the code:

1. **Run the backend server**:
   - Prerequisites:
     - Ensure you have a PostgreSQL database running on your machine or use the provided Docker Compose file to start a database container.
     - Set the environment variables to run the application and connect to the database.
      ```env
        export GYMSHARK_PORT=8000
        export GYMSHARK_DB_HOST=localhost
        export GYMSHARK_DB_PORT=5432
        export GYMSHARK_DB_USERNAME=changeusername
        export GYMSHARK_DB_PASSWORD=changepassword
        export GYMSHARK_DB_NAME=changedatabasename
        export GYMSHARK_ENABLE_DB_SSL=false
        export GYMSHARK_LOG_LEVEL=debug
        export GYMSHARK_FRONTEND_URL=http://localhost:8080
      ```
   - If you don't have a postgres instance running on your machine,
      you can use the provided docker-compose file to start a postgres container.
      To start the database container, use:
     ```sh
     make docker-run
     ```
     This command will attempt to start the database container using Docker Compose. If Docker Compose V2 is not available, it will fall back to Docker Compose V1.

   - Once the database container is running, you can start the backend server using:
     ```sh
     make run
     ```
   - This command will execute the Go application located at [`cmd/api/main.go`](./cmd/api/main.go).

2. **Start the frontend application**
   - Make sure you have an `.env` file to add the backend url before starting the frontend application.
     After creating the `.env` add the following in the file so it can be read before starting it.
     ```sh
     VITE_API_BASE_URL=http://localhost:8000
     ```
   - Run the following command to start the frontend application:
     ```sh
     npm run dev
     ```

3. **Shutdown Database Container**:
   - To stop the database container, use:
     ```sh
     make docker-down
     ```
   - This command will attempt to stop the database container using Docker Compose. If Docker Compose V2 is not available, it will fall back to Docker Compose V1.

4. **Lint the Code**:
   - To lint the code, use:
     ```sh
     make lint
     ```
   - This command will run the linting process to check for code quality issues.

By using these commands, you can easily manage the build, and run processes for the application.

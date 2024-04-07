# restaurant-reserve

---

# Project Name

This project includes building a Go application with a PostgreSQL database using Docker Compose.

## Getting Started

### Prerequisites

Make sure you have Docker and Docker Compose installed on your system.

### How to Use

1. Clone this repository to your local machine.

2. Update the environment variables in the `docker-compose.yml` file:

   - `POSTGRES_USER`: Username for PostgreSQL database.
   - `POSTGRES_PASSWORD`: Password for PostgreSQL database.
   - `POSTGRES_DB`: Name of the PostgreSQL database.

3. Build and run the project using the following command:

   ```bash
   docker-compose up --build
   ```

4. Access the API documentation at [http://localhost:8080/swagger/index.html#/](http://localhost:8080/swagger/index.html#/)

5. You can now interact with the API endpoints as documented in the Swagger UI.

## Environment Variables

To set up the PostgreSQL database, you can customize the following environment variables in the `docker-compose.yml` file:

- `POSTGRES_USER`: Replace `user` with your desired username.
- `POSTGRES_PASSWORD`: Replace `password` with your desired password.
- `POSTGRES_DB`: Replace `database` with your desired database name.


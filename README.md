# Introduction

This REST API is built using Go and PostgreSQL. It facilitates administrative tasks for teachers regarding their students.

The following endpoints are available:

- `/api/register`: Registers one or more students to a specified teacher.
- `/api/commonstudents`: Retrieves a list of common students associated with specified teachers.
- `/api/suspend`: Suspends a specific student.
- `/api/retrievefornotifications`: Retrieves a list of students associated with a specified teacher who have been mentioned.

# Setup

Before running this application, ensure that you have the following software installed:

- [Make](https://www.gnu.org/software/make/)
- [Docker](https://www.docker.com/)

> *This application has been tested with Go version 1.22.0.

# Running the Server

The server is accessible at `http://localhost:3000`.

To run the server locally, execute the following command:
```bash
make run
```

Additionally, a [Postman collection](OneCV-Assgn.postman_collection.json) is provided to assist with API testing.

# Cleanup local run

To cleanup generated artefacts and Docker containers, execute the following command:
```bash
make prod-clean
```

# Testing

Testing server is running at `http://localhost:3001`.

To test the API endpoints, execute the following command:
```bash
make test
```

>*Note: Docker Desktop must be running during testing

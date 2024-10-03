Project Wraith

This is a Go project built using the Fiber framework, a fast and lightweight web framework for Go. This project demonstrates basic usage of Fiber for building RESTful APIs.

Features

- User authentication and session management
- CRUD operations for user management
- Password reset functionality
- JSON and HTML responses
- Integration with external services (e.g., email and SMS)

Installation

1. Clone the repository:

   git clone https://github.com/yourusername/your-repo.git
   cd your-repo

2. Install dependencies:

   Ensure you have Go installed. Then, run:

   go mod tidy

3. Run the application:

   go run main.go

   The application will start on localhost:8080 by default.


Run Swagger

1. Run the Swagger CLI:

   go install github.com/swaggo/swag/cmd/swag@latest
   swag init
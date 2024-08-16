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

   The application will start on localhost:3000 by default.

API Endpoints

User Login

- URL: /user/login
- Method: POST
- Description: Authenticates a user and generates a session token.
- Request Body: JSON object with ID, Username, Email, Name, Phone, and Password.
- Responses:
    - 200 OK: Login successful, returns a session token.
    - 400 Bad Request: Invalid request body.
    - 401 Unauthorized: Invalid credentials.

User Registration

- URL: /user/register
- Method: POST
- Description: Registers a new user.
- Request Body: JSON object with ID, Username, Email, Name, Phone, and Password.
- Responses:
    - 200 OK: Registration successful.
    - 400 Bad Request: Invalid request body or registration error.

Get User Details

- URL: /user/{id}
- Method: GET
- Description: Retrieves user details based on the provided user ID.
- Path Parameters: id - User ID.
- Responses:
    - 200 OK: Returns user details.
    - 400 Bad Request: Invalid ID or request.
    - 404 Not Found: User not found.

Edit User Details

- URL: /user/edit
- Method: PUT
- Description: Updates user details.
- Request Body: JSON object with ID, Username, Email, Name, Phone, and Password.
- Responses:
    - 200 OK: User details updated successfully.
    - 400 Bad Request: Invalid request body or update error.

Remove User

- URL: /user/remove
- Method: DELETE
- Description: Removes a user based on the provided details.
- Request Body: JSON object with ID, Username, Email, Name, Phone, and Password.
- Responses:
    - 200 OK: User removed successfully.
    - 400 Bad Request: Invalid request body or removal error.

User Logout

- URL: /user/logout
- Method: POST
- Description: Logs out the user by expiring the session token.
- Responses:
    - 200 OK: Logout successful.
    - 401 Unauthorized: No session found.
    - 500 Internal Server Error: Failed to expire session.

Configuration

- JWT Secret: Set the JWT secret key in your environment variables or configuration file.
- Email and SMS Services: Ensure that the email and SMS services are correctly configured in your application.

License

This project is licensed under the MIT License. See the LICENSE file for details.

Contributing

Feel free to open issues and submit pull requests. Contributions are welcome!

Contact

For any questions or issues, please contact your-email@example.com.

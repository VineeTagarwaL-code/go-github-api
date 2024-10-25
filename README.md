
# Go GitHub API

A self-hosted Go-based GitHub API that allows users to fetch GitHub contributions data for a specified user. This API is designed for fast and reliable retrieval of GitHub contribution graphs, making it ideal for integrating with front-end applications.

## Features

- **Fetch GitHub Contribution Data**: Retrieve contributions data for a specified GitHub user.
- **Simple API**: Easy-to-use endpoints for seamless integration.
- **Self-Hosted**: Run the API on your own server for better control and security.
- **Built with Go**: Leveraging the speed and efficiency of Go for fast performance.

## Installation

1. **Clone the repository:**

   ```bash
   git clone https://github.com/VineeTagarwaL-code/go-github-api.git
   cd go-github-api
   ```

2. **Install dependencies:**

   Make sure you have Go installed. Run:

   ```bash
   go mod tidy
   ```

3. **Set environment variables:**

   Create a `.env` file in the root directory with the following variables:

   ```env
   GITHUB_ACCESS_TOKEN=your_github_access_token
   PORT=8080
   ```

   - `GITHUB_ACCESS_TOKEN`: Your GitHub personal access token (with required permissions).
   - `PORT`: The port you want the server to run on (default is 8080).

4. **Run the server:**

   ```bash
   go run main.go
   ```

   The server will start on `http://localhost:8080`.

## API Endpoints

### 1. **Get User Contributions**

   - **Endpoint**: `/api/{username}`
   - **Method**: `GET`
   - **Description**: Fetches the contributions data for a given GitHub username.
   - **Example**:

     ```bash
     curl http://localhost:8080/api/vineet
     ```

   - **Response**:

     ```json
     {
       "data": [
         {
           "date": "2024-10-24",
           "count": 5,
           "level": 3
         },
         ...
       ]
     }
     ```

## Project Structure

```
go-github-api/
│   main.go          # Entry point of the application
│   .env             # Environment variables
│   README.md        # Documentation
│
├───routes/          # API routes
├───controllers/     # Logic for handling API requests
├───models/          # Data models for GitHub contributions
└───utils/           # Utility functions
```

## Technologies Used

- **Go**: The primary programming language for building the API.
- **Go Modules**: For dependency management.
- **Net/HTTP**: Go's built-in package for handling HTTP requests.

## Contributing

1. Fork the repository.
2. Create a new branch (`git checkout -b feature-branch`).
3. Make your changes.
4. Commit your changes (`git commit -m 'Add some feature'`).
5. Push to the branch (`git push origin feature-branch`).
6. Open a pull request.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for more details.

## Contact

For any issues or suggestions, feel free to open an issue or reach out via [Vineet Agarwal](https://github.com/VineeTagarwaL-code).

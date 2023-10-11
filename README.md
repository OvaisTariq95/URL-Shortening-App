# URL-Shortening-App

This is a simple URL shortener application that allows you to create and manage shortened URLs.

## Prerequisites

Before running the application, make sure you have the following dependencies installed:

- Go (Golang) - [Download Go](https://golang.org/dl/)
- Git - [Download Git](https://git-scm.com/downloads)

## Getting Started
Clone this repository to your local machine:

   ```bash
   git clone https://github.com/yourusername/url-shortener.git
```
Change to the project directory:

  ```bash
cd url-shortener
```
 Run the application:
  ```bash
go run cmd/main.go
```
The application should now be running locally on http://127.0.0.1:8081. You can access it through a web browser or use API clients like Postman.

To change the host and port you need to open the main.go file & change 

```
Host   = "127.0.0.1"
Port   = "8081"
```

## Usage

  Shorten a URL:
To shorten a URL, send a POST request to /shorten with a JSON payload:

json
```
{
  "url": "http://www.example.com"
}
```
You will receive a response containing the shortened URL.

  Redirect to Original URL:

To redirect to the original URL associated with a shortcode, visit the url in the response generate in the POST call.

  Get Original URL:

To get the original URL associated with a shortcode, send a GET request to /original/{shortcode}. You will receive a JSON response with the original URL.

  Get All URLs

To get a list of all stored URLs, send a GET request to /all. You will receive a JSON response with a list of URLs.

  Contributing:

Feel free to contribute to this project by creating issues, suggesting improvements, or submitting pull requests.

  License:

This project is licensed under the MIT License - see the LICENSE file for details.

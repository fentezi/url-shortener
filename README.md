# URL Shortener

URL Shortener is a simple yet efficient web application that allows you to shorten long URLs into shorter and more manageable aliases. This project is built using Go and the Gin web framework, providing a robust and scalable solution for URL shortening.

## Features

- **Shorten URLs**: Easily shorten long URLs into short, easy-to-remember aliases.
- **Retrieve Original URLs**: Given a shortened alias, retrieve the original long URL.
- **Delete URL Aliases**: Remove previously created URL aliases when no longer needed.
- **Basic Authentication**: API endpoints are protected by basic authentication for added security.
- **SQLite Database Storage**: URL mappings are persistently stored in a SQLite database for efficient retrieval and management.

## Installation

To install and run the URL Shortener application locally, follow these steps:

1. **Clone the Repository**:

  ```bash
  git clone https://github.com/fentezi/url-shortener.git
  cd url-shortener
  go build
  ```

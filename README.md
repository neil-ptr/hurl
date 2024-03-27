# Project Name: hurl

## Overview

hurl is a command-line tool inspired by HTTPie, Postman and curl, designed for developers who prefer to manage and execute HTTP requests from the terminal. Unlike traditional tools that require manual input for each request, hurl allows users to store their HTTP requests in files and execute them directly from the command line. This approach enables easy version control, sharing among team members, and integration into scripts or CI/CD pipelines.

## Features

- **Save and Reuse Requests**: Store your HTTP requests in simple text files for easy reuse and version control.
- **Support for Multiple HTTP Methods**: GET, POST, PUT, DELETE, PATCH, and more.
- **Headers & Payload Handling**: Easily add headers and payloads to your requests.
- **Environment Variables**: Use environment variables in your request files for different deployment stages (e.g., development, staging, production).
- **Response Highlighting**: Color-coded response output for easy reading and debugging.
- **Request Chaining**: Chain requests together and use the response from one as input for another.
- **File Uploads**: Support for multipart file uploads.
- **History & Repeat**: Keep a history of your requests and repeat them with a single command.

## Installation
TODO

## Usage

1. **Creating a Request File**: Create a new file (e.g., `request.txt`) and write your HTTP request following the format:

    ```
    GET https://example.com/api/resource
    Authorization: Bearer your_token_here
    ```

2. **Executing a Request**: Use the `hurl` command followed by the path to your request file:

    ```bash
    hurl request.txt
    ```

3. **Using Variables**: To use environment variables in your requests, define them in your files like this:

    ```
    GET https://example.com/api/resource
    Authorization: Bearer {{API_TOKEN}}
    ```

    And execute your request like this:

    ```bash
    API_TOKEN=your_token_here hurl request.txt
    ```

4. **Viewing Response**: The response will be printed directly to your terminal, with syntax highlighting for JSON responses.

## Configuration

You can configure hurl globally by creating a `.hurlrc` file in your home directory. Available configurations include default headers, response timeout, and proxy settings.

## Contributing

TODO

## License

hurl is released under the MIT License. See the LICENSE file for more details.

---

For more information and detailed documentation, visit our [GitHub repository](https://github.com/yourusername/hurl).

Enjoy using hurl, and happy coding!


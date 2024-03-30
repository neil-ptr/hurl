# hurl 🤮

## Currently under construction 🚧

## Overview

Hurl is a command-line tool inspired by HTTPie, Postman, and curl, designed for developers who prefer to manage and execute HTTP requests from the terminal. Born out of a love for Neovim and the streamlined efficiency of terminal-based workflows, Hurl is crafted for those who thrive in the command-line environment. Unlike traditional tools that require manual input for each request, Hurl allows users to store their HTTP requests in files and execute them directly from the command line. Whether you're editing code, managing version control, or interacting with web services, Hurl is intended to keep you in the zone, without ever leaving the comforting embrace of the terminal.

## Features

- **Save and Reuse Requests**: Store your HTTP requests in simple text files for easy reuse and version control.
- **Support for Multiple HTTP Methods**: GET, POST, PUT, DELETE, PATCH, and more.
- **Headers & Payload Handling**: Easily add headers and payloads to your requests.
- **Environment Variables**: (TODO) Use environment variables in your request files for different deployment stages (e.g., development, staging, production).
- **Request Chaining**: (TODO) Chain requests together and use the response from one as input for another.
- **Response Highlighting**: Color-coded response output for easy reading and debugging.
- **File Uploads**: (TODO) Support for multipart file uploads. 
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



# hurl


## Motivation

You're a developer working in tmux and neovim developing backend services but everytime you want to invoke an endpoint by making an http call you need to move your right hand from the home row (oh god 🤮) and CMD + arrow key over to your Postman app to then use your mouse to send a request.

![cover](https://github.com/neil-and-void/hurl/assets/46465568/408c360b-36a8-4a9a-af4a-585a1854b8bd)


## Overview

Hurl is a command-line tool inspired by Postman and curl, designed for developers who prefer to manage and execute HTTP requests from the terminal. Born out of a love for Neovim and the streamlined efficiency of terminal-based workflows, Hurl is crafted for those who thrive in the command-line environment. Hurl allows users to store their HTTP requests in files and execute them directly from the command line. Whether you're editing code, managing version control, or interacting with web services, Hurl is intended to keep you in your flow state in the comforting embrace of your terminal.

## Features
https://github.com/neil-and-void/hurl/assets/46465568/adc2587e-844f-4459-bf90-cd5957e84ca5

- **Save and Reuse Requests**: Store your HTTP requests in simple text files for easy reuse and version control.
- **Support for Multiple HTTP Methods**: GET, POST, PUT, DELETE, and PATCH.
- **Headers & Payload Handling**: Easily add headers and payloads to your requests.
- **Environment Variables**: Use environment variables in your request files for different deployment stages (e.g., development, staging, production).
- **Response Highlighting**: Color-coded response output for easy reading and debugging.
- **File Uploads**: (TODO) Support for multipart file uploads. 
- **History & Repeat**: Keep a history of your requests and repeat them with a single command.

## Installation

```bash
brew update
brew tap neil-and-void/homebrew-hurl
brew install neil-and-void/homebrew-hurl/hurl
```

## Usage

The [examples](https://github.com/neil-and-void/hurl/tree/main/examples) folder has some requests to copy as a starting point to edit and tailor to your own needs.

1. **Creating a Request File**: Create a new file (e.g., `request.txt`) and write your HTTP request following the format:

   ```yaml
    POST https://example.com/api/resource
    Authorization: Bearer your_token_here
    ```
    
    Request bodies go below the headers separated by exactly 1 space and require a `Content-Type`. Otherwise the `Content-Type` is set as `text/plain`.
   
    ```yaml
    POST https://example.com/api/resource
    Authorization: Bearer your_token_here
    Content-Type: application/json

    {
        "test: 1
    }
    ```

3. **Executing a Request**: Use the `hurl` command followed by the path to your request file:

    ```bash
    $ hurl request.txt
    ```

4. **Using Variables**: To use environment variables in your requests, define them in your files like this:

    ```yaml
    POST {{BASE_URL}}/resource
    Authorization: Bearer {{API_TOKEN}}
    Content-Type: application/json

    {
        "test": {{SURE_WHY_NOT_HERE_TOO}}
    }
    ```

    And execute your request like this:

    ```bash
    $ API_TOKEN=your_token_here hurl request.txt
    ```

5. **Viewing Response**: The response will be printed directly to your terminal, with syntax highlighting for JSON responses.

```json
{
  "userId": 1,
  "id": 1,
  "title": "delectus aut autem",
  "completed": false
}
```

For large responses or non-human readable formats (.pdf, .png, .word...), you can also output response bodies to files using `-o` flag
```bash
$ hurl -o=./response.json examples/post.txt
```



## Docs

### Sending a Basic Request

Request files are meant to look as close to raw HTTP requests as possible without becoming too tedious to manage. A basic request has a request line and some headers. 
```yaml
GET http://wealthsimple.com          # [method] [url]
Authorization: Bearer jwt            # [header]: [value]
User-Agent: idk                      # [header]: [value]
```

### Requests with Bodies
If it has a body, it is separated with exactly 1 newline below the headers, similar to a raw http request. It is also recommended to have a `Content-Type` header, if one is not present then the header is set to `text/plain`.

```yaml
POST http://wealthsimple.com         # [method] [url]
Content-Type: application/json       # [header]: [value]
Authorization: Bearer jwt            
                                     # newline if there is a body
{                                    # body...
    "json": 123
}
```
### Single File Uploads
```yaml
POST http://wealthsimple.com         # [method] [url]
Content-Type: image/png              # [header]: [value]
Authorization: Bearer jwt            
                                     # newline if there is a body
@file=/path/to/file.png
```

### Multi Part (Use for Multiple File Upload) `multipart/form-data`
```yaml
POST http://wealthsimple.com                       # [method] [url]
Content-Type: multipart/form-data                  # [header]: [value]
Authorization: Bearer jwt            
                                                   # newline if there is a body
form-data; name="jj"; value="abrams"
form-data; name="bruh"; filename="path/to/file/image.png"
form-data; name="bruh2"; filename="path/to/file/image2.png"
```

### Form `application/x-www-form-urlencoded`
```yaml
POST http://wealthsimple.com                       # [method] [url]
Content-Type: application/x-www-form-urlencoded    # [header]: [value]
Authorization: Bearer jwt            
                                                   # newline if there is a body
name=John+Doe&age=30&city=New+York
```

### Environment Variables
```yaml
POST {{BASE_URL}}         # [method] [url]
Content-Type: application/json       # [header]: [value]
Authorization: Bearer jwt            
                                     # newline if there is a body
{                                    # body...
    "json": 123
}
```
To run you can either
```
$ BASE_URL=https://wealthsimple.com
```

or give a path to your `.env` file in your `hurl.json` file

```yaml
// hurl.json
{
    // path to your .env file
    "env": "/path/to/.env/file"
}
```

```yaml
# .env
BASE_URL=https://wealthsimple.com
```

### Flags
all flags need to come before the path to the request file.
* `-o=/path/to/file.json`: path to a file to output response body content
* `-v`: verbose out, prints all request and response headers in a format similar to a raw HTTP request and response

## Configuration
You can configure hurl by creating a `hurl.json` file in your current working directory. Available configurations include setting `.env` file path, default headers (TODO), response timeout (TODO). Below is an example config.
```yaml
{
    // path to your .env file
    "env": "/path/to/.env/file"
}
```

## License

hurl is released under the MIT License. See the LICENSE file for more details.

---

For more information and detailed documentation, visit our [GitHub repository](https://github.com/yourusername/hurl).

Enjoy using hurl, and happy coding!


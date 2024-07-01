# Bon Voyage CLI

`bon-voyage-cli` is a command-line tool for interacting with the `Bon Voyage` service. 
The tool also stores the authentication token securely in a temporary file.

## Features

- **User**: Authenticate, register or modify a user account.
- **Device**: List, configure or get a tunnel to connected devices
- **Session**: Change the username of the authenticated user.
- **Snippet**: Create code orlog snipt and share them with a read-only link.

## Installation

### Prerequisites

- Go programming language installed (version 1.22 or higher).

### Build

1. Clone the repository:
    ```sh
    git clone https://github.com/bitcrushtesting/bon-voyage-cli.git
    cd bon-voyage-cli
    ```

2. Build the binary:
    ```sh
    ./build.sh
    ```

## Usage

### Login

Authenticate with the server and store the token.

```sh
./bon-voyage-cli login <username>
```

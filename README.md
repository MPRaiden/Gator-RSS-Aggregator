# Gator RSS Aggregator

Gator is a command-line tool for managing RSS feeds and users, written in Go. It allows users to register, login, add and follow RSS feeds, and browse content. The project interfaces with a PostgreSQL database for persistent storage of user data, feeds, and posts.

## Table of Contents

- [Features](#features)
- [Installation](#installation)
- [Configuration](#configuration)
- [Usage](#usage)
- [Commands](#commands)
- [Contributing](#contributing)
- [License](#license)

## Features

- **User Management**: Register, login, view, and reset users.
- **RSS Feed Management**: Add, follow, unfollow, and browse RSS feeds.
- **Database Integration**: Utilizes PostgreSQL to store and retrieve feed and user data.
- **Command-line Interface**: Provides an easy-to-use CLI for interacting with the application.

## Installation

### Prerequisites

- Go 1.16 or later
- PostgreSQL database
- [lib/pq](https://github.com/lib/pq): PostgreSQL driver for Go
- Other dependencies can be installed via Go modules.

### Steps

1. Clone the repository:
    ```bash
    git clone https://github.com/your-username/gator
    cd gator
    ```

2. Install dependencies:
    ```bash
    go mod tidy
    ```

3. Build the project:
    ```bash
    go build
    ```

## Configuration

Gator requires a configuration file (`config.json`) located in the `internal/config` directory with the following structure:

```json
{
    "DBURL": "postgres://username:password@localhost:5432/database_name?sslmode=disable",
    "CurrentUserName": ""
}

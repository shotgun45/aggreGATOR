# aggreGATOR

aggreGATOR is a multi-user blog/RSS aggregator CLI written in Go. It uses PostgreSQL for persistent storage and supports user registration, feed management, and post aggregation.

## Prerequisites
- **Go** (version 1.20+ recommended)
- **PostgreSQL** (running locally or accessible via network)

## Installation
Install the CLI using Go:

```
go install github.com/shotgun45/aggreGATOR@latest
```

Or, if working locally:

```
go build -o gator .
```

## Configuration
Create a config file named `.gatorconfig.json` in your home directory. Example:

```
{
  "db_url": "postgres://postgres:postgres@localhost:5432/gator",
  "current_user_name": ""
}
```

- Replace the `db_url` with your actual PostgreSQL connection string.
- The `current_user_name` will be set when you log in or register.

## Usage
Run the CLI from your project directory:

```
go run . <command> [args]
```

Or, if you built the binary:

```
./gator <command> [args]
```

### Common Commands
- `register <username>`: Create a new user.
- `login <username>`: Log in as an existing user.
- `addfeed <name> <url>`: Add a new RSS feed and follow it.
- `feeds`: List all feeds.
- `follow <feed_url>`: Follow an existing feed.
- `browse [limit]`: Show recent posts for the current user (default limit is 2).
- `agg <duration>`: Start periodic aggregation (e.g., `agg 1m`).

## Tracking Changes
All changes should be tracked with Git. Commit regularly:

```
git add .
git commit -m "your message"
```

## Notes
- Make sure your PostgreSQL server is running and accessible.
- The CLI will create and migrate the database tables automatically if configured.
- For more commands and details, see the source code or run with no arguments for help.

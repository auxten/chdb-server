# chDB Server

This project is a Go-based web server that integrates with the [chDB-go](https://github.com/chdb-io/chdb-go) and provides a web interface for executing database queries. It embeds the static file `play.html` directly within the binary for easy distribution and deployment.

## Features

- **Query Execution**: Allows users to execute queries against a CHDB database.
- **Session Management**: Persists the database connection across multiple queries.
- **Static Content**: Serves a static HTML file (`play.html`) embedded within the executable.

## Prerequisites

Before you start, ensure you have the following installed:
- Go 1.16 or higher
- libchdb

```bash
curl -sL https://lib.chdb.io | bash
```

## Installation

1. Go Install

```bash
go install github.com/auxten/chdb-server@latest
$GOPATH/bin/chdb-server
```

2. Build from source

```bash
git clone https://github.com/auxten/chdb-server.git
make build
./chdb-server
```
### Session Path
You can configure the server by setting the following environment variables before running the server:

- DATA_PATH: Specifies the directory path for CHDB session data. Defaults to .chdb_data if not set.

## Usage


Open your browser and navigate to `http://localhost:8123` to start querying any data.

#### Create a new database

```sql
CREATE TABLE users (uid Int16, name String, age Int16) ENGINE=Memory;

INSERT INTO users VALUES (1231, 'John', 33);
INSERT INTO users VALUES (6666, 'Ksenia', 48);
INSERT INTO users VALUES (8888, 'Alice', 50);

SELECT * FROM users;
```

#### Query remote parquet data

```sql
SELECT RegionID, SUM(AdvEngineID), COUNT(*) AS c, AVG(ResolutionWidth), COUNT(DISTINCT UserID)
                        FROM url('https://datasets.clickhouse.com/hits_compatible/athena_partitioned/hits_0.parquet') GROUP BY RegionID ORDER BY c DESC LIMIT 10
```


# Yeo!
Yeo! is a tiny backup utility for PostgreSQL databases. The intention with Yeo! is not to substitute utilities such as `pg_dump` or `pg_restore` but to speed up trivial Postgres operations related to backing up and restoring databases.

<br>

<img src="examples/clone-demo.gif" width="350"/>

## Features

 With `yeo` you can:
* Back up a database using the `backup` command - known in Postgres terminology as **dump**. For more [here](https://www.postgresql.org/docs/current/app-pgdump.html).
* Revive a database using the `revive` command - known in Postgres terminology as **restore**. For more [here](https://www.postgresql.org/docs/current/app-pgrestore.html).
* Clone a database using the `clone` command - a combination of dump and restore in the Postgres terminology.

## Installation
> **Yeo! is currently only available for MacOS and Linux users via [Homebrew](https://brew.sh/)**


To install `yeo`:

```bash
# Installing libpq is essential
brew install libpq # Make sure libpq is in your PATH
brew tap jjcfrancisco/yeo # Adds the Github repository as a tap
brew install yeo
```

Alternatively, you can use the already-built binaries for MacOS [here](https://github.com/jjcfrancisco/yeo/releases/).

## Requirements
The only thing you need after the installation process is a `databases.json` file with credentials of databases. It must be populated with credentials of at least one database.

**IMPORTANT**: Yeo! will look for the `databases.json` in your user's home directory e.g. `/Users/joebloggs/databases.json`.

The `databases.json` must follow this JSON structure:
```json
{
    "databases": [
        {
            "name": "local",
            "isLocal": true,
            "database": "my_local_db",
            "port": "5432",
            "host": "localhost",
            "user": "yeo",
            "password": "mysecretpassword"
        },
        {
            "name": "development",
            "isLocal": false,
            "database": "dev",
            "port": "5432",
            "host": "yeo.host.com",
            "user": "yeo",
            "password": "mysecretpassword"
        }
    ]
}
```

## Usage
To back up a database:
```bash
# 'local' is the name set in the databases.json and can be personalised

yeo backup local db_backup.dump
```

To revive a database from a backup file:
```bash
# 'local' or 'prod' are names set in the database.json and can be personalised. The '--allow' flag allows to revive into non-local databases. 

yeo revive db_backup.dump local
yeo revive --allow db_backup.dump prod
```

To clone a database (backup + revive):
```bash
# 'local1', 'local2' or 'prod' are names set in the database.json and can be personalised. The '--allow' flag allows to revive into non-local databases.

yeo clone local1 local2
yeo clone --allow local1 prod
```

## Future improvements
* Validation for `databases.json` against schema.
* It should not be 'isLocal' but 'lock' for consistency reasons.

## License

See [`LICENSE`](./LICENSE)

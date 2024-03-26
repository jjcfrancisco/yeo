# Yeo!
Yeo! is a tiny backup utility for PostgreSQL databases.

## Features

 With `yeo` you can:
* Back up a database using the `backup` command - known in Postgres terminology as **dump**. For more [here](https://www.postgresql.org/docs/current/app-pgdump.html).
* Revive a database using the `revive` command - known in Postgres terminology as **restore**. For more [here](https://www.postgresql.org/docs/current/app-pgrestore.html).
* Clone a database using the `clone` command - a combination of dump and restore in the Postgres terminology.

## Installation
To install `yeo`:
1. It is **essential** to have Postgres utilities such as `pg_dump`, `pg_restore`, `createdb` or `dropdb` in your system. This are included in *libpq* for MacOS or the entire PostgreSQL suite for Windows:
```bash
# Homebrew (MacOS & Linux)
brew install libpq # Make sure libpq is in your PATH

# Chocolatey (Windows)
choco install postgresql # Easiest is to install postgres
```
2. Then install `yeo`:
```bash
go get github.com/jjcfrancisco/yeo
```

Alternatively, you can use the already-built binaries for MacOS & Windows [here](https://github.com/jjcfrancisco/yeo/releases/). Just bear in mind that when running the binary, you must be in the same directory as where the binary lives (and create a `databases.json`) and call it like this -> `./yeo [command]...`

For instance:
```bash
./yeo clone development local
``` 

## Requirements
The only thing you need after the installation process is a `databases.json` file with credentials of databases. It must be populated with credentials of at least one database.

The `databases.json` must follow this JSON structure:
```json
{
    "databases": [
        {
            "name": "local",
            "database": "my_local_db",
            "port": "5432",
            "host": "localhost",
            "user": "yeo",
            "password": "mysecretpassword"
        },
        {
            "name": "development",
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

yeo revive local db_backup.dump
yeo revive --allow prod db_backup.dump
```

To clone a database (backup + revive):
```bash
# 'local1', 'local2' or 'prod' are names set in the database.json and can be personalised. The '--allow' flag allows to revive into non-local databases.

yeo clone local1 local2
yeo clone --allow local1 prod
```

## Future improvements
* Validation for `databases.json` against schema.

## License

See [`LICENSE`](./LICENSE)

<p align="center">
  <img alt="GLoader logo" src="assets/logo.png" height="150" />
  <h3 align="center">GLoader</h3>
  <p align="center">Migrate your data across any source and destination with a single command!</p>
</p>

[![Go Reference](https://pkg.go.dev/badge/github.com/mohammadv184/gloader.svg)](https://pkg.go.dev/github.com/mohammadv184/gloader)
[![GitHub license](https://img.shields.io/github/license/mohammadv184/gloader)](https://github.com/mohammadv184/gloader/blob/main/LICENSE)
[![GitHub issues](https://img.shields.io/github/issues/mohammadv184/gloader)](https://github.com/mohammadv184/gloader/issues)
[![GitHub release](https://img.shields.io/github/release/mohammadv184/gloader.svg)](https://github.com/mohammadv184/gloader/releases)
[![Go Report Card](https://goreportcard.com/badge/github.com/mohammadv184/gloader)](https://goreportcard.com/report/github.com/mohammadv184/gloader)
[![gloader](https://snapcraft.io/gloader/badge.svg)](https://snapcraft.io/gloader)

---

GLoader is a powerful and flexible CLI tool for data migration between different databases. 
It provides a seamless way to migrate your data from any source database to any destination database,
Whether you are upgrading your database or moving data between different systems;
GLoader makes the process efficient and reliable.


# List of contents

- [List of supported databases](#list-of-supported-databases)
- [Getting started](#getting-started)
  - [Installation](#installation)
  - [Usage](#usage)
- [Glossary](#glossary)
- [Contributing](#contributing)
- [Security](#security)
- [Credits](#credits)
- [License](#license)

# List of supported databases
| Database \ As |       Source       |    Destination     |
|---------------|:------------------:|:------------------:|
| MySQL         | :white_check_mark: |     :x: (Soon)     |
| CockroachDB   |     :x: (Soon)     | :white_check_mark: |
| PostgreSQL    |     :x: (Soon)     |     :x: (Soon)     |
| MongoDB       |        :x:         |        :x:         |
| SQLite        |        :x:         |        :x:         |
| SQL Server    |        :x:         |        :x:         |
| Oracle        |        :x:         |        :x:         |
| Redis         |        :x:         |        :x:         |
| Cassandra     |        :x:         |        :x:         |
| Elasticsearch |        :x:         |        :x:         |
| Kafka         |        :x:         |        :x:         |
| RabbitMQ      |        :x:         |        :x:         |
| DynamoDB      |        :x:         |        :x:         |


**Note**: The database that is marked with :x: will be supported soon. 
However, if you have time and want to contribute, you can help us to support them faster by contributing to the project.

# Getting started

---

## Installation
### Debian based Linux distributions
```bash
curl -fsSL https://repo.gloader.tech/apt/gpg.key | sudo gpg --dearmor -o /etc/apt/trusted.gpg.d/gloader.gpg
echo "deb https://repo.gloader.tech/apt * *" > /etc/apt/sources.list.d/gloader.list
sudo apt update && sudo apt install gloader
```
### Fedora / RHEL / CentOS
```bash
echo '[gloader]
name=Gloader
baseurl=https://repo.gloader.tech/yum
enabled=1
gpgcheck=1
gpgkey=https://repo.gloader.tech/yum/gpg.key' | sudo tee /etc/yum.repos.d/gloader.repo
sudo yum install gloader
```
### Ubuntu
[![Get it from the Snap Store](https://snapcraft.io/static/images/badges/en/snap-store-black.svg)](https://snapcraft.io/gloader)

--- OR ---

```bash
sudo snap install gloader
```

### Homebrew
```bash
brew tap mohammadv184/gloader
brew install gloader
```

### Go install (Not recommended)
```bash
go install github.com/mohammadv184/gloader@latest
```

### Binary builds
You can download the binary builds from the [releases](https://github.com/mohammadv184/gloader/releases)

### Deb and RPM packages
You can download the deb and rpm packages from the [releases](https://github.com/mohammadv184/gloader/releases)


## Usage
```bash
gloader run <source-dsn> <destination-dsn> [flags]
flags:
      --end-offset stringToInt64           end offset for each table (default [])
  -e, --exclude strings                    exclude tables from migration
  -f, --filter stringToStringSlice         filter data to migrate
      --filter-all strings                 filter data to migrate (all tables)
  -h, --help                               help for run
  -r, --rows-per-batch uint                number of rows per batch (default 100)
  -s, --sort stringToStringSlice           sort data to migrate in ascending order
      --sort-all strings                   sort data to migrate in ascending order (all tables)
  -S, --sort-reverse stringToStringSlice   sort data to migrate in descending order
      --sort-reverse-all strings           sort data to migrate in descending order (all tables)
      --start-offset stringToInt64         start offset for each table (default [])
  -t, --table strings                      migrate only these tables
  -w, --workers uint                       number of workers (default 3)

```
#### Arguments
- **source-dsn**: The Data Source Name (DSN) used to connect to the source database.
- **destination-dsn**: The DSN used to connect to the destination database.
#### Flags
- **--start-offset**: The initial row offset for each table. This sets the starting point for migrating rows from the source to the destination.
- **--end-offset**: The final row offset for each table, limiting the number of rows migrated.
- **--exclude**: Exclude specific tables from the migration process.
- **--table**: Selectively migrate specific tables.
- **--filter**: Apply data filters to rows being migrated. Use operators such as `=`, `!=`, `>`, `>=`, `<`, and `<=`.
- **--filter-all**: Apply a universal data filter for all tables.
- **--sort**: Sort data in ascending order before migration. 
- **--sort-all**: Apply ascending sorting for all tables.
- **--sort-reverse**: Sort data in descending order before migration.
- **--sort-reverse-all**: Apply descending sorting for all tables.
- **--rows-per-batch**: Set the number of rows migrated per batch.
- **--workers**: Specify the number of parallel migration workers.
### Examples
#### Migrate all tables from the source to the destination
```bash
gloader run "mysql://root:root@tcp(localhost:3306)/source" "mysql://root:root@tcp(localhost:3306)/destination"
```
#### Migrate only the users table from the source to the destination
```bash
gloader run "mysql://root:root@tcp(localhost:3306)/source" "mysql://root:root@tcp(localhost:3306)/destination" --table users
```
#### Migrate all tables except the users table from the source to the destination
```bash
gloader run "mysql://root:root@tcp(localhost:3306)/source" "mysql://root:root@tcp(localhost:3306)/destination" --exclude users
```
#### Migrate all tables from the source to the destination with a filter
```bash
gloader run "mysql://root:root@tcp(localhost:3306)/source" "mysql://root:root@tcp(localhost:3306)/destination" --filter "id > 100"
```
#### Migrate all tables from the source to the destination with a filter for all tables
```bash
gloader run "mysql://root:root@tcp(localhost:3306)/source" "mysql://root:root@tcp(localhost:3306)/destination" --filter-all "id > 100"
```
#### Migrate all tables from the source to the destination with a sort
```bash
gloader run "mysql://root:root@tcp(localhost:3306)/source" "mysql://root:root@tcp(localhost:3306)/destination" --sort "id"
```
#### Migrate all tables from the source to the destination with a sort for all tables
```bash
gloader run "mysql://root:root@tcp(localhost:3306)/source" "mysql://root:root@tcp(localhost:3306)/destination" --sort-all "id"
```
#### Migrate all tables from the source to the destination with a sort-reverse
```bash
gloader run "mysql://root:root@tcp(localhost:3306)/source" "mysql://root:root@tcp(localhost:3306)/destination" --sort-reverse "id"
```
#### Migrate all tables from the source to the destination with a sort-reverse for all tables
```bash
gloader run "mysql://root:root@tcp(localhost:3306)/source" "mysql://root:root@tcp(localhost:3306)/destination" --sort-reverse-all "id"
```
#### Migrate all tables from the source to the destination with a start offset for each table
```bash
gloader run "mysql://root:root@tcp(localhost:3306)/source" "mysql://root:root@tcp(localhost:3306)/destination" --start-offset "users=100"
```
#### Migrate all tables from the source to the destination with a end offset for each table
```bash
gloader run "mysql://root:root@tcp(localhost:3306)/source" "mysql://root:root@tcp(localhost:3306)/destination" --end-offset "users=100"
```
#### Migrate all tables from the source to the destination with a rows per batch and workers
```bash
gloader run "mysql://root:root@tcp(localhost:3306)/source" "mysql://root:root@tcp(localhost:3306)/destination" --sort-all "id" --rows-per-batch 1000 --workers 10
```
# Glossary
- **Data**: A fundamental unit of information comprising a key, value, and data type. Data is the smallest entity in GLoader and represents the content being migrated.
- **DataType**: Denotes the type of data, such as strings, numbers, booleans, dates, JSON, etc., used in migration.
- **DataSet**: : A collection of data items representing a single row in a relational database.
- **DataBatch**: A group of DataSets migrated together in a single operation.
- **DataBuffer**: A storage space containing DataBatches fetched from the source database.
- **DataCollection**: An assembly of DataSets representing a table in a relational database.
- **DataMap**: A map associating DataCollections with their respective attributes, is equivalent to a table schema in a relational database.
- **Database**: A collection of DataCollections representing a database in a relational context.
- **Migration**: The process of relocating data from a source database to a target destination.

## Security

If you discover any security-related issues, please email mohammadv184@gmail.com instead of using the issue tracker.

## Credits

- [Mohammad Abbasi](https://github.com/mohammadv184)
- [All Contributors](../../contributors)

## License

The MIT License (MIT). Please see [License File](LICENSE) for more information.
module github.com/mohammadv184/gloader

go 1.19

require (
	github.com/go-sql-driver/mysql v1.7.0
	github.com/spf13/cobra v1.6.1
)

// it is a fork of github.com/lib/pq to fix the issue of the twice escaping
// ref: https://github.com/lib/pq/issues/1118
require github.com/mohammadv184/pq v0.1.2

require (
	github.com/inconshreveable/mousetrap v1.1.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
)

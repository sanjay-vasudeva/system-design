module conn_pool

go 1.24.1

require (
	filippo.io/edwards25519 v1.1.0 // indirect
    github.com/go-sql-driver/mysql v1.9.2
	github.com/sanjay-vasudeva/ioutil v1.0.0
    github.com/sanjay-vasudeva/queue v1.0.0
)

replace github.com/sanjay-vasudeva/ioutil v1.0.0 => ../../ioutil

replace github.com/sanjay-vasudeva/queue v1.0.0 => ../queue

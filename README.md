# mini
connect to a bitclout peer and stream recent transactions (before they are in a block)

use case:

- you want to read from the firehose of live real-time incoming tx's
- you don't want any badgerdb writes
- you don't want to validate any tx
- you simply want a stream of brand new, very fresh messages so you can inspect them
- some you might add to a chanel for further processing

# building

go mod tidy
go build

# running

./mini


# Conduit Connector Snowflake

### General

The Snowflake connector is one of [Conduit](https://github.com/ConduitIO/conduit) plugins. It provides the source
snowflake connector.

### Prerequisites

- [Go](https://go.dev/) 1.18
- (optional) [golangci-lint](https://github.com/golangci/golangci-lint) 1.45.2
- (optional) [mock](https://github.com/golang/mock) 1.6.0

### Configuration

The config passed to `Configure` can contain the following fields.

haris: which spec does this connection string follow? is it something by Snowflake itself (maybe their JDBC conn string?)
haris: also, is it possible to support other types of authentication mentioned here: https://docs.snowflake.com/en/user-guide/authentication.html
haris: it might make sense to make the batch size configurable too?

| name         | description                                                                                                                                                                                                                                     | required | example                                                |
|--------------|-------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|----------|--------------------------------------------------------|
| `connection` | Snowflake connection string.<br/>Supported formats:<br><code>user:password@my_organization-my_account/dbname/schemaname</code> or <br><code>username[:password]@hostname:port/dbname/schemaname </code><br><b>Important</b>: Schema is required | yes      | "user:password@my_organization-my_account/mydb/schema" |
| `table`      | The table name in snowflake db.                                                                                                                                                                                                                 | yes      | "users"                                                |
| `columns`    | Comma separated list of column names that should be included in each Record's payload. By default: all columns.                                                                                                                                 | no       | "id,name,age"                                          |
| `key`        | Column name that records should use for their `Key` fields.                                                                                                                                                                                     | yes      | "id"                                                   |



### How to build it

Run `make build`.

### Testing

Run `make test`.

### Snowflake Source

The Snowflake Source Connector connects to a snowflake with the provided configurations, using
`connection`, `table`,`columns`, `key`  and using snowflake driver. 
Source method `Configure`  parse the configurations.
 `Open` method is called to start the connection from the provided position get the
data from snowflake db. The `Read` return next record. The `Ack` method 
checks if the record with the position was recorded. The `Teardown` do gracefully shutdown.

#### Snapshot Iterator

haris: how are interrupted snapshots handled? 

The snapshot iterator starts getting data from the table using `select` query with `limit` and `offset`. Batch size is 1000,
offset value is zero for first time. Iterator save information from table to `data` slice variable.
Iterator `HasNext` method check if next element exist in `data` using variable `index` and if it is needed
change offset and run select query to get new data with new offset. Method `Next` gets next element and converts 
it to `Record` sets metadata variable `table`, set metadata variable `action` - `insert`, increase `index`.

#### CDC Iterator

haris: it's worth considering the option to skip snapshot and start with CDC. that's something some connectors offer.
but imho out of scope of this PR.

haris: can you briefly describe the metadata columns in the generated stream? also, can you please describe how is the 
stream/table populated? we don't have to go into implementation details, just high level.

CDC iterator starts working if snapshot iterator method `HasNext` return false.
CDC iterator uses snowflake stream (more information about streams https://docs.snowflake.com/en/user-guide/streams-intro.html) 
When source starts work first time iterator <b>creates</b> stream with name `conduit_stream_{table}` to `table` from
config, <b>creates</b> table for consuming stream with name `conduit_tracking_{table}`. 
This consuming table has the same schema as `table`  with additional metadata columns:
`METADATA$ACTION`, `METADATA$ISUPDATE`, `METADATA$ROW_ID`, `METADATA$TS`. Then iterator consume
data from stream using insert query to consuming table. Iterator run select query for get data
from consuming table using limit and offset and ordering by `METADATA$TS`. Batch size is 1000, offset value is zero for first time.
Iterator save information from table to 
`data` slice variable. Iterator `HasNext` method check if next element exist in `data` using variable
`index` and if it is needed change offset and run select query to get new data with new offset.
Method `Next` gets next element converts it to `Record` checks action(can be `insert`, `delete`, `update`)
using metadata columns `METADATA$ACTION`, `METADATA$ISUPDATE` and increase `index`.

#### Position

[comment]: <> (give an example of position)
Position has fields: `type` (`c` - CDC or `s`- Snapshot), `element`(index of element of current
offset), `offset`.

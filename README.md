# Burlesque

Burlesque is a [message processing queue](http://en.wikipedia.org/wiki/Message_queue) writen in Go. It exposes queues using the HTTP API that allows publishing messages and subscribing to them. See the [API](#api) section for more details.

Burlesque stores messages inside a [Kyoto Cabinet](http://fallabs.com/kyotocabinet/) database.

## Installation

OSX:
```
brew install go
brew install kyoto-cabinet
go get github.com/KosyanMedia/burlesque
```

Linux:
```
...
```

## Starting

Use the following arguments to the `burlesque` executable:

| Argument | Description | Defaults |
| -------- | ----------- | -------- |
| `-storage` | Kyoto Cabinet storage path (e.g. `storage.kch#msiz=524288000`) | `-` |
| `-environment` | Process environment: `development` or `production` | `development` |
| `-port` | Server HTTP port | `4401` |
| `-rollbar` | [Rollbar](https://rollbar.com/) token | ||

## Storage argument

`-storage` argument defines a way the data will be stored into a database. You can read more on Kyoto Cabinet database types [here](http://fallabs.com/kyotocabinet/spex.html#tutorial_dbchart).

If you need a **temporary in-memory** storage use the following symbols as the `-storage` value:

| Value | Database Type |
| ----- | ------------- |
| `-` | `ProtoHashDB` Prototype hash database. On-memory database implemented with `std::unorderd_map` |
| `+` | `ProtoTreeDB` Prototype tree database. On-memory database implemented with `std::map` |
| `:` | `StashDB` Stash database. On-memory database saving memory |
| `*` | `CacheDB` Cache hash database. On-memory database featuring [LRU](http://en.wikipedia.org/wiki/Cache_algorithms#Examples) deletion |
| `%` | `GrassDB` Cache tree database. On-memory database of B+ tree: cache with order |

In order to use a **persistent database** use the path to the database file (or directory) as the `-storage` argument value. File extension in the database path defines the type of the database created.

**Example:** `/path/to/my/storage.kch`

| File Extension | Database Type |
| -------------- | ------------- |
| `kch` | `HashDB` File hash database. File database of hash table: typical DBM |
| `kct` | `TreeDB` File tree database. File database of B+ tree: DBM with order |
| `kcd` | `DirDB` Directory hash database. Respective files in a directory of the file system |
| `kcf` | `ForestDB` Directory tree database. Directory database of B+ tree: huge DBM with order |
| `kcx` | `TextDB` Plain text database. Emulation to handle a plain text file as a database |

In addition to defining database type you can also add [tuning parameters](http://fallabs.com/kyotocabinet/spex.html#tips) to the `-storage` argument. Tuning parameters are separated by the `#` symbol, parameters' name and value are separated by the `=` symbol.

**Example:** `storage.kch#opts=c#zcomp=gz#msiz=524288000`

The table below describes tuning parameters.

| Parameter  | Description |
| ---------- | ----------- |
| `apow`     | Power of the alignment of record size |
| `bnum`     | Base hash table size (number of buckets of the hash table) |
| `capcnt`   | Capacity limit by the number of records (`capcnt=10000` means "keep in memory 10,000 records maximum) |
| `capsiz`   | Capacity limit by the size of records (`capsiz=536870912` means "keep in memory all the records that fit into 512 megabytes) |
| `dfunit`   | Unit step number of auto defragmentation |
| `fpow`     | Power of the capacity of the free block pool |
| `log`      | Path to the log file. Use `-` for the STDOUT, or `+` for the STDERR |
| `logkinds` | Kinds of logged messages. The value can be `debug`, `info`, `warn` or `error` |
| `logpx`    | Prefix of each log message |
| `msiz`     | Expected database memory usage |
| `opts`     | Additional options: `s`, `l` and `c` (can be specified together, e.g `lc`). `s` is for "small" and reduces the width of record addressing from 6 bytes to 4 bytes. As the result, the footprint for each record is reduced from 16 bytes to 12 bytes. However, it limits the maximum size of the database file up to 16GB. `l` is for "linear" and changes the data structure of the collision chain of hash table from binary tree to linear linked list. `c` enables compression of the record values. If the value is bigger than 1KB compression is effective. |
| `pccap`    | Capacity size of the page cache |
| `psiz`     | Page size |
| `rcomp`    | Comparator used to compare key names. `lex` for the lexical comparator, `dec` for the decimal comparator, `lexdesc` for the lexical descending comparator, or `decdesc` for the decimal descending comparator |
| `zcomp`    | Compression library: `zlib` for the [ZLIB raw](http://en.wikipedia.org/wiki/Zlib#Encapsulation) compressor, `def` for the ZLIB [deflate](http://en.wikipedia.org/wiki/DEFLATE) compressor, `gz` for the ZLIB [gzip](http://en.wikipedia.org/wiki/Gzip) compressor, `lzo` for the [LZO](http://en.wikipedia.org/wiki/Lempel%E2%80%93Ziv%E2%80%93Oberhumer) compressor, `lzma` for the [LZMA](http://en.wikipedia.org/wiki/Lempel%E2%80%93Ziv%E2%80%93Markov_chain_algorithm) compressor, or `arc` for the [Arcfour](http://en.wikipedia.org/wiki/RC4) cipher |
| `zkey`     | Cipher keyword used with compression |


The table below describes support of these parameters by the **in-memory** database types.

| Parameter  | `ProtoHashDB` | `ProtoTreeDB` | `StashDB` | `CacheDB` | `GrassDB` |
| ---------- | :-----------: | :-----------: | :-------: | :-------: | :-------: |
| `bnum`     | —             | —             | **Yes**   | **Yes**   | **Yes**   |
| `capcnt`   | —             | —             | —         | **Yes**   | —         |
| `capsiz`   | —             | —             | —         | **Yes**   | —         |
| `log`      | **Yes**       | **Yes**       | **Yes**   | **Yes**   | **Yes**   |
| `logkinds` | **Yes**       | **Yes**       | **Yes**   | **Yes**   | **Yes**   |
| `logpx`    | **Yes**       | **Yes**       | **Yes**   | **Yes**   | **Yes**   |
| `opts`     | —             | —             | —         | **Yes**   | **Yes**   |
| `pccap`    | —             | —             | —         | —         | **Yes**   |
| `psiz`     | —             | —             | —         | —         | **Yes**   |
| `rcomp`    | —             | —             | —         | —         | **Yes**   |
| `zcomp`    | —             | —             | —         | **Yes**   | **Yes**   |
| `zkey`     | —             | —             | —         | **Yes**   | **Yes**   |

The table below describes support of these parameters by the **persistent** database types.

| Parameter  | `HashDB` | `TreeDB` | `DirDB` | `ForestDB` | `TextDB` |
| ---------- | :------: | :------: | :-----: | :--------: | :------: |
| `apow`     | **Yes**  | **Yes**  | —       | —          | —        |
| `bnum`     | **Yes**  | **Yes**  | —       | —          | —        |
| `capcnt`   | —        | —        | —       | —          | —        |
| `capsiz`   | —        | —        | —       | —          | —        |
| `dfunit`   | **Yes**  | **Yes**  | —       | —          | —        |
| `fpow`     | **Yes**  | **Yes**  | —       | —          | —        |
| `log`      | **Yes**  | **Yes**  | **Yes** | **Yes**    | **Yes**  |
| `logkinds` | **Yes**  | **Yes**  | **Yes** | **Yes**    | **Yes**  |
| `logpx`    | **Yes**  | **Yes**  | **Yes** | **Yes**    | **Yes**  |
| `msiz`     | **Yes**  | **Yes**  | —       | —          | —        |
| `opts`     | **Yes**  | **Yes**  | **Yes** | **Yes**    | —        |
| `pccap`    | —        | **Yes**  | —       | **Yes**    | —        |
| `psiz`     | —        | **Yes**  | —       | **Yes**    | —        |
| `rcomp`    | —        | **Yes**  | —       | **Yes**    | —        |
| `zcomp`    | **Yes**  | **Yes**  | **Yes** | **Yes**    | —        |
| `zkey`     | **Yes**  | **Yes**  | **Yes** | **Yes**    | —        |

## API

## `/publish`

Publishes a message to the given queue. If there is a connection waiting to recieve a message from this queue, the message would be transfered directly to the awaiting connection.

Publication can be done via both `GET` and `POST` methods. To publish a message via `GET` method use the `queue` argument to pass queue name and the `msg` argument to pass message body. To publish a message via `POST` method pass message body via request body instead of the `msg` argument.

Server will respond with `OK` message.

**Example:**
```
/publish?queue=urgent&msg=Process+this+message+as+soon+as+possible!
```
Response:
```
OK
```

## `/subscribe`

Tries to fetch a message from one of the queues given. If there is a message at least in one of these queues, the message will be removed from the queue and returned as response body. The name of the queue from which the message was taken from will be provided inside a `Queue` response header.

Subscription is always done via `GET` method. To fetch a message from a queue use the name of the queue as the `queues` argument value. Multiple queue names could be passed separated with the `,` (quote) character.

**Example:**
```
/subscribe?queues=urgent,someday
```
Response:
```
Process this message as soon as possible!
```

## `/status`

Displays information about the queues, their messages and current subscriptions encoded in JSON format.

**Example:**
```
/status
```
Response:
```json
{
    "urgent": {
        "messages": 0,
        "subscriptions": 0
    },
    "someday": {
        "messages": 0,
        "subscriptions": 0
    }
}
```

## `/debug`

Displays debug information about the queue process. Currenty displays the number of goroutines only.

**Example:**
```
/debug
```
Response:
```json
{
    "goroutines": 13
}
```

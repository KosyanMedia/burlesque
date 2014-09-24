# Burlesque

Burlesque is a [message queue](http://en.wikipedia.org/wiki/Message_queue) server. It gives access to queues using the [pub/sub HTTP API](#api).

This server's main purpose is to provide an inter-process comunication tool with a memory efficient persistent storage for messages. These messages usually are delayed job descriptions serialized in JSON that are published by the application server and later on retrieved by application workers.

Subscription uses [long polling](http://en.wikipedia.org/wiki/Push_technology#Long_polling) technique. When application worker subscribes to a queue which is empty at the moment, connection is kept open until a client publishes a message to this queue, or subscription timeout is reached. If there is already a message in the queue it is removed from the queue and returned to the worker.

To store messages Burlesque uses [Kyoto Cabinet](http://fallabs.com/kyotocabinet/), which is a powerful DIY database. Usage of Kyoto Cabinet is thoroughly described in the [storage](#storage) section of this document.

#### Contents

* [Installation](#installation)
  * [Building on OS X](#building-on-os-x)
* [Starting](#starting)
* [Storage](#storage)
  * [In-memory databases](#in-memory-databases)
  * [Persistent databases](#persistent-databases)
  * [Tuning parameters](#tuning-parameters)
  * [Support of tuning parameters by databases](#support-of-tuning-parameters-by-databases)
* [What storage to choose](#what-storage-to-choose)
  * [Production](#production)
  * [Development](#development)
* [API](#api)
  * [Publish](#publish)
  * [Subscribe](#subscribe)
  * [Flush](#flush)
  * [Status](#status)
  * [Debug](#debug)

## Installation

Download and extract the [latest release](https://github.com/KosyanMedia/burlesque/releases). That's it.

### Building on OS X

First install [Homebrew](http://brew.sh/). Using Homebrew install Go language compiler and tools. Then install Kyoto Cabinet library.

```
brew install go
brew install kyoto-cabinet
go get github.com/KosyanMedia/burlesque
```

## Starting

The following arguments are supported by the `burlesque` executable:

| Argument | Description | Defaults |
| -------- | ----------- | -------- |
| `-storage` | Kyoto Cabinet storage path (e.g. `storage.kch#msiz=524288000`) | `-` |
| `-environment` | Process environment: `development` or `production` | `development` |
| `-port` | Server HTTP port | `4401` |
| `-rollbar` | [Rollbar](https://rollbar.com/) token | ||

#### Example

```bash
wget -O burlesque.zip https://github.com/KosyanMedia/burlesque/archive/1.0.0.zip
unzip burlesque.zip
./burlesque
```

By default, Burlesque starts on port `4401` and uses in-memory database `ProtoHashDB`.

## Storage

`-storage` argument defines the way data is stored in the database. You can read more on Kyoto Cabinet database types [here](http://fallabs.com/kyotocabinet/spex.html#tutorial_dbchart).

### In-memory databases

If you need a temporary in-memory storage use the following symbols as the `-storage` value:

| Value | Database Type |
| ----- | ------------- |
| `-` | `ProtoHashDB` Prototype hash database. In-memory database implemented with `std::unorderd_map` |
| `+` | `ProtoTreeDB` Prototype tree database. In-memory database implemented with `std::map` |
| `:` | `StashDB` Stash database. In-memory database saving memory |
| `*` | `CacheDB` Cache hash database. In-memory database featuring [LRU](http://en.wikipedia.org/wiki/Cache_algorithms#Examples) deletion |
| `%` | `GrassDB` Cache tree database. In-memory database of B+ tree: cache with order |

#### Example: `-`

### Persistent databases

In order to use a persistent database, use the path to the database file (or directory) as the `-storage` argument value. File extension in the database path defines the type of the database created.

| File Extension | Database Type |
| -------------- | ------------- |
| `kch` | `HashDB` File hash database. File database of hash table: typical DBM |
| `kct` | `TreeDB` File tree database. File database of B+ tree: DBM with order |
| `kcd` | `DirDB` Directory hash database. Respective files in a directory of the file system |
| `kcf` | `ForestDB` Directory tree database. Directory database of B+ tree: huge DBM with order |
| `kcx` | `TextDB` Plain text database. Emulation to handle a plain text file as a database |

#### Example: `/path/to/my/storage.kch`

### Tuning parameters

When the database type is defined, you can also add [tuning parameters](http://fallabs.com/kyotocabinet/spex.html#tips) to the `-storage` argument. Tuning parameters are separated by the `#` symbol, parameters' name and value are separated by the `=` symbol.

The table below describes tuning parameters.

| Parameter  | Description |
| ---------- | ----------- |
| `apow`     | Power of the record size alignment |
| `bnum`     | Base hash table size (number of buckets of the hash table) |
| `capcnt`   | Capacity limit by the number of records (`#capcnt=10000` means "keep in memory 10,000 records maximum") |
| `capsiz`   | Capacity limit by the size of records (`#capsiz=536870912` means "keep in memory all the records that fit into 512 megabytes") |
| `dfunit`   | Unit step number of auto defragmentation (`#dfunit=8` means "run defragmentation every 8 fragmentations detected"). |
| `fpow`     | Power of the free block pool capacity |
| `log`      | Path to the log file. Use `-` for the STDOUT, or `+` for the STDERR |
| `logkinds` | Kinds of logged messages. The value can be `debug`, `info`, `warn` or `error` |
| `logpx`    | Prefix of each log message |
| `msiz`     | Expected database memory usage |
| `opts`     | Additional options: `s`, `l` and `c` (can be specified together, e.g `lc`). `s` stands for "small" and reduces the width of record address from 6 bytes to 4 bytes. As a result, the footprint for each record is reduced from 16 bytes to 12 bytes. However, it limits the maximum size of the database file to 16GB. `l` stands for "linear" and changes the data structure of the collision chain of hash table from binary tree to linear linked list. `c` enables compression of the record values. If the value is bigger than 1KB compression is effective. |
| `pccap`    | Capacity size of the page cache |
| `psiz`     | Page size |
| `rcomp`    | Comparator used to compare key names. `lex` for the lexical comparator, `dec` for the decimal comparator, `lexdesc` for the lexical descending comparator, or `decdesc` for the decimal descending comparator |
| `zcomp`    | Compression library: `zlib` for the [ZLIB raw](http://en.wikipedia.org/wiki/Zlib#Encapsulation) compressor, `def` for the ZLIB [deflate](http://en.wikipedia.org/wiki/DEFLATE) compressor, `gz` for the ZLIB [gzip](http://en.wikipedia.org/wiki/Gzip) compressor, `lzo` for the [LZO](http://en.wikipedia.org/wiki/Lempel%E2%80%93Ziv%E2%80%93Oberhumer) compressor, `lzma` for the [LZMA](http://en.wikipedia.org/wiki/Lempel%E2%80%93Ziv%E2%80%93Markov_chain_algorithm) compressor, or `arc` for the [Arcfour](http://en.wikipedia.org/wiki/RC4) cipher |
| `zkey`     | Cipher keyword used with compression |

#### Example: `storage.kch#opts=c#zcomp=gz#msiz=524288000`

### Support of tuning parameters by databases

The table below describes support of these parameters by the **in-memory** database types.

| Parameter  | `ProtoHashDB` | `ProtoTreeDB` | `StashDB` | `CacheDB` | `GrassDB` |
| ---------- | :-----------: | :-----------: | :-------: | :-------: | :-------: |
| `bnum`     |               |               | •         | •         | •         |
| `capcnt`   |               |               |           | •         |           |
| `capsiz`   |               |               |           | •         |           |
| `log`      | •             | •             | •         | •         | •         |
| `logkinds` | •             | •             | •         | •         | •         |
| `logpx`    | •             | •             | •         | •         | •         |
| `opts`     |               |               |           | •         | •         |
| `pccap`    |               |               |           |           | •         |
| `psiz`     |               |               |           |           | •         |
| `rcomp`    |               |               |           |           | •         |
| `zcomp`    |               |               |           | •         | •         |
| `zkey`     |               |               |           | •         | •         |

The table below describes support of these parameters by the **persistent** database types.

| Parameter  | `HashDB` | `TreeDB` | `DirDB` | `ForestDB` | `TextDB` |
| ---------- | :------: | :------: | :-----: | :--------: | :------: |
| `apow`     | •        | •        |         |            |          |
| `bnum`     | •        | •        |         |            |          |
| `dfunit`   | •        | •        |         |            |          |
| `fpow`     | •        | •        |         |            |          |
| `log`      | •        | •        | •       | •          | •        |
| `logkinds` | •        | •        | •       | •          | •        |
| `logpx`    | •        | •        | •       | •          | •        |
| `msiz`     | •        | •        |         |            |          |
| `opts`     | •        | •        | •       | •          |          |
| `pccap`    |          | •        |         | •          |          |
| `psiz`     |          | •        |         | •          |          |
| `rcomp`    |          | •        |         | •          |          |
| `zcomp`    | •        | •        | •       | •          |          |
| `zkey`     | •        | •        | •       | •          |          |

## What storage to choose

### Production

For production usage it is strongly recommended to choose a **persistent** database. Burlesque uses Kyoto Cabinet as a persistent hash-table, which means `HashDB` would be a smart choice.

If the average message size is expected to be more than 1KB then compression should be considered as an option. To enable compression you need to pass `opts` tuning parameter to the database path with value `c` (`#opts=c`) in it, you also need to define compression algorithm using the `zcomp` parameter (e.g `#zcomp=gz`).

You can define maximum memory limit; when the limit is reached new records are swapped to disk. Memory limit is defined by value of `msiz` parameter in bytes (e.g `#msiz=524288000`)

So, to use a persistent hash database with enabled compression and 512MB memory limit the `-storage` argument value should be `storage.kch#opts=c#zcomp=gz#msiz=524288000`.

#### Further tuning

If queues are kept empty or relatively small, `bnum` option might be considered (e.g `#bnum=1000`)

### Development

If development database doesn't need to be persisted consider using `ProtoHashDB` (which locks the whole table), `StashDB` (locks record) or `CacheDB` (locks record using a mutex). By default `ProtoHashDB` is used.

## API

All endpoints exposed by the API are described below.

## Publish

This endpoint is used to publish messages to a queue. If there is a connection waiting to recieve a message from this queue, the message will be handed directly to the awaiting worker.

Publication can be done via both `GET` and `POST` methods. Both methods use `queue` argument to pass queue name. When using `GET` method pass message body with `msg` argument. To publish a message via `POST` method pass message body via request body instead of the `msg` argument.

In case of success, server will respond with status 200 and `OK` message. Otherwise, there will be status 500 and `FAIL` message.

#### Example
```bash
$ curl '127.0.0.1:4401/publish?queue=urgent' -d \
  'Process this message as soon as possible!'
```
Response
```
OK
```

## Subscribe

This endpoint is used to try and fetch a message from one of the queues given. If at least one of these queues contains a message, this message will be removed from the queue and returned as a response body. The name of the queue where this message was taken from will be provided as a `Queue` response header.

Subscription is always done via `GET` method. To fetch a message from a queue use the name of the queue as the `queues` argument value. Multiple queue names could be passed separated with the comma character.

#### Example
```bash
$ curl '127.0.0.1:4401/subscribe?queues=urgent,someday'
```
Response
```
Process this message as soon as possible!
```

## Flush

This endpoint is used to fetch all messages from all of the given queues. All messages are encoded into a single JSON document.

#### Example
```bash
$ curl '127.0.0.1:4401/flush?queues=urgent,someday' > dump.json
$ cat dump.json
```
Result
```json
[
    {
        "queue": "urgent",
        "message": "Process this message as soon as possible!"
    },
    {
        "queue": "someday",
        "message": "Process this message in your spare time"
    }
]
```

## Status

This endpoint is used to display information about the queues, their messages and current subscriptions encoded in JSON format.

#### Example
```bash
$ curl '127.0.0.1:4401/status'
```
Response
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

## Debug

This endpoint is used to display debug information about Burlesque process. Currenty displays the number of goroutines only.

#### Example
```bash
$ curl '127.0.0.1:4401/debug'
```
Response
```json
{
    "gomaxprocs": 1,
    "goroutines": 12,
    "kyoto_cabinet": {
        "count": 0,
        "path": "-",
        "realtype": 16,
        "size": 0,
        "type": 16
    },
    "version": "0.2.0"
}
```

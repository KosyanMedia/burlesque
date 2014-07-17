# Burlesque

Burlesque is a [message processing queue](http://en.wikipedia.org/wiki/Message_queue) writen in [Go](http://golang.org/). It exposes queues using the [pub/sub HTTP API](#api).

The general purpose of this queue is to provide tool for inter-process comutication with a memory efficient persisted storage for messages (usually a delayed job description serialized in JSON) published by the application server and later retrieved by other application workers.

Subscription is done using [long polling](http://en.wikipedia.org/wiki/Push_technology#Long_polling) technique. When application worker subscribes to a queue which is empty at the moment, connection is kept open until another client publishes a message to this queue, or the first client disconnects. If there is a message in the queue it will be removed from the queue and returned to the client.

Burlesque uses [Kyoto Cabinet](http://fallabs.com/kyotocabinet/) to store messages, which is a powerfull DIY database. Usage of Kyoto Cabinet is thoroughly described in the [storage](#storage) section of this document.

#### Contents

* [Installation](#installation)
  * [Building on OSX](#building-on-osx)
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
  * [Status](#status)
  * [Debug](#debug)

## Installation

Download and extract the [latest release](https://github.com/KosyanMedia/burlesque/releases). That's it.

### Building on OSX

First install [Homebrew](http://brew.sh/). Using Homebrew install Go language compiler and tools. Then install Kyoto Cabinet library.

```
brew install go
brew install kyoto-cabinet
go get github.com/KosyanMedia/burlesque
```

## Starting

Use the following arguments to the `burlesque` executable:

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

By default Burlesque starts on port `4401` in development mode and uses in-memory database `ProtoHashDB`.

## Storage

`-storage` argument defines a way the data will be stored into a database. You can read more on Kyoto Cabinet database types [here](http://fallabs.com/kyotocabinet/spex.html#tutorial_dbchart).

### In-memory databases

If you need a temporary in-memory storage use the following symbols as the `-storage` value:

| Value | Database Type |
| ----- | ------------- |
| `-` | `ProtoHashDB` Prototype hash database. On-memory database implemented with `std::unorderd_map` |
| `+` | `ProtoTreeDB` Prototype tree database. On-memory database implemented with `std::map` |
| `:` | `StashDB` Stash database. On-memory database saving memory |
| `*` | `CacheDB` Cache hash database. On-memory database featuring [LRU](http://en.wikipedia.org/wiki/Cache_algorithms#Examples) deletion |
| `%` | `GrassDB` Cache tree database. On-memory database of B+ tree: cache with order |

#### Example: `-`

### Persistent databases

In order to use a persistent database use the path to the database file (or directory) as the `-storage` argument value. File extension in the database path defines the type of the database created.

| File Extension | Database Type |
| -------------- | ------------- |
| `kch` | `HashDB` File hash database. File database of hash table: typical DBM |
| `kct` | `TreeDB` File tree database. File database of B+ tree: DBM with order |
| `kcd` | `DirDB` Directory hash database. Respective files in a directory of the file system |
| `kcf` | `ForestDB` Directory tree database. Directory database of B+ tree: huge DBM with order |
| `kcx` | `TextDB` Plain text database. Emulation to handle a plain text file as a database |

#### Example: `/path/to/my/storage.kch`

### Tuning parameters

In addition to defining database type you can also add [tuning parameters](http://fallabs.com/kyotocabinet/spex.html#tips) to the `-storage` argument. Tuning parameters are separated by the `#` symbol, parameters' name and value are separated by the `=` symbol.

The table below describes tuning parameters.

| Parameter  | Description |
| ---------- | ----------- |
| `apow`     | Power of the alignment of record size |
| `bnum`     | Base hash table size (number of buckets of the hash table) |
| `capcnt`   | Capacity limit by the number of records (`#capcnt=10000` means "keep in memory 10,000 records maximum) |
| `capsiz`   | Capacity limit by the size of records (`#capsiz=536870912` means "keep in memory all the records that fit into 512 megabytes) |
| `dfunit`   | Unit step number of auto defragmentation (`#dfunit=8` means "run defragmentation every 8 fragmentations detected"). |
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

#### Example: `storage.kch#opts=c#zcomp=gz#msiz=524288000`

### Support of tuning parameters by databases

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

## What storage to choose

### Production

For production usage it is strongly recommended to choose a **persistent** database. Internally Burlesque uses Kyoto Cabinet as a persisted hash-table, so using `HashDB` would be a smart choise.

If the average message size expected to be more than 1KB then compression should be considered as an option. To enable compression you need to pass `opts` tuning parameter to the database path with value `c` (`#opts=c`), you also need to define compression algorithm using the `zcomp` parameter (e.g `#zcomp=gz`).

You can define maximum memory limit; when the limit is reached new records are swapped to disk. Memory limit is defined by passing `msiz` parameter with value in bytes (e.g `#msiz=524288000`)

So, to use a persisted hash database with enabled compression and 512MB memory limit the `-storage` argument value is `storage.kch#opts=c#zcomp=gz#msiz=524288000`.

#### Further tuning

If queues are kept empty all at relatively small size, `bnum` option might be considered (e.g `#bnum=1000`)

### Development

If development database don't need to be persisted consider using `ProtoHashDB` (which locks the whole table), `StashDB` (locks record) or `CacheDB` (locks record using a mutex). By default `ProtoHashDB` is used.

## API

All endpoints exposed by the API are described below.

## Publish

Publishes a message to the given queue. If there is a connection waiting to recieve a message from this queue, the message would be transfered directly to the awaiting connection.

Publication can be done via both `GET` and `POST` methods. To publish a message via `GET` method use the `queue` argument to pass queue name and the `msg` argument to pass message body. To publish a message via `POST` method pass message body via request body instead of the `msg` argument.

Server will respond with `OK` message.

#### Example
```bash
curl '127.0.0.1:4401/publish?queue=urgent' -d 'Process this message as soon as possible!'
```
Response
```
OK
```

## Subscribe

Tries to fetch a message from one of the queues given. If there is a message at least in one of these queues, the message will be removed from the queue and returned as response body. The name of the queue from which the message was taken from will be provided inside a `Queue` response header.

Subscription is always done via `GET` method. To fetch a message from a queue use the name of the queue as the `queues` argument value. Multiple queue names could be passed separated with the comma character `,`.

#### Example
```bash
curl '127.0.0.1:4401/subscribe?queues=urgent,someday'
```
Response
```
Process this message as soon as possible!
```

## Status

Displays information about the queues, their messages and current subscriptions encoded in JSON format.

#### Example
```bash
curl 127.0.0.1:4401/status
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

Displays debug information about the queue process. Currenty displays the number of goroutines only.

#### Example
```basg
curl 127.0.0.1:4401/debug
```
Response
```json
{
    "goroutines": 13
}
```

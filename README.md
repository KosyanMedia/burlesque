# Burlesque

Burlesque is a [message processing queue](http://en.wikipedia.org/wiki/Message_queue) writen in Go. It exposes queues using the HTTP API that allows publishing messages and subscribing to them. See the [API](#API) section for more details.

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
| `-storage` | Kyoto Cabinet storage path (e.g. `storage.kch#zcomp=gz#capsiz=524288000`) | `-` |
| `-environment` | Process environment: `development` or `production` | `development` |
| `-port` | Server HTTP port | `4401` |
| `-rollbar` | [Rollbar](https://rollbar.com/) token | ||

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
```
{
	"goroutines": 13
}
```
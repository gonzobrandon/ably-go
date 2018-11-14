## [Ably Go](https://www.ably.io)

A Go client library for [www.ably.io](https://ably.io), the realtime messaging service.

## Feature support

This library targets the Ably 1.1 client library specification.

The 1.1 specification includes multiple new features, and this library implements a subset of those features as follows.

| Feature | Spec reference | Supported |
| --- | --- | --- |
| Push notifications admin API | RSH1 | No |
| Push notifications target API | RSH2 | No |
| JWT authentication | multiple | No |
| Idempotent REST publishing | RSL1j, RSL1k | Yes |
| Fallback host affinity | RSC15f | Yes |
| `ErrorInfo.href` help | TI4, TI5 | Yes |
| request() enhancements | RSC19 | Yes |
| Transient realtime publishing | RTL6c | No |
| Exception reporting | RSC20 | No |

The 1.1 specification also contains numerous spec clarifications, corrections and minor feature additions.
Refer to [complete diff for the 1.1 specification](https://github.com/ably/docs/tree/master/content/client-lib-development-guide/versions/features-1-0__1-1.diff)
for more details.

It is intended that this library is upgraded incrementally, with 1.1 feature support expanded in successive minor
releases. If there are features that are currently missing that are a high priority for your use-case then please
[contact Ably customer support](https://support.ably.io). Pull Requests are also welcomed.

## Installation

```bash
~ $ go get -u github.com/ably/ably-go/ably
```

## Using the Realtime API

### Creating a client

```go
client, err := ably.NewRealtimeClient(ably.NewClientOptions("xxx:xxx"))
if err != nil {
	panic(err)
}

channel := client.Channels.Get("test")
```

### Subscribing to a channel for all events

```go
sub, err := channel.Subscribe()
if err != nil {
	panic(err)
}

for msg := range sub.MessageChannel() {
	fmt.Println("Received message:", msg)
}
```

### Subscribing to a channel for `EventName1` and `EventName2` events

```go
sub, err := channel.Subscribe("EventName1", "EventName2")
if err != nil {
	panic(err)
}

for msg := range sub.MessageChannel() {
	fmt.Println("Received message:", msg)
}
```

### Publishing to a channel

```go
// send request to a server
res, err := channel.Publish("EventName1", "EventData1")
if err != nil {
	panic(err)
}

// await confirmation
if err = res.Wait(); err != nil {
	panic(err)
}
```

### Announcing presence on a channel

```go
// send request to a server
res, err := channel.Presence.Enter("presence data")
if err != nil {
	panic(err)
}

// await confirmation
if err = res.Wait(); err != nil {
	panic(err)
}
```

### Announcing presence on a channel on behalf of other client

```go
// send request to a server
res, err := channel.Presence.EnterClient("clientID", "presence data")
if err != nil {
	panic(err)
}

// await confirmation
if err = res.Wait(); err != nil {
	panic(err)
}
```

### Getting all clients present on a channel

```go
clients, err := channel.Presence.Get(true)
if err != nil {
	panic(err)
}

for _, client := range clients {
	fmt.Println("Present client:", client)
}
```

### Subscribing to all presence messages

```go
sub, err := channel.Presence.Subscribe()
if err != nil {
	panic(err)
}

for msg := range sub.PresenceChannel() {
	fmt.Println("Presence event:", msg)
}
```

### Subscribing to 'Enter' presence messages only

```go
sub, err := channel.Presence.Subscribe(proto.PresenceEnter)
if err != nil {
	panic(err)
}

for msg := range sub.PresenceChannel() {
	fmt.Println("Presence event:", msg)
}
```

## Using the REST API

### Introduction

All examples assume a client and/or channel has been created as follows:

```go
client, err := ably.NewRestClient(ably.NewClientOptions("xxx:xxx"))
if err != nil {
	panic(err)
}

channel := client.Channel("test")
```

### Publishing a message to a channel

```go
err = channel.Publish("HelloEvent", "Hello!")
if err != nil {
	panic(err)
}
```

### Querying the History

```go
page, err := channel.History(nil)
for ; err == nil; page, err = page.Next() {
	for _, message := range page.Messages() {
		fmt.Println(message)
	}
}
if err != nil {
	panic(err)
}
```

### Presence on a channel

```go
page, err := channel.Presence.Get(nil)
for ; err == nil; page, err = page.Next() {
	for _, presence := range page.PresenceMessages() {
		fmt.Println(presence)
	}
}
if err != nil {
	panic(err)
}
```

### Querying the Presence History

```go
page, err := channel.Presence.History(nil)
for ; err == nil; page, err = page.Next() {
	for _, presence := range page.PresenceMessages() {
		fmt.Println(presence)
	}
}
if err != nil {
	panic(err)
}
```

### Generate Token and Token Request

```go
client.Auth.RequestToken()
client.Auth.CreateTokenRequest()
```

### Fetching your application's stats

```go
page, err := client.Stats(&ably.PaginateParams{})
for ; err == nil; page, err = page.Next() {
	for _, stat := range page.Stats() {
		fmt.Println(stat)
	}
}
if err != nil {
	panic(err)
}
```

## Known limitations (work in progress)

As the library is actively developed couple of features are not there yet:

- Realtime connection recovery is not implemented
- Realtime connection failure handling is not implemented
- ChannelsOptions and CipherParams are not supported when creating a Channel
- Realtime Ping function is not implemented

## Release process

This library uses [semantic versioning](http://semver.org/). For each release, the following needs to be done:



* Create a branch for the release, named like `release-1.0.6`
* Replace all references of the current version number with the new version number and commit the changes
* Run [`github_changelog_generator`](https://github.com/skywinder/Github-Changelog-Generator) to update the [CHANGELOG](./CHANGELOG.md): `github_changelog_generator -u ably -p ably-go --header-label="# Changelog" --release-branch=release-1.0.6 --future-release=v1.0.6` 
* Commit [CHANGELOG](./CHANGELOG.md)
* Add a tag and push to origin such as `git tag v1.0.6; git push origin v1.0.6`
* Make a PR against `develop`
* Once the PR is approved, merge it into `develop`
* Fast-forward the master branch: `git checkout master && git merge --ff-only develop && git push origin master`


## Support and feedback

Please visit http://support.ably.io/ for access to our knowledgebase and to ask for any assistance.

You can also view the [community reported Github issues](https://github.com/ably/ably-go/issues).

## Contributing

Because this package uses `internal` packages, all fork development has to happen under `$GOPATH/src/github.com/ably/ably-go` to prevent `use of internal package not allowed` errors.

1. Fork `github.com/ably/ably-go`
2. go to the `ably-go` directory: `cd $GOPATH/src/github.com/ably/ably-go`
3. add your fork as a remote: `git remote add fork git@github.com:your-username/ably-go`
4. create your feature branch: `git checkout -b my-new-feature`
5. commit your changes (`git commit -am 'Add some feature'`)
6. ensure you have added suitable tests and the test suite is passing: `make test`
7. push to the branch: `git push fork my-new-feature`
8. create a new Pull Request

## License

Copyright (c) 2016 Ably Real-time Ltd, Licensed under the Apache License, Version 2.0.  Refer to [LICENSE](LICENSE) for the license terms.

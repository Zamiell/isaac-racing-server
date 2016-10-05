isaac-racing-server
=================

Description
-----------

This is the server software for the Binding of Isaac: Afterbirth+ racing mod. Normally a single player game, the mod and server allow players to be able to race each other in real time.

The server is written in [Go](https://golang.org/) and uses WebSockets to communicate with the client. It leverages [Auth0](https://auth0.com/) for authentication and uses a [SQLite](https://sqlite.org/) database to keep track of the races.



Install
-------

* Install Go (you need to be able to run the `go` command).
* Install SQLite3 (you need to be able to run the `sqlite3` command).
* `go get github.com/Zamiell/isaac-racing-server`
* `cd $GOPATH/Zamiell/isaac-racing-server`
* `sqlite3 database.sqlite < install/database_schema.sql`
* Open up the `main.go` file and change the constants near the top of the file to your liking.
* Create a `.env` file in the current directory with the following contents:

```
SESSION_SECRET=some_long_random_string
AUTH0_CLIENT_ID=the_client_id_from_auth0
AUTH0_CLIENT_SECRET=the_client_secret_from_auth0
```



Run
---

* `cd $GOPATH/Zamiell/isaac-racing-server`
* `go run *.go`

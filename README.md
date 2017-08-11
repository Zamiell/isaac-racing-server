isaac-racing-server
===================

Additional Information
----------------------

If you are not a developer, please visit [the website for Racing+](https://isaacracing.net/).

<br />



Description
-----------

This is the server software for Racing+, a Binding of Isaac: Afterbirth+ racing platform. Normally a single player game, the Lua mod, client, and server allow players to be able to race each other in real time.

The server is written in [Go](https://golang.org/) and uses WebSockets to communicate with the client. It leverages [Auth0](https://auth0.com/) for authentication and uses a [SQLite](https://sqlite.org/) database to keep track of the races.

You may also be interested in [the client repository](https://github.com/Zamiell/isaac-racing-client) or [the Lua mod](https://github.com/Zamiell/isaac-racing-client/tree/master/mod).

<br />



Code Layout
-----------

* Logging stuff is in the `logger` directory.
* Twitch.tv stuff is in the `twitch` directory.
* Discord stuff is in the `discord` directory.
* WebSocket command logic is in the `websocket` directory.
* Database logic is in the `models` directory.
* HTML templates are in the `views` directory.
* Webpage logic is in the `controllers` directory.

<br/>



Install
-------

These instructions assume you are running Ubuntu 16.04 LTS. Some adjustment will be needed for Windows installations.

* Install Go:
  * `sudo add-apt-repository ppa:longsleep/golang-backports`
  * `sudo apt update`
  * `sudo apt install golang-go -y`
* Install SQLite3:
   * `sudo apt install sqlite3 -y`
* Clone the server:
  * `go get github.com/Zamiell/isaac-racing-server`
  * `cd $GOPATH/Zamiell/isaac-racing-server`
* Set up the database:
  * `sqlite3 database.sqlite < install/database_schema.sql`
* Set up the configuration:
  * `nano main.go` (change the constants near the top of the file to your liking)
* Set up the environment values:
  * `cp .env_template .env`
  * `nano .env` (fill in the values)

<br />



Run
---

* `cd $GOPATH/Zamiell/isaac-racing-server`
* `go run *.go`

<br />




Compile / Build
---------------

* `go install` (this creates `$GOPATH/bin/isaac-racing-server`)

<br />



Install HTTPS (optional)
------------------------

* `apt-install letsencrypt`
* `letsencrypt certonly --standalone -d isaacracing.net -d www.isaacracing.net` (this creates `/etc/letsencrypt/live/isaacracing.net`)

Later, to renew the certificate:

* `RENEW_DIR=/root/isaac-racing-server/letsencrypt && mkdir -p $RENEW_DIR && letsencrypt renew --webroot --webroot-path $RENEW_DIR && rm -rf $RENEW_DIR`

<br />



Install as a service (optional)
-------------------------------

* Install Supervisor (for example, on Ubuntu 16.04):
  * `apt install supervisor`
  * `systemctl enable supervisor` (http://unix.stackexchange.com/questions/281774/ubuntu-server-16-04-cannot-get-supervisor-to-start-automatically)
* Copy the configuration files:
  * `cp $GOPATH/Zamiell/isaac-racing-server/install/supervisord/supervisord.conf /etc/supervisord/supervisord.conf`
  * `cp $GOPATH/Zamiell/isaac-racing-server/install/supervisord/isaac-racing-server.conf /etc/supervisord/conf.d/isaac-racing-server.conf`
* Start it: `systemctl start supervisor`

Later, to manage the service:

* Start it: `supervisorctl start isaac-racing-server`
* Stop it: `supervisorctl stop isaac-racing-server`
* Restart it: `supervisorctl restart isaac-racing-server`

<br />

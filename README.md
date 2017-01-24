isaac-racing-server
=================

Additional Information
----------------------

If you are not a developer, please visit [the website for Racing+](https://isaacracing.net/).



Description
-----------

This is the server software for Racing+, a Binding of Isaac: Afterbirth+ racing mod. Normally a single player game, the mod and server allow players to be able to race each other in real time.

The server is written in [Go](https://golang.org/) and uses WebSockets to communicate with the client. It leverages [Auth0](https://auth0.com/) for authentication and uses a [SQLite](https://sqlite.org/) database to keep track of the races.



Install
-------

* Install Go (you need to be able to run the `go` command).
* Install SQLite3 (you need to be able to run the `sqlite3` command).
   * On Ubuntu: `apt install sqlite3`
* `go get github.com/Zamiell/isaac-racing-server`
* `cd $GOPATH/Zamiell/isaac-racing-server`
* `sqlite3 database.sqlite < install/database_schema.sql`
* `nano main.go`
  * Change the constants near the top of the file to your liking.
* `cp .env_template .env && nano .env`
  * Fill in the values.



Run
---

* `cd $GOPATH/Zamiell/isaac-racing-server`
* `go run *.go`



Compile / Build
---------------

* `go install` (this creates `$GOPATH/bin/isaac-racing-server`)



Install HTTPS (optional)
------------------------

* `apt-install letsencrypt`
* `letsencrypt certonly --standalone -d isaacracing.net -d www.isaacracing.net` (this creates `/etc/letsencrypt/live/isaacracing.net`)

Later, to renew the certificate:

* `RENEW_DIR=/root/isaac-racing-server/letsencrypt && mkdir -p $RENEW_DIR && letsencrypt renew --webroot --webroot-path $RENEW_DIR && rm -rf $RENEW_DIR`



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

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



Install
-------

These instructions assume you are running Ubuntu 16.04 LTS. Some adjustment will be needed for Windows installations.

* Install Go:
  * `sudo add-apt-repository ppa:longsleep/golang-backports`
  * `sudo apt update`
  * `sudo apt install golang-go -y`
* Install [MariaDB](https://mariadb.org/) and set up a user:
  * `sudo apt install mariadb-server -y`
  * `sudo mysql_secure_installation`
    * Follow the prompts.
  * `sudo mysql -u root -p`
    * `CREATE DATABASE isaac;`
    * `CREATE USER 'isaacuser'@'localhost' IDENTIFIED BY '1234567890';` (change the password to something else)
    * `GRANT ALL PRIVILEGES ON isaac.* to 'isaacuser'@'localhost';`
* Clone the server:
  * `mkdir -p $GOPATH/Zamiell`
  * `cd $GOPATH/Zamiell/`
  * `git clone https://github.com/Zamiell/isaac-racing-server.git` (or clone a fork, if you are doing development work)
* Set up environment variables:
  * `cp .env_defaults .env`
  * `nano .env`
    * Change the `DB_HOST`, `DB_USER`, and `DB_PASS` values accordingly.
    * Create a random 64 digit alphanumeric string for `SESSION_SECRET`.
    * `SENTRY_SECRET`, `TWITCH_OAUTH`, and `DISCORD_TOKEN` can be left blank.
* Import the database schema:
  * `mysql -uisaacuser -p1234567890 < install/database_schema.sql` (change the password accordingly)
* Set up the some configuration variables:
  * `nano src/main.go` (change the constants near the top of the file to your liking)

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

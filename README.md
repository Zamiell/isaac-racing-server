# isaac-racing-server

If you are not a developer, please visit [the website for Racing+](https://isaacracing.net/).

<br />

## Description

This is the server software for Racing+, a Binding of Isaac: Repentance racing platform. Normally a single player game, the Lua mod, client, and server allow players to be able to race each other in real time.

The server is written in [Go](https://golang.org/) and uses WebSockets to communicate with the client. It leverages the [Steam API](https://partner.steamgames.com/doc/webapi_overview) for authentication and uses a [MariaDB](https://mariadb.com/) database to keep track of the races.

You may also be interested in [the client repository](https://github.com/Zamiell/isaac-racing-client) or [the Lua mod](https://github.com/Zamiell/isaac-racing-client/tree/master/mod).

<br />

## Install

These instructions assume you are running Ubuntu 16.04 LTS. Some adjustment will be needed for Windows or MacOS installations.

- Install [Go](https://golang.org/):
  - `sudo add-apt-repository ppa:longsleep/golang-backports` (if you don't do this, it will install a version of Go that is very old)
  - `sudo apt update`
  - `sudo apt install golang-go -y`
  - `mkdir "$HOME/go"`
  - `export GOPATH=$HOME/go && echo 'export GOPATH=$HOME/go' >> ~/.profile`
  - `export PATH=$PATH:$GOPATH/bin && echo 'export PATH=$PATH:$GOPATH/bin' >> ~/.profile`
- Install [MariaDB](https://mariadb.org/):
  - `sudo apt install mariadb-server -y`
  - `sudo mysql_secure_installation`
    - Follow the prompts.
- Set up a MariaDB database and a MariaDB user:
  - `sudo mysql -u root -p`
    - `CREATE DATABASE isaac;`
    - `CREATE USER 'isaacuser'@'localhost' IDENTIFIED BY '1234567890';` (change the password to something else)
    - `GRANT ALL PRIVILEGES ON isaac.* to 'isaacuser'@'localhost';`
    - `FLUSH PRIVILEGES;`
- Clone the repository:
  - `cd [the path where you want the code to live]` (optional)
  - If you already have an SSH key pair and have the public key attached to your GitHub profile, then use the following command to clone the repository via SSH:
    - `git clone git@github.com:Zamiell/isaac-racing-server.git`
  - If you do not already have an SSH key pair, then use the following command to clone the repository via HTTPS:
    - `git clone https://github.com/Zamiell/isaac-racing-server.git`
  - Or, if you are doing development work, then clone your forked version of the repository. For example:
    - `git clone https://github.com/[Your_Username]/isaac-racing-server.git`
- Enter the cloned repository:
  - `cd isaac-racing-server`
- Build the server, which will automatically download install all of the Go dependencies:
  - `./build.sh`
- Set up environment variables:
  - `cp .env_template .env`
  - `nano .env`
    - Create a random 64 digit alphanumeric string for `SESSION_SECRET`.
    - Change the `DB_PASS` value accordingly.
    - If you want to be able to login to the WebSocket server, set a value for `STEAM_WEB_API_KEY`. (You can get it from the [Steam community portal](https://steamcommunity.com/dev/apikey).)
    - The rest of the values can be left blank.
- Import the database schema:
  - `mysql -uisaacuser -p < install/database_schema.sql` <!-- cspell:disable-line -->

<br />

## Run

- To re-compile and run the server, simply run the `run.sh` script.
  - The re-compiled binary is called `isaac-racing-server` and is located in the root of the repository.
- If you are on Linux or MacOS, sudo might be necessary because the server listens on port 80 and/or 443.

<br />

## Install HTTPS (optional)

- `apt-install letsencrypt`
- `letsencrypt certonly --standalone -d isaacracing.net -d www.isaacracing.net` (this creates `/etc/letsencrypt/live/isaacracing.net`)

Later, to renew the certificate:

- `RENEW_DIR=/root/isaac-racing-server/letsencrypt && mkdir -p $RENEW_DIR && letsencrypt renew --webroot --webroot-path $RENEW_DIR && rm -rf $RENEW_DIR`

<br />

## Install as a service (optional)

- Install Supervisor:
  - `apt install supervisor`
  - `systemctl enable supervisor` (this is needed due to [a quirk in Ubuntu 16.04](http://unix.stackexchange.com/questions/281774/ubuntu-server-16-04-cannot-get-supervisor-to-start-automatically))
- Copy the configuration files:
  - `mkdir -p "/etc/supervisord"`
  - `cp "/root/isaac-racing-server/install/supervisord/supervisord.conf" "/etc/supervisord/supervisord.conf"`
  - `cp "/root/isaac-racing-server/install/supervisord/isaac-racing-server.conf" "/etc/supervisord/conf.d/isaac-racing-server.conf"`
- Start it: `systemctl start supervisor`

Later, to manage the service:

- Start it: `supervisorctl start isaac-racing-server`
- Stop it: `supervisorctl stop isaac-racing-server`
- Restart it: `supervisorctl restart isaac-racing-server`

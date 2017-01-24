#!/usr/bin/perl

# Imports
use strict;
use warnings;

# Configuration
my $directory = "/root/isaac-racing-server";
my $databaseName = "database.sqlite";
my $schemaName = "database_schema.sql";

# Remove the old database, if present
system "touch $directory/$databaseName";
system "rm -f $directory/$databaseName";

# Install the database
system "sqlite3 $directory/$databaseName < $directory/install/$schemaName";

# Rebuild Go dependencies
# (shouldn't be necessary)
#system "rm -rf \$GOPATH/pkg";
#system "go build -i";

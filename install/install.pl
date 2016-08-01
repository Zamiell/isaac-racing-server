#!/usr/bin/perl

# Imports
use strict;
use warnings;
use Cwd 'abs_path';

# Global variables
my $directory = abs_path($0);
if ($directory =~ /(.+)\/.+$/) {
	$directory = $1;
} else {
	die "Can't parse the script directory.\n";
}

# Install the database
system "sqlite3 $directory/../database.sqlite < $directory/database_schema.sql";
#system "sqlite3 $directory/../database.sqlite < $directory/seeds.sql";

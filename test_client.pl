#!/usr/bin/perl

# Imports
use strict;
use warnings;

# Variables
my $command;
my $output;
my $username;
my $password = "asdf";
my $accessToken;
my $cookie;
my $clientID = "tqY8tYlobY4hc16ph5B61dpMJ1YzDaAR";

# Validate command-line arguments
if (scalar @ARGV != 1) {
	die "Must provide test user number.\n";
}
if ($ARGV[0] == 1) {
	$username = "zamiel";
} elsif ($ARGV[0] == 2) {
	$username = "zamiel2";
} else {
	die "Invalid test user number.\n";
}

# Login (1/3)
$command = "curl https://isaacserver.auth0.com/oauth/ro --data 'grant_type=password&username=$username&password=$password&client_id=$clientID&connection=Isaac-Server-DB-Connection' --verbose 2>&1";
$output = `$command`;
if ($output =~ /{"access_token":"(.+)","token_type":"bearer"}/) {
	$accessToken = $1;
	print "access_token: $accessToken\n";
} else {
	die "\nFailed to parse the login response for step 1:\n\n$output\n";
}

# Login (2/3)
$command = "curl https://isaacitemtracker.com/login -H \"Content-Type: application/json\" --data \"{\\\"access_token\\\":\\\"$accessToken\\\",\\\"token_type\\\":\\\"bearer\\\"}\" --verbose 2>&1";
$output = `$command`;
if ($output =~ /Set-(Cookie: isaac.sid=.+; HttpOnly; Secure)/) {
	$cookie = $1;
	print "$cookie\n";
} else {
	die "Failed to parse the login response for step 2: $output\n";
}

# Login (3/3)
$command = "wscat --connect https://isaacitemtracker.com/ws --header \"$cookie\"";
system $command;

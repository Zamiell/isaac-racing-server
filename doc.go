/*

TODO
----

- move database erorr messages to parent

- ruleset: http://pastebin.com/mZYnxa0F
- race data commands: http://pastebin.com/PG8TpLqh

- userprofile
- leaderboard
- achievement logic


CLIENT TODO
-----------

- Test connecting to server when race status is "starting"



Misc. Notes
-----------

Get line count:
find . -name "*.go" | xargs cat | wc -l



Database
--------

Open the database:
sqlite3 database.sqlite

Get all of the users:
SELECT * FROM users;

Make people admins:
UPDATE users SET admin = 2 WHERE username = 'zamiel';
UPDATE users SET admin = 2 WHERE username = 'chronometrics';



Register
--------

Register:
curl https://isaacserver.auth0.com/dbconnections/signup -H "Content-Type: application/json" --data '{"client_id":"tqY8tYlobY4hc16ph5B61dpMJ1YzDaAR","email":"zamiel@zamiel.com","username":"zamiel","password":"1","connection":"Username-Password-Authentication"}' --verbose

Register response:
{"_id":"5770687c52fa77db5cea97ba","email_verified":false,"email":"zamiel@zamiel.com","username":"zamiel"}



Login
-----

Login (1/2):
curl https://isaacserver.auth0.com/oauth/ro --data "grant_type=password&username=zamiel&password=1&client_id=tqY8tYlobY4hc16ph5B61dpMJ1YzDaAR&connection=Username-Password-Authentication" --verbose

Login response (1/2):
{"access_token":"bWEPnEwPCLOBLbAL","token_type":"bearer"}

Login (2/2):
curl https://isaacitemtracker.com/login -H "Content-Type: application/json" --data "{\"access_token\":\"Nb7BG78AqgQHzAbk\",\"token_type\":\"bearer\"}" --verbose

Login response (2/2):
Set-Cookie: isaac.sid=MTQ2NzkxMjE1MHxEdi1CQkFFQ180SUFBUkFCRUFBQWVfLUNBQU1HYzNSeWFXNW5EQXNBQ1d4dloyZGxaRjlwYmdSaWIyOXNBZ0lBQVFaemRISnBibWNNQ2dBSVlYVjBhREJmYVdRR2MzUnlhVzVuREJvQUdEVTNOemN3Tm1Kak1qTTFNemMwWXprd05tVTFaakUwWXdaemRISnBibWNNQ2dBSWRYTmxjbTVoYldVR2MzUnlhVzVuREFnQUJucGhiV2xsYkE9PXwa-22yvtIQnLdqtlTmM0C8Y5czV5jbjMAnxdkyoY4eEw==; Path=/; Expires=Sat, 06 Aug 2016 17:22:30 GMT; Max-Age=2592000

Open a WebSocket connection with wscat:
COOKIE="isaac.sid=MTQ2ODEyMzIyN3xEdi1CQkFFQ180SUFBUkFCRUFBQVh2LUNBQU1HYzNSeWFXNW5EQXNBQ1d4dloyZGxaRjlwYmdSaWIyOXNBZ0lBQVFaemRISnBibWNNQ2dBSWRYTmxjbTVoYldVR2MzUnlhVzVuREFrQUIzcGhiV2xsYkRJR2MzUnlhVzVuREFjQUJXRmtiV2x1QTJsdWRBUUNBQUE9fAuYVPDRsKH6i90gVTKt3zYhK_h936q0FS6usbdO9GA7; Path=/; Domain=isaacitemtracker.com; Expires=Sun, 17 Jul 2016 04:00:27 GMT; Max-Age=604800; HttpOnly; Secure" && wscat --connect https://isaacitemtracker.com/ws --header "Cookie: $COOKIE"



Incoming WebSocket commands - chat
----------------------------------

Join a new chat room:
roomJoin {"room":"fartchannel"}

Leave a chat room:
roomLeave {"room":"fartchannel"}

Send a message to a chat room:
roomMessage {"room":"global","msg":"i poopd"}

Send a message to a chat room for a race:
roomMessage {"room":"_race_1","msg":"gg"}

Send a private message:
privateMessage {"name":"zamiel","msg":"private message lol"}

Get a list of all of the current rooms:
roomListAll {}



Outgoing WebSocket commands - chat
----------------------------------

When you join a new chat room, you get the list of people in it:
roomList {"room":"global","users":[{"name":"zamiel","admin":0,"squelched":0},{"name":"zamiel2","admin":0,"squelched":0}]}

When you join a new chat room, you get the chat history for the past 50 messages (or all the messages if its a "_race_#" channel):
roomHistory {"room":"global","history":[{"name":"zamiel","msg":"MrDestructoid","datetime":1469662590}]}

Someone else joined a chat room that you are in:
roomJoined {"room":"global","user":{"name":"zamiel2","admin":0,"squelched":0}}

Someone else left a chat room that you are in:
roomLeft {"room":"global","name":"chronometrics"}

Someone sent a message to a chat room that you are in:
roomMessage {"room":"global","name":"zamiel",msg":"i poopd"}

Someone sent you a private message:
privateMessage {"name":"chronometrics","msg":"i lit the candle"}

When a list of all the chat rooms is requested:
roomListAll [{"room":"global","numUsers":1}]

Someone got squelched:
roomSetSquelched {"room":"global","username":"cmondinger","squelched":1}

Someone got unsquelched:
roomSetSquelched {"room":"global","username":"cmondinger","squelched":0}

Someone got promoted:
roomSetAdmin {"room":"global","username":"sillypears","admin":1}

Someone got demoted:
roomSetAdmin {"room":"global","username":"sillypears","admin":0}



Incoming WebSocket commands - race
----------------------------------

Create a race:
raceCreate {}

Create a race with both optional arguments:
raceCreate {"name":"dee's race","ruleset":"diversity"}

Join a race:
raceJoin {"id":1}

Leave a race:
raceLeave {"id":1}

Ready in a race:
raceReady {"id":1}

Unready in a race:
raceUnready {"id":1}

Change a ruleset in a race (if you are the race captain):
raceRuleset {"id":3,"ruleset":{"type":"unseeded","character":4,"goal":"chest","seed":"-","instantStart":0}

Finish a race:
raceFinish {"id":1}

Quit a race:
raceQuit {"id":1}

Comment in a race:
raceComment {"id":1,"msg":"died to mom"}

Got a new item:
raceItem {"id":1,"itemID":100}

Got to a new floor:
raceFloor {"id":1,"floor":2}



Outgoing WebSocket commands - race
----------------------------------

When someone joins a new race or is already in an existing race upon connection:
racerList {"id":6,"racers":[{"name":"zamiel","status":"not ready","datetime_joined":1469178564,"datetime_finished":0,"place":0,"comment":"-","items":[],"floor":1}]}

When a new race is created:
raceCreated {"id":3,"name":"-","status":"open","ruleset":{"type":"unseeded","character":4,"goal":"chest","seed":"-","instantStart":0},"datetime_created":1469660053,"datetime_started":0,"captain":"zamiel","racers":["zamiel"]}

When someone joins a race:
raceJoin {"id":1,"name":"zamiel"}

When someone leaves a race:
raceLeft {"id":1,"name":"zamiel"}

When the race ruleset changes:
raceSetRuleset {"id":3,"ruleset":{"type":"unseeded","character":4,"goal":"chest","seed":"-","instantStart":0}}

When a race is starting (time is in UnixNano() format):
raceStart {"id":10,"time":1469147515988023769}

When someone readies up:
racerSetStatus {"id":1,"name":"zamiel","status":"ready"}

When someone unreadies:
racerSetStatus {"id":1,"name":"zamiel","status":"not ready"}

When someone finishes:
racerSetStatus {"id":1,"name":"zamiel","status":"finished"}

When someone quits:
racerSetStatus {"id":1,"name":"zamiel","status":"quit"}

When someone gets a new item:
racerAddItem {"id":1,"name":"zamiel","item":{"id":100,"floor":1}}

When someone gets to a new floor:
racerSetFloor {"id":1,"name":"zamiel","floor":2}



Incoming WebSocket commands - profile
-------------------------------------

Get the profile of a user:
profileGet {"name":"zamiel2"}

Set the username to a new stylization:
profileSetUsername {"name":"zAmIeL2"}



Outgoing WebSocket commands - profile
-------------------------------------

When a profile is requested:
profile { TODO }



Incoming WebSocket commands - admin
-----------------------------------

Ban a user:
adminBan {"name":"zamiel2"}

Unban a user:
adminUnban {"name":"zamiel2"}

Ban an IP:
adminBanIP {"ip":"1.2.3.4"}

Unban an IP:
adminUnbanIP {"ip":"1.2.3.4"}

Squelch a user:
adminSquelch {"name":"zamiel2"}

Unsquelch a user:
adminUnsquelch {"name":"zamiel2"}

Promote a user:
adminPromote {"name":"zamiel2"}

Demote a user:
adminDemote {"name":"zamiel2"}



Outgoing WebSocket commands - admin
-----------------------------------

You got banned:
error {"type":"adminBan","msg":"You have been banned. If you think this was a mistake, please contact the administration to appeal."}



Incoming WebSocket commands - miscellaneous
-------------------------------------------

Logout:
logout {}



Outgoing WebSocket commands - miscellaneous
-------------------------------------------

You did something right:
success {"type":"raceCreate","input":{"room":"","msg":"","name":"","ruleset":{"type":"","character":0,"goal":"","seed":"","instantStart":0},"id":0,"comment":"","itemID":0,"floor":0,"ip":""}}

An error occured:
error {"type":"logout","msg":"You have logged on from somewhere else, so I'll disconnect you here."}



*/

package main

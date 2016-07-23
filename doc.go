/*

TODO
----

- discover chat room list
- userprofile
- send last chat when joining race
- send last chat when joining channel
- achievements



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



WebSocket chat commands
-----------------------

Join chat channel:
roomJoin {"name":"fartchannel"}

Leave chat channel:
roomLeave {"name":"fartchannel"}

Send message:
roomMessage {"to":"global","msg":"i poopd"}

Send message to a race channel:
roomMessage {"to":"_race_1","msg":"gg"}

Send private message:
privateMessage {"to":"zamiel","msg":"private message lol"}

Get a list of all of the current rooms:
roomGetAll {}



WebSocket race commands
-----------------------

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

Change a ruleset in a race:
raceRuleset {"id":1,"ruleset":"diversity"}

Finish a race:
raceDone {"id":1}

Quit a race:
raceQuit {"id":1}

Comment in a race:
raceComment {"id":1,"msg":"died to mom"}

Got a new item:
raceItem {"id":1,"item_id":"100"}

Got to a new floor:
raceFloor {"id":1,"floor":2}



WebSocket profile commands
--------------------------

Get the profile of a user:
profileGet {"name":"zamiel2"}

Set the username to a new stylization:
profileSetUsername {"name":"zAmIeL2"}



WebSocket admin commands
------------------------

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



WebSocket miscellaneous commands
--------------------------------

Logout:
logout {}



WebSocket responses back from the server
----------------------------------------

An error occured:
error {"type":"logout","msg":"You have logged on from somewhere else, so I'll disconnect you here."}

You did something right:
success {"type":"raceCreate","msg":{"name":"poop2","ruleset":"diversity","id":1}}

Whenever a chat room is updated:
roomList {"room":"global","users":[{"name":"zamiel","admin":0,"squelched":0,"status":"","datetime_joined":0,"datetime_finished":0,"place":0,"comment":"","floor":0}]}

When a list of chat rooms is requested:
roomListAll [{"room":"global","numUsers":1}]

Whenever someone joins or leaves a race, a race changes status, or a race changes ruleset:
raceList [{"id":3,"name":"-","status":"open","ruleset":"unseeded","datetime_created":1469177311,"datetime_started":0,"captain":"zamiel","racers":["zamiel"]}]

Whenever someone does something inside of a race:
racerList {"id":6,"racers":[{"name":"zamiel","status":"not ready","datetime_joined":1469178564,"datetime_finished":0,"place":0,"comment":"-","items":[],"floor":1}]}

When a race is starting:
raceStart {"id":10,"time":1469147515988023769}

When a profile is requested:
profile { TODO }

*/

package main

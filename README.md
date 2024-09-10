# go-hellofresh/test
HelloFresh take home test in golang

To run service:
a CompileDaemon must be running.  To get install the daemon the following may be used
go install github.com/githubnemo/CompileDaemon

Once the CompileDaemon is installed,

working under the assumption that the app was built using:
cd test
go build -o jlacovTest

CompileDaemon -command=./jlacovTest

The daemon will kick off listening to port 9080 (per the .env file values)

POST
----

A POST to http://localhost:9080/event with a text body similar to

725982365063,0.2171728015,1647950356
1725982376083,0.5417924523,1536703272
1725982377072,0.0167227983,1424228047
1725982377080,0.6819181442,1531345087
1725982369069,0.8325055838,1201069418
1725982373074,0.0086158421,1474243545

Is expected to return a http status code of 202
If there are processing issues a different code will be presented


GET
---

A GET call to  http://localhost:9080/stats will return a string that has 5 parts.  
1) Total number of events found occurring in the last minute
2) A sum of the X values for all events found
3) An average of the X values for all the events found
4) A sum of the Y values for all events found
5) An average of the Y values for all the events found

If no events are found within the last minute, the following string is returned
"0,0.0000000000,0.0000000000,0,0"
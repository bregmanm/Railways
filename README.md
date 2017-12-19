# railways
Here is a simple JSON API implementation of stations and trains account system. HTTP server on start contains the empty list of stations. Each station contains the unique name and list of trains which stay on it at the moment. Each train contains the unique name. User can add and remove stations and add and remove trains to/from stations. In addition user can move train from one station to other. Here is the API description:

Description                   Path              Method   Body                       Example

Add station                   /station          POST     {"name":"<station name>"}  {"name":"London"}
Remove station                /station          DELETE   {"name":"<station name>"}  {"name":"London"}
Get list of stations          /stations         GET      No                         {"stations":["London", "Paris"]}
Add train to station          /{station}/train  POST     {"name":"<train name>"}    {"name":"Blue Express"} 
Remove train from station     /{station}/train  DELETE   {"name":"<train name>"}    {"name":"Blue Express"} 
Get list of trains on station /{station}/trains GET      No                         {"trains":["Blue Express", "Red Arrow"]}
Move train from one station to other
                              /trip             POST     {"FromStation":"<station name>",     {"FromStation":"London",
                                                          "ToStation":"<station name>"         "ToStation":"Paris",
                                                          "Train":"<train name>"}              "Train":"Blue Express"}

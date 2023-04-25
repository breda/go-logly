# Logly
Logly is a very simple logging server, you send logs to it and you can fetch them back.
It stores all data into a data file (data.db) and includes indexing options to speed up lookups.

# TODO:
This is a list of enhacenemnts that need to be worked on.

1. Add node roles: master, slave & forwarder.
    - Forwarder: forwards requests to masters & slaves.
    - Slave: Can run queries but is not responsible for the cluster.
    - Master: Same as slave but is responsible for cluster: can eject slaves is they misbehave.
2. Better log structure: time, level, message, origin, application.
3. Bulk fetch queries: by time period, by origin, by application, by range of IDs.
4. Index is currently only in memory and is rebuilt every time: find a way to serlaize the index into a file, with fast lookups still.
5. mTLS everywhere.
6. ?


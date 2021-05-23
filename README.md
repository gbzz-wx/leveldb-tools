# leveldb-tools
command line tools for operating lelveldb
-source :db's path
-target: db's path
-action: get ,put ,copy,delete
-key : key 
-value : value

get value case: start width 'person_'
leveldb-tools -action=get -source=/root/db -key=person_

copy value from a to b

leveldb-tools -action=copy -source=/root/db -key=person_  -target=/root/db2



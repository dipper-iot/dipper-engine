# Rule Engine Example Redis Queue


## Run Dipper Engine with redis queue
```shell
$ dipper-engine -q redis  --redis-host 127.0.0.1:6379 --iq --oq
Start Dipper Engine

Rules: 5
-----------------------------------------------------------
No Rule Name   Worker Status
1  switch      1      enable
2  log-core    1      enable
3  arithmetic  1      enable
4  fork        1      enable
5  conditional 1      enable
-----------------------------------------------------------

```


## Run Example 

```shell
$ go build
$ redis-queue.exe
Result To Queue Output: 
{
   "id": 430992873538715659,
   "chan_id": "test-1",
   "data": {
     "default": {
       "a": 10,
       "b": 20,
       "d": 5
     }
   },
   "result": {
     "3": {
       "session_id": 430992873538715659,
       "chan_id": "test-1",
       "id_node": "3",
       "from_engine": "log-core",
       "data": {
         "default": {
           "a": 10,
           "b": 20,
           "c": 30,
           "d": 5
         }
       },
       "branch_main": "",
       "next": [],
       "time": "2022-10-22T13:49:31.2850393+07:00",
       "type": "success",
       "error": null,
       "Debug": false
     },
     "4": {
       "session_id": 430992873538715659,
       "chan_id": "test-1",
       "id_node": "4",
       "from_engine": "conditional",
       "data": {
         "default": {
           "a": 10,
           "b": 20,
           "c": 30,
           "cond_a_b": 0,
           "d": 5
         }
       },
       "branch_main": "",
       "next": [],
       "time": "2022-10-22T13:49:31.2836959+07:00",
       "type": "success",
       "error": null,
       "Debug": true
     }
   }
 }
```


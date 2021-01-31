## Query plan
EXPLAIN format=json SELECT id, username FROM users WHERE first_name LIKE 'Bobby' and second_name LIKE 'Chase' ORDER BY id;

# no index testing
```json
{
  "query_block": {
    "select_id": 1,
    "cost_info": {
      "query_cost": "206040.60"
    },
    "ordering_operation": {
      "using_filesort": false,
      "table": {
        "table_name": "users",
        "access_type": "index",
        "key": "PRIMARY",
        "used_key_parts": [
          "id"
        ],
        "key_length": "8",
        "rows_examined_per_scan": 992553,
        "rows_produced_per_join": 12251,
        "filtered": "1.23",
        "cost_info": {
          "read_cost": "203590.34",
          "eval_cost": "2450.26",
          "prefix_cost": "206040.60",
          "data_read_per_join": "15M"
        },
        "used_columns": [
          "id",
          "username",
          "first_name",
          "second_name"
        ],
        "attached_condition": "((`db`.`users`.`first_name` like 'Bobby') and (`db`.`users`.`second_name` like 'Chase'))"
      }
    }
  }
}

```

## wrk -t1 -c1 -d10s  --latency 'http://0.0.0.0:3000/account/search_user?firstname=Bobby&secondname=Chase'

Running 10s test @ http://0.0.0.0:3000/account/search_user?firstname=Bobby&secondname=Chase <br>
1 threads and 1 connections <br>
Thread Stats   Avg      Stdev     Max   +/- Stdev  <br>
Latency   362.28ms   13.20ms 422.48ms   92.59%  <br>
Req/Sec     2.41      0.57     3.00     51.85%  <br>
Latency Distribution  <br>
50%  359.55ms  <br>
75%  362.74ms  <br> 
90%  367.26ms  <br>
99%  422.48ms  <br>
27 requests in 10.10s, 6.06KB read  <br>
Requests/sec:      2.67  <br>
Transfer/sec:     614.58B  <br>

## wrk -t1 -c10 -d10s  --latency 'http://0.0.0.0:3000/account/search_user?firstname=Bobby&secondname=Chase'
Running 10s test @ http://0.0.0.0:3000/account/search_user?firstname=Bobby&secondname=Chase  <br>
1 threads and 10 connections  <br>
Thread Stats   Avg      Stdev     Max   +/- Stdev  <br>
Latency   627.75ms  145.47ms 942.75ms   62.82%  <br>
Req/Sec    17.56     12.17    60.00     72.22%  <br>
Latency Distribution  <br>
50%  648.16ms  <br>
75%  747.85ms  <br>
90%  807.96ms  <br>
99%  874.06ms  <br>
156 requests in 10.10s, 35.04KB read  <br>
Requests/sec:     15.45  <br>
Transfer/sec:      3.47KB  <br>

## wrk -t1 -c100 -d10s  --latency 'http://0.0.0.0:3000/account/search_user?firstname=Bobby&secondname=Chase'
Running 10s test @ http://0.0.0.0:3000/account/search_user?firstname=Bobby&secondname=Chase  <br>
1 threads and 100 connections  <br>
Thread Stats   Avg      Stdev     Max   +/- Stdev  <br>
Latency     0.00us    0.00us   0.00us     nan%  <br>
Req/Sec    74.23     70.74   262.00     84.62%  <br>
Latency Distribution  <br>
50%    0.00us  <br>
75%    0.00us  <br>
90%    0.00us  <br>
99%    0.00us  <br>
100 requests in 10.01s, 22.46KB read  <br>
Socket errors: connect 0, read 0, write 0, timeout 100  <br>
Requests/sec:      9.99  <br> 
Transfer/sec:      2.24KB  <br>

## wrk -t1 -c1000 -d10s  --latency 'http://0.0.0.0:3000/account/search_user?firstname=Bobby&secondname=Chase'
Running 10s test @ http://0.0.0.0:3000/account/search_user?firstname=Bobby&secondname=Chase  <br>
1 threads and 1000 connections  <br>
Thread Stats   Avg      Stdev     Max   +/- Stdev  <br>
Latency     1.53s     0.00us   1.53s   100.00%  <br>
Req/Sec     0.33      0.58     1.00     66.67%  <br>
Latency Distribution  <br>
50%    1.53s  <br>
75%    1.53s  <br>
90%    1.53s  <br>
99%    1.53s  <br>
3 requests in 10.10s, 690.00B read  <br>
Socket errors: connect 748, read 3203, write 0, timeout 2  <br> 
Requests/sec:      0.30  <br>
Transfer/sec:      68.34B  <br>

## index on first_name + index on second_name
create index users_first_name_index
on users (first_name);
create index users_first_name_index
on users (second_name);
```json
{
  "query_block": {
    "select_id": 1,
    "cost_info": {
      "query_cost": "16.35"
    },
    "ordering_operation": {
      "using_filesort": true,
      "cost_info": {
        "sort_cost": "1.00"
      },
      "table": {
        "table_name": "users",
        "access_type": "index_merge",
        "possible_keys": [
          "users_first_name_index",
          "users_second_name_index"
        ],
        "key": "intersect(users_first_name_index,users_second_name_index)",
        "key_length": "52,52",
        "rows_examined_per_scan": 1,
        "rows_produced_per_join": 1,
        "filtered": "100.00",
        "cost_info": {
          "read_cost": "15.15",
          "eval_cost": "0.20",
          "prefix_cost": "15.35",
          "data_read_per_join": "1K"
        },
        "used_columns": [
          "id",
          "username",
          "first_name",
          "second_name"
        ],
        "attached_condition": "((`db`.`users`.`first_name` like 'Bobby') and (`db`.`users`.`second_name` like 'Chase'))"
      }
    }
  }
}
```

## Running 10s test @ http://0.0.0.0:3000/account/search_user?firstname=Bobby&secondname=Chase
1 threads and 1 connections  <br>
Thread Stats   Avg      Stdev     Max   +/- Stdev  <br>
Latency     3.66ms    1.52ms  14.85ms   87.65%  <br>
Req/Sec   279.21     77.58   390.00     65.00%  <br>
Latency Distribution  <br>
50%    3.11ms  <br>
75%    3.91ms  <br>
90%    5.69ms  <br>
99%    9.57ms  <br>
2783 requests in 10.01s, 625.09KB read  <br>
Requests/sec:    277.88  <br>
Transfer/sec:     62.42KB  <br>

## Running 10s test @ http://0.0.0.0:3000/account/search_user?firstname=Bobby&secondname=Chase
1 threads and 10 connections  <br>
Thread Stats   Avg      Stdev     Max   +/- Stdev  <br>
Latency     9.73ms    4.46ms  63.22ms   79.88%  <br>
Req/Sec     1.05k   136.49     1.30k    65.00%  <br>
Latency Distribution  <br>
50%    8.66ms  <br>
75%   11.55ms  <br>
90%   15.11ms  <br>
99%   25.14ms  <br>
10468 requests in 10.00s, 2.30MB read  <br>
Requests/sec:   1046.37  <br>
Transfer/sec:    235.02KB  <br>

## Running 10s test @ http://0.0.0.0:3000/account/search_user?firstname=Bobby&secondname=Chase
1 threads and 100 connections  <br>
Thread Stats   Avg      Stdev     Max   +/- Stdev  <br>
Latency   106.75ms   68.96ms 566.87ms   71.19%  <br>
Req/Sec     0.98k   151.06     1.28k    76.00%  <br>
Latency Distribution  <br>
50%   95.72ms  <br>
75%  142.66ms  <br>
90%  197.48ms  <br>
99%  323.41ms  <br>
9739 requests in 10.02s, 2.14MB read  <br>
Requests/sec:    972.20  <br>
Transfer/sec:    218.37KB  <br>

## Running 10s test @ http://0.0.0.0:3000/account/search_user?firstname=Bobby&secondname=Chase
1 threads and 1000 connections  <br>
Thread Stats   Avg      Stdev     Max   +/- Stdev  <br>
Latency   501.41ms  447.83ms   1.85s    76.02%  <br>
Req/Sec   337.46    299.30     1.20k    75.00%  <br>
Latency Distribution  <br>
50%  359.51ms  <br>
75%  657.92ms  <br>
90%    1.31s  <br>
99%    1.78s  <br>
2883 requests in 10.04s, 647.55KB read  <br>
Socket errors: connect 0, read 2888, write 9, timeout 77  <br>
Requests/sec:    287.07  <br>
Transfer/sec:     64.48KB  <br>

## index on (first_name, second_name)
create index users_first_name_second_name_index
on users (first_name, second_name);
```json
{
  "query_block": {
    "select_id": 1,
    "cost_info": {
      "query_cost": "3.41"
    },
    "ordering_operation": {
      "using_filesort": true,
      "cost_info": {
        "sort_cost": "1.00"
      },
      "table": {
        "table_name": "users",
        "access_type": "range",
        "possible_keys": [
          "users_first_name_second_name_index"
        ],
        "key": "users_first_name_second_name_index",
        "used_key_parts": [
          "first_name",
          "second_name"
        ],
        "key_length": "104",
        "rows_examined_per_scan": 1,
        "rows_produced_per_join": 1,
        "filtered": "100.00",
        "index_condition": "((`db`.`users`.`first_name` like 'Bobby') and (`db`.`users`.`second_name` like 'Chase'))",
        "cost_info": {
          "read_cost": "2.21",
          "eval_cost": "0.20",
          "prefix_cost": "2.41",
          "data_read_per_join": "1K"
        },
        "used_columns": [
          "id",
          "username",
          "first_name",
          "second_name"
        ]
      }
    }
  }
}

```
## Running 10s test @ http://0.0.0.0:3000/account/search_user?firstname=Bobby&secondname=Chase
1 threads and 1 connections  <br>
Thread Stats   Avg      Stdev     Max   +/- Stdev  <br>
Latency     2.12ms  464.94us   9.25ms   91.71%  <br>
Req/Sec   476.01     37.98   575.00     72.00%  <br>
Latency Distribution  <br>
50%    2.05ms  <br>
75%    2.20ms  <br>
90%    2.41ms  <br>
99%    4.22ms  <br>
4741 requests in 10.01s, 1.04MB read  <br>
Requests/sec:    473.81  <br>
Transfer/sec:    106.42KB  <br>

## Running 10s test @ http://0.0.0.0:3000/account/search_user?firstname=Bobby&secondname=Chase
1 threads and 10 connections  <br>
Thread Stats   Avg      Stdev     Max   +/- Stdev  <br>
Latency     7.87ms    4.64ms  41.13ms   86.57% <br>
Req/Sec     1.34k   382.54     1.96k    61.00%  <br>
Latency Distribution <br>
50%    6.20ms <br>
75%    9.62ms <br>
90%   13.51ms <br>
99%   25.81ms <br>
13321 requests in 10.00s, 2.92MB read  <br>
Requests/sec:   1331.59  <br>
Transfer/sec:    299.09KB  <br>

## Running 10s test @ http://0.0.0.0:3000/account/search_user?firstname=Bobby&secondname=Chase
1 threads and 100 connections <br>
Thread Stats   Avg      Stdev     Max   +/- Stdev <br>
Latency   113.49ms   64.18ms 482.12ms   66.56% <br>
Req/Sec     0.90k   327.57     1.53k    59.60% <br>
Latency Distribution <br>
50%  113.29ms <br>
75%  147.43ms <br>
90%  196.75ms <br>
99%  293.12ms <br>
8866 requests in 10.01s, 1.94MB read <br>
Requests/sec:    885.63 <br>
Transfer/sec:    198.92KB <br>

## Running 10s test @ http://0.0.0.0:3000/account/search_user?firstname=Bobby&secondname=Chase
1 threads and 1000 connections <br>
Thread Stats   Avg      Stdev     Max   +/- Stdev <br>
Latency   170.15ms   99.27ms   1.75s    77.34% <br>
Req/Sec   631.24    218.49     1.14k    77.63% <br>
Latency Distribution <br>
50%  207.24ms <br>
75%  228.04ms <br>
90%  246.00ms <br>
99%  352.32ms <br>
4888 requests in 10.02s, 1.07MB read  <br>
Socket errors: connect 0, read 26, write 0, timeout 96  <br>
Requests/sec:    487.61  <br>
Transfer/sec:    109.52KB  <br>

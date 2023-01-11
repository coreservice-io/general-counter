### DailyGCounterModel

```
//daily_g_counter_agg
//id is a string, format:
[date]:[gkey]:[gtype]
//gkey can be anyting like your userid, accountid ,etc
//gtype can be anything like 'user_credit','account_credit',etc
```

```
PUT /{project_name}_g_counter_daily_agg
{
    "mappings": {
        "dynamic": "false",
        "properties": {
            "sql_id": {
                "type": "long"
            },
            "id": {
                "type": "keyword",
                "ignore_above": 256
            },
            "gkey": {
                "type": "keyword",
                "ignore_above": 256
            },
            "gtype": {
                "type": "keyword",
                "ignore_above": 256
            },
            "date": {
                "type": "date",
                "format": "yyyy-MM-dd"
            },
            "amount": {
               "type": "long"
            }
        }
    }
}
```

//g_counter_detail
//id same as sql_id
```
PUT /{project_name}_g_counter_detail
{
    "mappings": {
        "dynamic": "false",
        "properties": {
            "sql_id": {
                "type": "long"
            },
            "id":{
                "type": "keyword",
                "ignore_above": 256
            },
            "gkey": {
                "type": "keyword",
                "ignore_above": 256
            },
            "gtype": {
                "type": "keyword",
                "ignore_above": 256
            },
            "datetime": {
                "type": "date",
                "format": "yyyy-MM-dd HH:mm:ss"
            },
            "amount": {
               "type": "long"
            },
            "msg": {
                "type": "text",
                "index":false
            }
        }
    }
}
```


# WSTunnel 

## INTRO

simple websocket to tcp proxy program let you safely expose internal services to your enviroment or deal with 

some reachable problem

## INSTALL

simple download and compile

```

$ git clone https://github.com/griffinexe/ws2TCP.git

$ cd ws2TCP

$ go build

```

if commands run with no errors, you will get WSTunnel binary executable file

## CONFIG

The program can running in 2 roles: server and client, controlled by config.json file

client config.json example:

```

{
    "client":{
        "upstream":"ws(s)://your.upstream.site/path_in_server_config",
        "listenmap":{
            "service1":["127.0.0.1:9000", "tcp"],
            "service2":["127.0.0.1:9001", "tcp"],
            "service3":["127.0.0.1:9002", "tcp"]
        },
        "acl":{
            "service1":"any_secret_strings_match_entry_in_server_config",
            "service2":"to_prevent_unautnorized_access",
            "service3":"if_not_match_connection_will_close"
        },
        "proxy":{
            "enabled":false,
            "url":"proxy_address"
        }
    }
}

```

server config.json example:

```

{
    "server":{
        "listen":"127.0.0.1:8088",
        "path":"/ws",
        "servicemap":{
            "service1":["127.0.0.1:6667", "tcp"],
            "service2":["127.0.0.1:3000", "tcp"],
            "service3":["127.0.0.1:8085", "tcp"]
        },
       "acl":{
            "service1":"any_secret_strings",
            "service2":"to_prevent_unautnorized_access",
            "service3":"if_not_match_connection_will_close"
        },
        "tls":{
            "enabled":false,
            "keyfile":"",
            "certfile":""
        }
    }
}

```

## START

simple start and leave it in background, config.json must exist

```
$ WSTunnel
```
#Loguru
Its really simple to use and bind context to the logger but, it does not have a way to not duplicate the logs in the console. 
Besides that, the time to render the logs is really out of sync with the other libraries.


output 
```
2023-02-13 09:06:51.447 | ERROR    | __main__:<module>:16 - ERROR
{"text": "2023-02-13 09:06:51.447 | ERROR    | __main__:<module>:16 - ERROR\n", "record": {"elapsed": {"repr": "0:00:00.237521", "seconds": 0.237521}, "exception": null, "extra": {"ip": "192.168.0.1", "user": "someone"}, "file": {"name": "loguru_lib.py", "path": "/root/playground/python/log_lib/loguru_lib.py"}, "function": "<module>", "level": {"icon": "❌", "name": "ERROR", "no": 40}, "line": 16, "message": "ERROR", "module": "loguru_lib", "name": "__main__", "process": {"id": 2567, "name": "MainProcess"}, "thread": {"id": 140496342511616, "name": "MainThread"}, "time": {"repr": "2023-02-13 09:06:51.447363-03:00", "timestamp": 1676290011.447363}}}
2023-02-13 09:06:51.447 | INFO     | __main__:<module>:15 - Contextualize your logger easily
{"text": "2023-02-13 09:06:51.447 | INFO     | __main__:<module>:15 - Contextualize your logger easily\n", "record": {"elapsed": {"repr": "0:00:00.237765", "seconds": 0.237765}, "exception": null, "extra": {"ip": "192.168.0.1", "user": "someone"}, "file": {"name": "loguru_lib.py", "path": "/root/playground/python/log_lib/loguru_lib.py"}, "function": "<module>", "level": {"icon": "ℹ️", "name": "INFO", "no": 20}, "line": 15, "message": "Contextualize your logger easily", "module": "loguru_lib", "name": "__main__", "process": {"id": 2567, "name": "MainProcess"}, "thread": {"id": 140496342511616, "name": "MainThread"}, "time": {"repr": "2023-02-13 09:06:51.447607-03:00", "timestamp": 1676290011.447607}}}
2023-02-13 09:06:51.447 | ERROR    | __main__:<module>:16 - ERROR
{"text": "2023-02-13 09:06:51.447 | ERROR    | __main__:<module>:16 - ERROR\n", "record": {"elapsed": {"repr": "0:00:00.238085", "seconds": 0.238085}, "exception": null, "extra": {"ip": "192.168.0.1", "user": "someone"}, "file": {"name": "loguru_lib.py", "path": "/root/playground/python/log_lib/loguru_lib.py"}, "function": "<module>", "level": {"icon": "❌", "name": "ERROR", "no": 40}, "line": 16, "message": "ERROR", "module": "loguru_lib", "name": "__main__", "process": {"id": 2567, "name": "MainProcess"}, "thread": {"id": 140496342511616, "name": "MainThread"}, "time": {"repr": "2023-02-13 09:06:51.447927-03:00", "timestamp": 1676290011.447927}}}
```

# Structlog

```
{"message_id": "1234", "event": "This is a test message", "level": "info"}
```

# Logging standard library
Its not structured and it does not have a way to bind context to the logger.
Using that will be hard to filter the logs in the future.

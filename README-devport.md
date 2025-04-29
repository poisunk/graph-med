### dev port



#### service port

| service name | api service port(1xxx) | rpc service port(2xxx) | other service port(3xxx) |
|--------------| ---------------------- | ---------------------- |--------------------------|
| usercenter   | 1001                   | 2001                   |                          |
| captcha      | 1002                   | 2002                   |                          |
| chat         | 1003                   | 2003                   |                          |
| mqueue       | -                      | -                      |                          |



#### Prometheus Port

⚠️Online containers are separate, so online all set to the same port on it, local because in a container development, to prevent port conflicts

| service name     | prometheus port |
|------------------| --------------- |
| usercenter-api   | 4001            |
| usercenter-rpc   | 4002            |
| captcha-api      | 4003            |
| captcha-rpc      | 4004            |
| chat-api         | 4005            |
| chat-rpc         | 4006            |

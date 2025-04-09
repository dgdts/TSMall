# *** Project

## introduce

- Use the [Hertz](https://github.com/cloudwego/hertz/) framework
- Integration of pprof, cors, recovery, access_log, gzip and other extensions of Hertz.
- Generating the base code for unit tests.
- Provides basic profile functions.
- Provides the most basic MVC code hierarchy.

## Directory structure

|  catalog   | introduce  |
|  ----  | ----  |
| conf  | Configuration files |
| cmd/main.go  | Startup file |
| hertz_gen  | Hertz generated model |
| biz/handler  | Used for request processing, validation and return of response. |
| biz/service  | The actual business logic. |
| biz/router  | Routing and middleware registration |
| biz/utils  | Wrapped some common methods |

## How to run

```shell
./build_hz.sh
./script/run.sh
```
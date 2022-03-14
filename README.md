# wisdom-server

## How does it work
1) Server sends random number >= 2 to client*
2) Client should find prime factors of this number and send them back to server
3) Server checks solution and sends random "wise" quote

\* As more connections from IP come, as more complex the number is. Number's complexity is defined by it's byte-length. Complexity recovers over time.

[Client's code.](https://github.com/don2quixote/wisdom-client)

## ENV Configuration:
| name                        | type   | description                                                |
| --------------------------- | ------ | ---------------------------------------------------------- |
| PORT                        | int    | Port to laynch TCP server                                  |
| COMPLEXITY_FACTOR           | float  | Shows how fast does complexity grow                        |
| MAX_COMPLEXITY              | int    | Limits the maximum complexity (length of challange number) |
| COMPLEXITY_DURATION_SECONDS | int    | Time to recover complextity level                          |

## How to launch
Using Makefile (inconvenient change env configuration due to rebuilds):
```
make run
```

Or using docker manually:d
```
docker build -f build/Dockerfile -t wisdom-server .
docker run --name wisdom-server --rm \
    -e PORT=4444 \
    -e COMPLEXITY_FACTOR=0.25 \
    -e MAX_COMPLEXITY=16 \
    -e COMPLEXITY_DURATION_SECONDS=44 \
    -p 4444:4444 wisdom-server
```
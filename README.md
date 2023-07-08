# anagram-finder-trendyol

# Assumptions
- all words are in lowercase
- all words are in english
- there cannot be more than 255 identical characters in a word. (uint8 limit)

# Setup
```
go install github.com/golang/mock/mockgen@v1.6.0
curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.53.3
```

# How to run
After building the project, you can run the binary with the following command
```
Usage of ./anagram-finder-trendyol:
  -file string
        Path to the text file to be processed. It is required if url is not given
  -log-level uint
        Log level (0-6). Default is 4 (info) (default 4)
  -redis-db int
        Redis db for storage. The given db will be flushed before the application starts
  -redis-host string
        Redis host
  -redis-password string
        Redis password
  -redis-port int
        Redis port
  -storage-type string
        Storage type (local or redis). Default is local (default "local")
  -url string
        URL to the text file to be processed. It is required if file is not given
  -words-channel-size int
        Words channel size. Default is 8 (default 8)
  -worker-pool-size int
        Worker pool size. Default is 16 (default 16)
```
example:
```./anagram-finder-trendyol --url https://raw.githubusercontent.com/Trendyol/assignments/main/anagrams.txt --storage-type redis --redis-host localhost --redis-port 6379 --log-level 4 ```
# How to run with docker
``` 
    if you want to use redis as storage, you should run redis first
docker run --name redis -d redis

docker build -t anagram-trendyol .
docker run anagram-trendyol --url https://raw.githubusercontent.com/Trendyol/assignments/main/anagrams.txt --storage-type redis --redis-host 172.17.0.2 --redis-port 6379 --log-level 4
```

# How to run tests
```
./scripts/setup_mocks.sh
go test ./... -cover
```
# Benchmark
```
With Redis (200.000 words) 
    hash and set operations - 6.385868885s
    hash and set operations with concurrency (8 goroutines)  - 955.66091ms
    hash and set operations with concurrency (16 goroutines) - 647.531846ms
    hash and set operations with concurrency (32 goroutines) - 557.25159ms
    read operations with pipe - 1.312478984s
With Local (200.000 words) hashmap with mutex
    hash and set operations - 63.747096ms
    hash and set operations with concurrency - 61.895841ms
    read operations with pipe - 645.10004ms

```

# K8s
```
I have not tested it but I think it makes sense to run the application in a k8s cluster.
- use env maps for configuration
- pass env maps to the container
```

# Notes
.golangci.yaml from https://gist.github.com/maratori/47a4d00457a92aa426dbd48a18776322#file-golangci-yml
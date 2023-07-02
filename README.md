# anagram-finder-trendyol

# Assumptions
- all words are in lowercase
- all words are in english
- there cannot be more than 255 identical characters in a word. (uint8 limit)

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



.golangci.yaml from https://gist.github.com/maratori/47a4d00457a92aa426dbd48a18776322#file-golangci-yml
## zcache

#### quick start
```go
package main

import (
	"fmt"
	"zcache"
)

func main() {
    cache := zcache.NewCache()
    v := 123
    cache.Set("a", v)
    if v != cache.Get("a") {
        fmt.Printf("Get invalid value: %+v", v)
    }
    cache.Del("a")
    if cache.Get("a") != nil {
		fmt.Println("Del failed")
    }
}
```

#### Multiple servers
server list:
- localhost:5205
- localhost:5206
- localhost:5207

On server1 start one 
```go
c1 = NewDCache(
    "localhost:5205",
    []string{"localhost:5206", "localhost:5207"},
    WithLru(),
    WithSize(1000),
)
```
On server2 start other one
```go
c2 = NewDCache(
    "localhost:5206",
    []string{"localhost:5205", "localhost:5207"},
    WithLru(),
    WithSize(1000),
)
```
On server3 start last
```go
c3 = NewDCache(
    "localhost:5207",
    []string{"localhost:5205", "localhost:5206"},
    WithLru(),
    WithSize(1000),
)
```

On server1, set a value:
```go
c1.Set("a", 123)
```

On another server, this key will be none:
```go
c2.Get("a") == nil
```

Now, sever2 has newest value, you should delete it before set:
```go
c2.Del("a")         // sync to server1 and server3
c2.Set("a", 456)
```

On server1:
```go
c1.Get("a") == nil
```


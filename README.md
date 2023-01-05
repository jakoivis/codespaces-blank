# utcoffset

Converts UTC offset formats ±[hh]:[mm], ±[hh][mm], ±[hh] e.g. +02:34, +0234, +02 to seconds. 
Allows using case insensitive UTC prefix e.g. UTC+02:34, UTC+0234, UTC+02

## Example usage

```go
package main

import (
	"fmt"
	"time"

	"github.com/jakoivis/utcoffset"
)

func main() {
	offset := "+0200"

	if utcoffset.IsValidUtcOffset(offset) {
		fmt.Printf("%v is valid\n", offset)
	} else {
		fmt.Printf("%v is not valid\n", offset)
	}

	seconds, err := utcoffset.UtcOffsetSeconds(offset)

	if err != nil {
		fmt.Println(err)
	}

	fmt.Printf("UTC offset %v is in seconds %v\n", offset, seconds)

	location := time.FixedZone("myzone", seconds)

	fmt.Println(time.Now().In(location))
}
```

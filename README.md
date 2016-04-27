# wit.ai-go

Go library for [wit.ai](https://wit.ai/) HTTP API

## How to install

```bash
$ go get -u github.com/meinside/wit.ai-go
```

## Usage

```go
package main

import (
	"fmt"

	witai "github.com/meinside/wit.ai-go"
)

const (
	Token = "YOUR-WIT.AI-APP-SPECIFIC-TOKEN-HERE"
)

func main() {
	token := Token
	c := witai.NewClient(&token)

	//c.Verbose = false
	c.Verbose = true // for verbose messages

	// message
	if result, err := c.QueryMessage("how's the weather today?", nil, "", "", 1); err == nil {
		fmt.Printf("query message result: %+v\n", result)
	} else {
		fmt.Printf("%s\n", err)
	}

	// speech
	if result, err := c.QuerySpeechMp3("/some/where/test_voice.mp3", nil, "", "", 1); err == nil {
		fmt.Printf("query speech result: %+v\n", result)
	} else {
		fmt.Printf("%s\n", err)
	}
}
```

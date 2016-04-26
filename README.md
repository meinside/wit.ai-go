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

	// message
	if result, err := c.Message("how's the weather today?", nil, "", "", 1); err == nil {
		fmt.Printf("message result = %+v\n", result)
	} else {
		fmt.Printf("message error = %s\n", err)
	}

	// speech
	if result, err := c.SpeechMp3("/some/path/test_voice.mp3", nil, "", "", 1); err == nil {
		fmt.Printf("speech result = %+v\n", result)
	} else {
		fmt.Printf("speech error = %s\n", err)
	}
}
```

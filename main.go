package main

import (
	"fmt"
	"log"

	hook "github.com/robotn/gohook"
)

const exampleJSON = "./example.json"

func main() {
  
  ps, err := loadPhraseSet(exampleJSON)
  if err != nil {
    log.Fatal(err)
  }
  
  chw := parseWheel(ps)
  setupAppListener(chw)
}


func setupAppListener(chw ChatWheelI) {
  fmt.Println("--- Press q to enter into active mode ---")

  hook.Register(hook.KeyDown, []string{}, func(e hook.Event) {
    chw.ReactOnKey(e)
  })

	s := hook.Start()
	<-hook.Process(s)
}


package main

import (
	"fmt"
	"log"

	hook "github.com/robotn/gohook"
)

func main() {
  
  chosenFile := getCurrWheelFile()
  ps, err := loadPhraseSet(chosenFile)
  if err != nil {
    log.Fatal(err)
  }
  
  chw := parseWheel(ps)
  setupAppListener(chw)
}


func setupAppListener(chw ChatWheelI) {
  msg := fmt.Sprintf(`--- [biblade] ---
1) Press "%v" to enter into active mode
2) To choose a phrase simply press a corresponding key that
is specified in square brackets;
3) To exit from active mode without confirmation please
press any other key`, string(getActivatingChar()))
  fmt.Println(msg)

  hook.Register(hook.KeyDown, []string{}, func(e hook.Event) {
    chw.ReactOnKey(e)
  })

	s := hook.Start()
	<-hook.Process(s)
}


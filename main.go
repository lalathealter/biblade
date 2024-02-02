package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/atotto/clipboard"
	"github.com/go-vgo/robotgo"
	hook "github.com/robotn/gohook"
)

const exampleJSON = "./example.json"

func main() {
  
  ps, err := loadPhraseSet(exampleJSON)
  if err != nil {
    log.Fatal(err)
  }

  
  chw := produceChatWheelFrom(ps)
  setupKeyboardListener(chw)
}

func readClipboard() string {
  res, err := clipboard.ReadAll()
  if err != nil {
    log.Fatal(err)
  }

  return res
}

type PhraseSet map[string]string

type ReactOptionI interface {
  GetTag() string
  Response()
}

type ReactMap map[rune]ReactOptionI
func (rmap ReactMap) ReactOnKey(ev hook.Event) error {
  kchar := ev.Keychar

  f, ok := rmap[kchar] 
  if !ok {
    return ErrNoReactionMethod
  }
  f.Response()
  return nil
}

func (rmap ReactMap) IntroduceChatOptions() {
  for _, reopt := range rmap {
    robotgo.TypeStr(reopt.GetTag())
  }
}

var ErrNoReactionMethod = errors.New("couldn't execute command;")

func produceChatWheelFrom(ps PhraseSet) ChatWheelI {
  i := 1
  rmap := ReactMap{}
  fullLen := 0
  for tag, phrase := range ps {
    id := strconv.Itoa(i)
    opt := formIntroOptionText(tag, id)

    rmap[rune('0'+i)] = formReactOption(opt, phrase, &fullLen)
    i++
    fullLen += len(opt)
  }
  return rmap
}


type ReactOption struct {
  Tag string
  Content string
  FrameSize *int
}

func (reopt ReactOption) GetTag() string {
  return reopt.Tag
}

func (reopt ReactOption) Response() {
  removePreviousCharacters(reopt.FrameSize)
  robotgo.TypeStr(reopt.Content)
}

func formReactOption(optTag string, v string, fullLen *int) ReactOptionI {
  reopt := ReactOption{optTag, v, fullLen}
  return reopt
}

func removePreviousCharacters(n *int) {
  robotgo.KeyDown(robotgo.Shift)
  for i := 0; i < *n + 1; i++ {
    robotgo.KeyPress(robotgo.Left)
  }
  robotgo.KeyUp(robotgo.Shift)
  robotgo.KeyPress(robotgo.Backspace)
}

const MAX_INTRO_LEN = 12
func formIntroOptionText(key, id string) string {
  if len(key) > MAX_INTRO_LEN {
    key = key[:MAX_INTRO_LEN-3] + "..."
  }
  return fmt.Sprintf("[%v] %v ", id, key)
}

func loadPhraseSet(path string) (PhraseSet, error) {
  f, err := os.ReadFile(path)
  if err != nil {
    return nil, err
  }
  
  phs := PhraseSet{}
  err = json.Unmarshal(f, &phs) 
  return phs, err
}

type ChatWheelI interface {
  ReactOnKey(e hook.Event) error
  IntroduceChatOptions()
}




func setupKeyboardListener(chw ChatWheelI) {
  fmt.Println("--- Press q to enter into active mode ---")
  activeMode := false
  one := 1

  hook.Register(hook.KeyDown, []string{}, func(e hook.Event) {
    if activeMode {
      chw.ReactOnKey(e)
      activeMode = false
    }
  })
	hook.Register(hook.KeyDown, []string{"q"}, func(e hook.Event) {
    activeMode = true
    removePreviousCharacters(&one)
    chw.IntroduceChatOptions()
	})

	s := hook.Start()
	<-hook.Process(s)
}



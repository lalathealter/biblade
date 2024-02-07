package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"

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

  
  chw := parseWheel(ps)
  setupKeyboardListener(chw)
}

func readClipboard() string {
  res, err := clipboard.ReadAll()
  if err != nil {
    log.Fatal(err)
  }

  return res
}


type PhraseSet [][2]any // string | PhraseSet

type ReactOptionI interface {
  GetTag() string
  Response()
}


type ErrNoReactionForKey struct {
  Key rune
}
func (keyerr ErrNoReactionForKey) Error() string {
  return fmt.Sprint("Ignoring key ", keyerr.Key)
}


func deductIndFrom(keyChar rune) int {
  return int(keyChar - '0') - 1
}

func makeKey(n int) rune {
  return rune('0'+n + 1)
}



type WheelItemI interface {
  Response() *WheelFrame
  GetTag() string
}


type WheelFrame struct {
  Key rune
  Prompt string
  Items []WheelItemI
  FrameSize int
}

func (wf *WheelFrame) Response() *WheelFrame {
  for _, wi := range wf.Items {
    robotgo.TypeStr(wi.GetTag())
  }
  return wf
}

func (wf *WheelFrame) GetTag() string {
  return makeTag(wf.Key, wf.Prompt)
}

func (wf *WheelFrame) addItem(whi WheelItemI) {
  wf.Items = append(wf.Items, whi)
  wf.FrameSize += len(whi.GetTag())
}

type WheelChatOption struct {
  Key rune
  Prompt string
  Text string
}

func (wco WheelChatOption) GetTag() string {
  return makeTag(wco.Key, wco.Prompt)
}

func (wco WheelChatOption) Response() *WheelFrame {
  robotgo.TypeStr(wco.Text)
  return nil
}


const MAX_INTRO_LEN = 10
func makeWheelChatOption(key rune, prompt string, phrase string) WheelChatOption {
  if len(prompt) > MAX_INTRO_LEN {
    prompt = prompt[:MAX_INTRO_LEN-2] + ".."
  }
  return WheelChatOption{key, prompt, phrase}
}

type WheelController struct {
  Start *WheelFrame
  Current *WheelFrame
}

func (wc *WheelController) getCurrOpts() []WheelItemI {
  return wc.Current.Items
}

func (wc *WheelController) ReactOnKey(ev hook.Event) error {
  kchar := ev.Keychar
  if isBackspace(kchar) {
    return nil
  }

  if wc.Current == nil {
    if kchar == ACTIVATING_CHAR {
      wc.Current = wc.Start.Response()
    }
    return nil
  }

  i := deductIndFrom(kchar)
  removePreviousCharacters(wc.Current.FrameSize)
  opts := wc.getCurrOpts()
  if isOutOfBounds(i, opts) {
    wc.Current = nil
    return ErrNoReactionForKey{kchar}
  }
  whi := opts[i].Response()
  wc.Current = whi
  return nil
}

func makeTag(key rune, prompt string) string {
  return fmt.Sprintf("[%v] %v ", string(key), prompt)
}

func makeWheelFrame(key rune, prompt string) *WheelFrame {
  wf := WheelFrame{key, prompt, nil, 0}
  return &wf
}

func (wc *WheelController) addItem(nextI int, it WheelItemI) {
  if nextI >= 5 - 1 {
    var oldKey rune
    it, oldKey = reassignAndSwapKeys(it, makeKey(0))
    slider := makeWheelFrame(oldKey, ">>")
    wc.Current.addItem(slider)
    wc.Current = slider
  }
  wc.Current.addItem(it)
}

func reassignAndSwapKeys(it WheelItemI, toKey rune) (WheelItemI, rune) {
  var oldKey rune
  var wi WheelItemI
  switch it.(type) {
  case *WheelFrame:
    oldF := it.(*WheelFrame)
    oldKey = oldF.Key
    oldPrompt := oldF.Prompt
    wi = makeWheelFrame(toKey, oldPrompt)
  case WheelChatOption:
    oldF := it.(WheelChatOption)
    oldKey = oldF.Key
    oldPrompt := oldF.Prompt
    oldContent := oldF.Text
    wi = makeWheelChatOption(toKey, oldPrompt, oldContent)
  default:
    log.Fatal("Encountered wrong type while trying to parse data")
  }
  return wi, oldKey
}

func isOutOfBounds[T any](i int, arr []T) bool {
  return i < 0 || i >= len(arr)
}

var ErrParsePhraseSet = errors.New("Couldn't parse the file")

func parseWheelFrameInto(wc *WheelController, ps PhraseSet) {
  if len(ps) == 0 {
    log.Fatal("Encountered empty phrase set")
  }

  for _, p := range ps {
    prompt := p[0].(string)
    nextI := len(wc.Current.Items)
    key := makeKey(nextI)

    phrase, isString := p[1].(string)
    if isString {
      wco := makeWheelChatOption(key, prompt, phrase)
      wc.addItem(nextI, wco)
      continue 
    } 

    anyArr, isArr := p[1].([]any)
    if !isArr {
      log.Fatal(ErrParsePhraseSet)
    }

    pset := parseAnyArrIntoPhraseSet(anyArr)

    wf := makeWheelFrame(key, prompt)
    wc.addItem(nextI, wf)
    oldCur := wc.Current

    wc.Current = wf
    parseWheelFrameInto(wc, pset)
    wc.Current = oldCur
  }
}

func parseAnyArrIntoPhraseSet(anyArr []any) PhraseSet {
  pset := make(PhraseSet, 0)
  for _, v := range anyArr {
    slice, ok := v.([]any)
    if !ok {
      log.Fatal(ErrParsePhraseSet)
    }
    pair := [2]any{}
    for i := range pair {
      pair[i] = slice[i]
    }

    pset = append(pset, pair)
  }
  return pset
}

func parseWheel(ps PhraseSet) ChatWheelI {
  wc := new(WheelController)
  wc.Current = makeWheelFrame('0', "")
  wc.Start = wc.Current
  parseWheelFrameInto(wc, ps)
  wc.Current = nil
  return wc
}

func removePreviousCharacters(n int) {
  robotgo.KeyDown(robotgo.Shift)
  for i := 0; i < n+1; i++ {
    robotgo.KeyPress(robotgo.Left)
  }
  robotgo.KeyUp(robotgo.Shift)
  robotgo.KeyPress(robotgo.Backspace)
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
  ReactOnKey(hook.Event) (error)
}


func setupKeyboardListener(chw ChatWheelI) {
  fmt.Println("--- Press q to enter into active mode ---")

  modeActive := false
  hook.Register(hook.KeyDown, []string{}, func(e hook.Event) {
    chw.ReactOnKey(e)
  })
  hook.Register(hook.KeyDown, []string{"q"}, func(e hook.Event) {
    if !modeActive {
      modeActive = true
    }
  })

	s := hook.Start()
	<-hook.Process(s)
}

func isBackspace(char rune) bool {
  return rune(8) == char
}

const ACTIVATING_CHAR = 'q'


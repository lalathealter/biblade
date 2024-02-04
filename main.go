package main

import (
	"encoding/json"
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

  
  chw := produceReactCollectionFrom(ps)
  setupKeyboardListener(chw)
}

func readClipboard() string {
  res, err := clipboard.ReadAll()
  if err != nil {
    log.Fatal(err)
  }

  return res
}

type PhraseSet [][2]string

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

type ReactMap struct {
  Reactions []ReactOptionI
  FrameSize int
}

func makeReactMap(sectionCap int) *ReactMap {
  rmap := ReactMap{
    make([]ReactOptionI, 0, sectionCap),
    0,
  }
  return &rmap
}

func (rmap *ReactMap) addReactOption(reopt ReactOptionI) {
  rmap.Reactions = append(rmap.Reactions, reopt)
}

func (rmap *ReactMap) GetReaction(keyChar rune) (ReactOptionI, error) {
  ind := deductIndFrom(keyChar)
  if ind < 0 || ind >= len(rmap.Reactions) {
    return nil, ErrNoReactionForKey{keyChar}
  }
  return rmap.Reactions[ind], nil
}

func deductIndFrom(keyChar rune) int {
  return int(keyChar - '0') - 1
}

func makeKey(n int) rune {
  return rune('0'+n + 1)
}

type ReactCollection struct {
  CurrentSection int
  Sections []*ReactMap
  SectionCap int
  activatingChar rune
  isActive bool
}

func (rc *ReactCollection) getCurrSection() ReactMap {
  return *rc.Sections[rc.CurrentSection]
}

func (rc *ReactCollection) setCurrSection(n int) {
  rc.CurrentSection = n % rc.SectionCap
}

func (rc *ReactCollection) SetActiveMode(v bool) {
  rc.isActive = v
}

func (rc *ReactCollection) ReactOnKey(ev hook.Event) error {
  kchar := ev.Keychar
  rmap := rc.getCurrSection()
  if !rc.isActive {
    if kchar == rc.activatingChar {
      rc.SetActiveMode(true)
      rc.IntroduceChatOptions(0)
    }
    return nil
  } else if kchar == rc.activatingChar {
    rc.SetActiveMode(false)
    removePreviousCharacters(rmap.FrameSize)
    return nil
  }

  react, err := rmap.GetReaction(kchar)
  if err != nil {
    return err
  }

  removePreviousCharacters(rmap.FrameSize)
  react.Response()
  _, needsFollowingActions := react.(MoveOption)
  rc.SetActiveMode(needsFollowingActions)

  return nil
}

func (rc *ReactCollection) IntroduceChatOptions(n int) {
  rc.setCurrSection(n)
  for _, reopt := range rc.getCurrSection().Reactions {
    robotgo.TypeStr(reopt.GetTag())
  }
}

func (rc *ReactCollection) addSection(section *ReactMap) {
  rc.Sections = append(rc.Sections, section)
}

func produceReactCollectionFrom(ps PhraseSet) *ReactCollection {
  rcoll := &ReactCollection{0, nil, 5, 'q', false}
  sCap := rcoll.SectionCap

  reactI := 0
  currSect := makeReactMap(sCap)
  sectI := 0
  for n, pair := range ps {
    var reopt ReactOptionI
    if reactI >= sCap - 1 && n < len(ps) - 1 {
      reopt = formNextSetButton(makeKey(reactI), sectI+1, rcoll)
      currSect.addReactOption(reopt)
      currSect.FrameSize += len(reopt.GetTag())
      rcoll.addSection(currSect)
      
      sectI++
      currSect = makeReactMap(sCap)
      reactI = 0
    }

    tag, phrase := pair[0], pair[1]
    key := makeKey(reactI)
    opt := formIntroOptionText(tag, key)

    currSect.FrameSize += len(opt)
    reopt = formReactOption(opt, phrase)
    currSect.addReactOption(reopt)
    reactI++
  }

  sectsN := 1 + len(ps) / (sCap - 1)
  if sectI < sectsN {
    rcoll.addSection(currSect)
  }

  return rcoll
}


type MoveOption struct {
  Tag string
  GoToNextSection func()
}

func (mo MoveOption) GetTag() string {
  return mo.Tag
}

func (mo MoveOption) Response() {
  mo.GoToNextSection()
}

func formNextSetButton(key rune, n int, rcoll *ReactCollection) MoveOption {
  tag := formIntroOptionText(">>", key)
  return MoveOption{ 
    tag, func() { rcoll.IntroduceChatOptions(n)},
  }
}

type ReactOption struct {
  Tag string
  Content string
}

func (reopt ReactOption) GetTag() string {
  return reopt.Tag
}

func (reopt ReactOption) Response() {
  robotgo.TypeStr(reopt.Content)
}

func formReactOption(optTag string, v string) ReactOptionI {
  reopt := ReactOption{optTag, v}
  return reopt
}

func removePreviousCharacters(n int) {
  robotgo.KeyUp(robotgo.Ctrl)
  robotgo.KeyDown(robotgo.Shift)
  for i := 0; i < n + 2; i++ {
    // deleting both activating characters and frame; hence n + 2
    robotgo.KeyPress(robotgo.Left)
  }
  robotgo.KeyUp(robotgo.Shift)
  robotgo.KeyPress(robotgo.Backspace)
}

const MAX_INTRO_LEN = 12
func formIntroOptionText(name string, key rune) string {
  if len(name) > MAX_INTRO_LEN {
    name = name[:MAX_INTRO_LEN-3] + "..."
  }
  return fmt.Sprintf("[%v] %v ", string(key), name)
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
}




func setupKeyboardListener(chw ChatWheelI) {
  fmt.Println("--- Press q to enter into active mode ---")

  hook.Register(hook.KeyDown, []string{}, func(e hook.Event) {
    chw.ReactOnKey(e)
  })

	s := hook.Start()
	<-hook.Process(s)
}



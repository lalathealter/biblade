package main

import (
	"fmt"

	"github.com/go-vgo/robotgo"
	hook "github.com/robotn/gohook"
)

type ErrNoReactionForKey struct {
	Key rune
}

func (keyerr ErrNoReactionForKey) Error() string {
	return fmt.Sprint("Ignoring key ", keyerr.Key)
}

type ChatWheelI interface {
	ReactOnKey(hook.Event) error
}

type WheelController struct {
	Start   *WheelFrame
	Current *WheelFrame
}

func (wc *WheelController) getCurrOpts() []WheelItemI {
	return wc.Current.Items
}

const ACTIVATING_CHAR = 'q'

func (wc *WheelController) ReactOnKey(ev hook.Event) error {
	kchar := ev.Keychar
	if isBackspace(kchar) {
		return nil
	}

	if wc.Current == nil {
		if kchar == ACTIVATING_CHAR {
			removePreviousCharacters(1)
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

func removePreviousCharacters(n int) {
	for i := 0; i < n+1; i++ {
		robotgo.KeyPress(robotgo.Left, robotgo.Shift)
	}
	robotgo.KeyPress(robotgo.Backspace)
}

func (wc *WheelController) addItem(nextI int, it WheelItemI) {
	if nextI >= 5-1 {
		var oldKey rune
		it, oldKey = reassignAndSwapKeys(it, makeKey(0))
		slider := makeWheelFrame(oldKey, ">>")
		wc.Current.addItem(slider)
		wc.Current = slider
	}
	wc.Current.addItem(it)
}

type WheelItemI interface {
	Response() *WheelFrame
	GetTag() string
}

func makeWheelFrame(key rune, prompt string) *WheelFrame {
	wf := WheelFrame{key, prompt, nil, 0}
	return &wf
}

type WheelFrame struct {
	Key       rune
	Prompt    string
	Items     []WheelItemI
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
	Key    rune
	Prompt string
	Text   string
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

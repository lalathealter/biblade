package main

import (
	"encoding/json"
	"errors"
	"log"
	"os"
)

type PhraseSet [][2]any // string | PhraseSet

func loadFileInto[T any](path string) (T, error) {
	f, err := os.ReadFile(path)
	if err != nil {
		return *new(T), err
	}

	phs := new(T)
	err = json.Unmarshal(f, &phs)
	return *phs, err
}

func loadPhraseSet(path string) (PhraseSet, error) {
  return loadFileInto[PhraseSet](path)
}

func parseWheel(ps PhraseSet) ChatWheelI {
	wc := new(WheelController)
	wc.Current = makeWheelFrame('0', "")
	wc.Start = wc.Current
	parseWheelFrameInto(wc, ps)
	wc.Current = nil
	return wc
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

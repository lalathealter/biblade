package main

import (
	"log"
)

const pathToSettings = "./settings.json"

var Settings map[string]any

const ACTIVATING_CHAR_KEY = "key"
func getActivatingChar() rune {
  str := Settings[ACTIVATING_CHAR_KEY].(string)
  return rune(str[0])
}

const FRAME_CAP_KEY = "frameCap"
func getFrameCap() int {
  return dealWithInt(Settings[FRAME_CAP_KEY], 4)
}

const INTRO_LEN_KEY = "introLen"
func getIntroLen() int {
  return dealWithInt(Settings[INTRO_LEN_KEY], 4)
}

// Golang parses integer literal from json as float64 (for some reason);
// so this helps to convert any values from settings map to a proper int
func dealWithInt(v any, clampV int) int {
  f := v.(float64)
  i := int(f)
  return clampPositive(i, clampV)
}

const WHEEL_FILE_KEY = "wheelFile"
func getCurrWheelFile() string {
  return Settings[WHEEL_FILE_KEY].(string)
}

func init() {
  sets, err := loadSettings(pathToSettings)
  if err != nil {
    log.Fatal(err)
  }

  Settings = sets
}

func loadSettings(path string) (map[string]any, error) {
  return loadFileInto[map[string]any](path)
}

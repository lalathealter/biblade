package main

import "fmt"

func isOutOfBounds[T any](i int, arr []T) bool {
	return i < 0 || i >= len(arr)
}

func clampPositive(v int, lower int) int {
  if v < lower {
    v = lower
  }
  return v
}

func isBackspace(char rune) bool {
	return rune(8) == char
}

func makeTag(key rune, prompt string) string {
	return fmt.Sprintf("[%v] %v ", string(key), prompt)
}

func deductIndFrom(keyChar rune) int {
	return int(keyChar-'0') - 1
}

func makeKey(n int) rune {
	return rune('0' + n + 1)
}

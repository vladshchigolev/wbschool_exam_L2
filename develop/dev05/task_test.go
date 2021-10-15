package main

import "testing"

func TestAfter(t *testing.T) {
	phrase := "солнце"
	text := "стол рука чашка\nдом солнце игра\nвысоко мышь луч\nцифра жук зебра\nмышь его говори"
	count := 1
	if After(phrase, text, count) != "дом солнце игра\nвысоко мышь луч\n" {
		t.Error("результат не совпал с ожидаемым значением", After(phrase, text, count))
	}
}

func TestBefore(t *testing.T) {
	phrase := "жук"
	text := "стол рука чашка\nдом солнце игра\nвысоко мышь луч\nцифра жук зебра\nмышь его говори"
	count := 3
	if Before(phrase, text, count) != "стол рука чашка\nдом солнце игра\nвысоко мышь луч\nцифра жук зебра" {
		t.Error("результат не совпал с ожидаемым значением:", Before(phrase, text, count))
	}
}

func TestContextText(t *testing.T) {
	phrase := "солнце"
	text := "стол рука чашка\nдом солнце игра\nвысоко мышь луч\nцифра жук зебра\nмышь его говори"
	count := 1
	if ContextText(phrase, text, count) != "стол рука чашка\nдом солнце игра\nвысоко мышь луч" {
		t.Error("результат не совпал с ожидаемым значением:", ContextText(phrase, text, count))
	}
}

func TestCount(t *testing.T) {
	text := "стол рука чашка\nдом солнце игра\nвысоко мышь луч\nцифра жук зебра\nмышь его говори"
	if Count(text) != 5 {
		t.Error("результат не совпал с ожидаемым значением:", Count(text))
	}
}

func TestIgnoreCase(t *testing.T) {
	phrase := "солНЦЕ"
	text := "стол рука чашка\nдом сОлнце игра\nвысоко мышь луч\nцифра жук зебра\nмышь его говори"
	after = 2
	if IgnoreCase(phrase, text, after, before, contextText) != "дом солнце игра\nвысоко мышь луч\nцифра жук зебра\n" {
		t.Error("результат не совпал с ожидаемым значением:", IgnoreCase(phrase, text, after, before, contextText))
	}
}

func TestInvert(t *testing.T) {
	phrase := "мышь"
	text := "стол рука чашка\nдом солнце игра\nвысоко мышь луч\nцифра жук зебра\nмышь его говори"
	if Invert(phrase, text) != "стол рука чашка\nдом солнце игра\nцифра жук зебра\n" {
		t.Error("результат не совпал с ожидаемым значением:", Invert(phrase, text))
	}
}
package main

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestSortCharsInWords(t *testing.T) {
	testCases := []struct {
		in  string
		out string
	}{
		{
			"пятка",
			"акптя",
		},
		{
			"листок",
			"иклост",
		},
	}

	for _, v := range testCases {
		result := SortCharsInWord(v.in)
		assert.Equal(t, result, v.out)
	}
}

func TestFindAnagrams(t *testing.T) {
	testCases := []struct {
		in  []string
		out map[string][]string
	}{
		{
			[]string{"пятка", "пятак", "тяпка", "листок", "слиток", "столик", "ЛИСТОК", "ТяПкА"},
			map[string][]string{
				"пятка":  []string{"пятак", "пятка", "тяпка"},
				"листок": []string{"листок", "слиток", "столик"},
			},
		},
		{
			[]string{"бурак", "бурка", "рубка", "ДЕКОР", "ДОКЕР", "КРЕДО"},
			map[string][]string{
				"бурак" : []string{"бурак", "бурка", "рубка"},
				"ДЕКОР" : []string{"декор", "докер", "кредо"},
			},
		},
	}

	for _, v := range testCases {
		result := FindAnagrams(v.in)
		assert.Equal(t, result, v.out)
	}
}
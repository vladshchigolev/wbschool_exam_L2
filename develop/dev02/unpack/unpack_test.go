package unpack

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestHelper_Unpack(t *testing.T) {
	testCases := []struct {
		name    string
		actual  string
		expected  string
		isValid bool
	}{
		{
			name:    "Обычная распаковка",
			actual:  "a4bc2d5e",
			expected:  "aaaabccddddde",
			isValid: true,
		},
		{
			name:    "Обычная распаковка без цифр",
			actual:  "abcd",
			expected:  "abcd",
			isValid: true,
		},
		{
			name:    "Некорректная строка",
			actual:  "45",
			isValid: false,
		},
		{
			name:    "Пустая строка",
			actual:  "",
			expected:  "",
			isValid: true,
		},
		{
			name:    "Корректная, с экранированными цифрами",
			actual:  `qwe\4\5`,
			expected:  "qwe45",
			isValid: true,
		},
		{
			name:    "С экранированной цифрой, которую нужно распаковать",
			actual:  `qwe\45`,
			expected:  "qwe44444",
			isValid: true,
		},
		{
			name:    "С экранированной обратной косой чертой, которую нужно распаковать",
			actual:  `qwe\\5`,
			expected:  `qwe\\\\\`,
			isValid: true,
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) { // Run() вызывает анонимную функцию как сабтест в новой горутине
			act, err := Unpack(test.actual)
			if test.isValid {
				assert.NoError(t, err)
				assert.Equal(t, act, test.expected)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

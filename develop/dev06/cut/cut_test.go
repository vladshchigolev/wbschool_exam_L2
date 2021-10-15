package cut

import (
	"reflect"
	"testing"
)

func TestSetFields(t *testing.T) {
	mc := &ManCut{options: GetDefaultOptions()}

	testCases := []struct {
		name          string
		delimeter     string
		separated     bool
		fields        []int
		data          []string
		expected      []string
	}{
		{
			name:          "в качестве разделителя пробел",
			delimeter:     " ",
			separated: false,
			fields:        []int{0, 1, 2},
			data: []string{
				`стол рука чашка`,
				`дом солнце игра`,
			},
			expected: []string{
				`стол рука чашка`,
				`дом солнце игра`,
			},
		},
		{
			name:          "в качестве разделителя звездочка, только первые 2 поля",
			delimeter:     "*",
			separated: false,
			fields:        []int{0, 1},
			data: []string{
				`стол*рука*чашка`,
				`дом*солнце*игра`,
			},
			expected: []string{
				`стол*рука`,
				`дом*солнце`,
			},
		},
		{
			name:          "только строки с разделителем, только первые 2 поля",
			delimeter:     "#",
			separated: true,
			fields:        []int{0, 1},
			data: []string{
				`стол#рука#чашка`,
				`дом*солнце*игра`,
			},
			expected: []string{
				`стол#рука`,
			},
		},
	}

	for _, test := range testCases {
		t.Run(test.name, func(t *testing.T) {
			mc.result = []string{}
			mc.options.delimeter = test.delimeter
			mc.options.separated = test.separated
			mc.options.fields = test.fields
			mc.data = test.data

			SelectFields(mc)

			if !reflect.DeepEqual(mc.result, test.expected) {
				t.Errorf("actual %v, expected %v", mc.result, test.expected)
			}
		})
	}

}
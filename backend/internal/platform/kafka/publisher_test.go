package kafka

import (
	"reflect"
	"testing"
)

func TestParseBrokers(t *testing.T) {
	cases := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "single broker",
			input:    "localhost:9092",
			expected: []string{"localhost:9092"},
		},
		{
			name:     "multiple brokers",
			input:    "kafka-1:9092, kafka-2:9092 ,kafka-3:9092",
			expected: []string{"kafka-1:9092", "kafka-2:9092", "kafka-3:9092"},
		},
		{
			name:     "empty values ignored",
			input:    ", , localhost:9092, ",
			expected: []string{"localhost:9092"},
		},
		{
			name:     "empty input",
			input:    "",
			expected: []string{},
		},
	}

	for _, tc := range cases {
		got := ParseBrokers(tc.input)
		if !reflect.DeepEqual(got, tc.expected) {
			t.Fatalf("%s: expected %v, got %v", tc.name, tc.expected, got)
		}
	}
}

package config

import (
	"testing"
)

func TestSortFiles(t *testing.T) {
	tests := []struct {
		name     string
		input    []string
		expected []string
	}{
		{
			name:     "Numeric ordering (1, 2, 10)",
			input:    []string{"10_posts.sql", "1_init.sql", "2_users.sql"},
			expected: []string{"1_init.sql", "2_users.sql", "10_posts.sql"},
		},
		{
			name:     "Timestamp ordering",
			input:    []string{"20231024000000_init.sql", "20231023000000_setup.sql"},
			expected: []string{"20231023000000_setup.sql", "20231024000000_init.sql"},
		},
		{
			name:     "Mixed numeric and non-numeric",
			input:    []string{"schema.sql", "1_init.sql", "2_users.sql"},
			expected: []string{"1_init.sql", "2_users.sql", "schema.sql"},
		},
		{
			name:     "Non-numeric lexicographic fallback",
			input:    []string{"b.sql", "a.sql", "c.sql"},
			expected: []string{"a.sql", "b.sql", "c.sql"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			files := make([]string, len(tt.input))
			copy(files, tt.input)
			sortFiles(files)

			for i := range files {
				if files[i] != tt.expected[i] {
					t.Errorf("at index %d: expected %s, got %s", i, tt.expected[i], files[i])
				}
			}
		})
	}
}

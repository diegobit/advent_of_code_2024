package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFindUncorruptedMul(t *testing.T) {
	corruptedString := "xmul(2,4)%&mul[3,7]!@^do_not_mul(5,5)+mul(32,64]then(mul(11,8)mul(8,5))"

	expectedUncorruptedMuls := []string{"mul(2,4)", "mul(5,5)", "mul(11,8)", "mul(8,5)"}

	uncorruptedMuls := ExtractUncorruptedMuls(corruptedString)

	assert.Equal(t, expectedUncorruptedMuls, uncorruptedMuls)
}

func TestProcessMul(t *testing.T) {
	muls := []string{"mul(2,4)", "mul(5,5)", "mul(11,8)", "mul(8,5)"}

	result := ProcessMuls(muls)

	expected := 161

	assert.Equal(t, expected, result)
}

func BenchmarkProcessMul(b *testing.B) {

	for range b.N {
		path := "./data/input_3.txt"

		corruptedStrings, err := ReadCorruptedStrings(path)
		require.NoError(b, err)
		ComputeFirstResult(corruptedStrings)
	}

}

func TestExtractMulsWithInstructions(t *testing.T) {
	corruptedString := "xmul(2,4)&mul[3,7]!^don't()_mul(5,5)+mul(32,64](mul(11,8)undo()?mul(8,5))"

	expectedUncorruptedStrings := []string{"mul(2,4)", "don't()", "mul(5,5)", "mul(11,8)", "do()", "mul(8,5)"}

	uncorruptedMulsWithInstructions := ExtractUncorruptedMulsWithInstructions(corruptedString)

	assert.Equal(t, expectedUncorruptedStrings, uncorruptedMulsWithInstructions)
}

func TestProcessMulsWithInstructions(t *testing.T) {
	testCases := []struct {
		description          string
		mulsWithInstructions []string
		expected             int
		expectedProcess      bool
	}{
		{
			description:          "mul - don't - do - mul",
			mulsWithInstructions: []string{"mul(2,4)", "don't()", "mul(5,5)", "mul(11,8)", "do()", "mul(8,5)"},
			expected:             48,
			expectedProcess:      true,
		},
		{
			description:          "do - mul - don't",
			mulsWithInstructions: []string{"do()", "mul(8,5)", "don't()", "mul(5,5)", "mul(11,8)"},
			expected:             40,
			expectedProcess:      false,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.description, func(t *testing.T) {
			got, gotProcess := ProcessMulsWithInstructions(testCase.mulsWithInstructions, true)
			assert.Equal(t, testCase.expected, got)
			assert.Equal(t, testCase.expectedProcess, gotProcess)
		})
	}
}

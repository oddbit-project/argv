package argv

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

type ArgStructInt struct {
	Arg1 int8  `argv:"arg1"`
	Arg2 int32 `argv:"arg2"`
	Arg3 int   `argv:"arg3"`
	Arg4 int64 `argv:"arg4"`
}

type ArgStructUint struct {
	Arg1 uint8  `argv:"arg1"`
	Arg2 uint32 `argv:"arg2"`
	Arg3 uint   `argv:"arg3"`
	Arg4 uint64 `argv:"arg4"`
}

type ArgStructFloat struct {
	Arg1 float32 `argv:"arg1"`
	Arg2 float64 `argv:"arg2"`
}

type ArgStructTime struct {
	Arg1 time.Time `argv:"arg1"`
}

type ArgStructBool struct {
	Arg1 bool `argv:"arg1"`
	Arg2 bool `argv:"arg2"`
}

type ArgStructString struct {
	Arg1 string   `argv:"arg1"`
	Arg2 []string `argv:"arg2"`
}

func TestParseArgvInitErrors(t *testing.T) {

	payload := make([]string, 0)
	dest := &ArgStructInt{}

	// empty payload, should return ErrEmptyArgs
	err := ParseArgv(dest, payload)
	assert.ErrorIs(t, ErrEmptyArgs, err)

	// odd arg count, should return ErrEmptyArgs
	payload = []string{"param1", "value1", "param2"}
	err = ParseArgv(dest, payload)
	assert.ErrorIs(t, ErrInvalidParameterCount, err)

	// even arg count, should work
	payload = []string{"param1", "value1", "param2"}
	err = ParseArgv(dest, payload)
	assert.ErrorIs(t, ErrInvalidParameterCount, err)
}

func TestParseArgvInt(t *testing.T) {
	testCases := []struct {
		name           string
		args           []string
		expected       error
		expectedValues ArgStructInt
	}{
		{
			name:           "missing arg1",
			args:           []string{"arg2", "2"},
			expected:       ErrMissingValue("arg1"),
			expectedValues: ArgStructInt{},
		},
		{
			name:           "missing arg2",
			args:           []string{"arg1", "2"},
			expected:       ErrMissingValue("arg2"),
			expectedValues: ArgStructInt{},
		},
		{
			name:           "missing arg3",
			args:           []string{"arg1", "2", "arg2", "3"},
			expected:       ErrMissingValue("arg3"),
			expectedValues: ArgStructInt{},
		},
		{
			name:           "invalid arg1 value",
			args:           []string{"arg1", "xxx"},
			expected:       fmt.Errorf("error parsing arg arg1: strconv.ParseInt: parsing \"xxx\": invalid syntax"),
			expectedValues: ArgStructInt{},
		},
		{
			name:           "invalid arg4 value",
			args:           []string{"arg1", "2", "arg2", "3", "arg3", "45", "arg4", "potato"},
			expected:       fmt.Errorf("error parsing arg arg4: strconv.ParseInt: parsing \"potato\": invalid syntax"),
			expectedValues: ArgStructInt{},
		},
		{
			name:           "overflow arg1",
			args:           []string{"arg1", "2000"},
			expected:       fmt.Errorf("error parsing arg arg1: strconv.ParseInt: parsing \"2000\": value out of range"),
			expectedValues: ArgStructInt{},
		},
		{
			name:     "valid parsing",
			args:     []string{"arg1", "127", "arg2", "37000", "arg3", "4532", "arg4", "-476800"},
			expected: nil,
			expectedValues: ArgStructInt{
				Arg1: 127,
				Arg2: 37000,
				Arg3: 4532,
				Arg4: -476800,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			dest := &ArgStructInt{}
			err := ParseArgv(dest, tc.args)
			if tc.expected != nil {
				if err == nil {
					t.Errorf("expected %v, got <nil>", tc.expected)
				} else if err.Error() != tc.expected.Error() {
					t.Errorf("expected %v, got %v", tc.expected, err)
				}
			}
			if tc.expected == nil {
				assert.Equal(t, &tc.expectedValues, dest)
			}
		})
	}

}

func TestParseArgvUint(t *testing.T) {
	testCases := []struct {
		name           string
		args           []string
		expected       error
		expectedValues ArgStructUint
	}{
		{
			name:           "missing arg1",
			args:           []string{"arg2", "2"},
			expected:       ErrMissingValue("arg1"),
			expectedValues: ArgStructUint{},
		},
		{
			name:           "missing arg2",
			args:           []string{"arg1", "2"},
			expected:       ErrMissingValue("arg2"),
			expectedValues: ArgStructUint{},
		},
		{
			name:           "missing arg3",
			args:           []string{"arg1", "2", "arg2", "3"},
			expected:       ErrMissingValue("arg3"),
			expectedValues: ArgStructUint{},
		},
		{
			name:           "invalid arg1 value",
			args:           []string{"arg1", "xxx"},
			expected:       fmt.Errorf("error parsing arg arg1: strconv.ParseUint: parsing \"xxx\": invalid syntax"),
			expectedValues: ArgStructUint{},
		},
		{
			name:           "invalid arg4 value",
			args:           []string{"arg1", "2", "arg2", "3", "arg3", "45", "arg4", "potato"},
			expected:       fmt.Errorf("error parsing arg arg4: strconv.ParseUint: parsing \"potato\": invalid syntax"),
			expectedValues: ArgStructUint{},
		},
		{
			name:           "overflow arg1",
			args:           []string{"arg1", "2000"},
			expected:       fmt.Errorf("error parsing arg arg1: strconv.ParseUint: parsing \"2000\": value out of range"),
			expectedValues: ArgStructUint{},
		},
		{
			name:           "negative arg2",
			args:           []string{"arg1", "100", "arg2", "-5", "arg3", "2345", "arg4", "324654"},
			expected:       fmt.Errorf("error parsing arg arg2: strconv.ParseUint: parsing \"-5\": invalid syntax"),
			expectedValues: ArgStructUint{},
		},
		{
			name:     "valid parsing",
			args:     []string{"arg1", "127", "arg2", "37000", "arg3", "4532", "arg4", "476800"},
			expected: nil,
			expectedValues: ArgStructUint{
				Arg1: 127,
				Arg2: 37000,
				Arg3: 4532,
				Arg4: 476800,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			dest := &ArgStructUint{}
			err := ParseArgv(dest, tc.args)
			if tc.expected != nil {
				if err == nil {
					t.Errorf("expected %v, got <nil>", tc.expected)
				} else if err.Error() != tc.expected.Error() {
					t.Errorf("expected %v, got %v", tc.expected, err)
				}
			}
			if tc.expected == nil {
				assert.Equal(t, &tc.expectedValues, dest)
			}
		})
	}

}

func TestParseArgvFloat(t *testing.T) {
	testCases := []struct {
		name           string
		args           []string
		expected       error
		expectedValues ArgStructFloat
	}{
		{
			name:           "missing arg1",
			args:           []string{"arg2", "2.0"},
			expected:       ErrMissingValue("arg1"),
			expectedValues: ArgStructFloat{},
		},
		{
			name:           "missing arg2",
			args:           []string{"arg1", "1.3E13"},
			expected:       ErrMissingValue("arg2"),
			expectedValues: ArgStructFloat{},
		},
		{
			name:           "invalid arg1 value",
			args:           []string{"arg1", "xxx"},
			expected:       fmt.Errorf("error parsing arg arg1: strconv.ParseFloat: parsing \"xxx\": invalid syntax"),
			expectedValues: ArgStructFloat{},
		},
		{
			name:           "invalid arg2 value",
			args:           []string{"arg1", "2", "arg2", "potato"},
			expected:       fmt.Errorf("error parsing arg arg2: strconv.ParseFloat: parsing \"potato\": invalid syntax"),
			expectedValues: ArgStructFloat{},
		},
		{
			name:           "overflow arg1",
			args:           []string{"arg1", "10E140"},
			expected:       fmt.Errorf("error parsing arg arg1: strconv.ParseFloat: parsing \"10E140\": value out of range"),
			expectedValues: ArgStructFloat{},
		},
		{
			name:     "valid parsing",
			args:     []string{"arg1", "45.72", "arg2", "-23453453.32423"},
			expected: nil,
			expectedValues: ArgStructFloat{
				Arg1: 45.72,
				Arg2: -23453453.32423,
			},
		},
		{
			name:     "valid parsing, scientific notation",
			args:     []string{"arg1", "-2.345E8", "arg2", "8.3732E18"},
			expected: nil,
			expectedValues: ArgStructFloat{
				Arg1: -2.345e+08,
				Arg2: 8.3732e+18,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			dest := &ArgStructFloat{}
			err := ParseArgv(dest, tc.args)
			if tc.expected != nil {
				if err == nil {
					t.Errorf("expected %v, got <nil>", tc.expected)
				} else if err.Error() != tc.expected.Error() {
					t.Errorf("expected %v, got %v", tc.expected, err)
				}
			}
			if tc.expected == nil {
				assert.Equal(t, &tc.expectedValues, dest)
			}
		})
	}
}

func TestParseArgvTime(t *testing.T) {
	validTime1, _ := time.Parse(time.RFC3339, "2006-01-02T15:04:05Z07:00")
	validTime2, _ := time.Parse(time.RFC3339, "2023-05-25 00:10:01-02:00")

	testCases := []struct {
		name           string
		args           []string
		expected       error
		expectedValues ArgStructTime
	}{
		{
			name:           "missing arg1",
			args:           []string{"arg2", "2.0"},
			expected:       ErrMissingValue("arg1"),
			expectedValues: ArgStructTime{},
		},
		{
			name:           "invalid arg1",
			args:           []string{"arg1", "2.0"},
			expected:       ErrInvalidValue("arg1", fmt.Errorf("parsing time \"2.0\" as \"2006-01-02T15:04:05Z07:00\": cannot parse \"2.0\" as \"2006\"")),
			expectedValues: ArgStructTime{},
		},
		{
			name:           "invalid arg1",
			args:           []string{"arg1", "potato"},
			expected:       ErrInvalidValue("arg1", fmt.Errorf("parsing time \"potato\" as \"2006-01-02T15:04:05Z07:00\": cannot parse \"potato\" as \"2006\"")),
			expectedValues: ArgStructTime{},
		},
		{
			name:     "valid parsing",
			args:     []string{"arg1", "2023-05-25 00:10:01-02:00"},
			expected: nil,
			expectedValues: ArgStructTime{
				Arg1: validTime2,
			},
		},
		{
			name:     "valid parsing, scientific notation",
			args:     []string{"arg1", "2006-01-02T15:04:05Z07:00"},
			expected: nil,
			expectedValues: ArgStructTime{
				Arg1: validTime1,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			dest := &ArgStructTime{}
			err := ParseArgv(dest, tc.args)
			if tc.expected != nil {
				if err == nil {
					t.Errorf("expected %v, got <nil>", tc.expected)
				} else if err.Error() != tc.expected.Error() {
					t.Errorf("expected %v, got %v", tc.expected, err)
				}
			}
			if tc.expected == nil {
				assert.Equal(t, &tc.expectedValues, dest)
			}
		})
	}
}

func TestParseArgvBool(t *testing.T) {

	testCases := []struct {
		name           string
		args           []string
		expected       error
		expectedValues ArgStructBool
	}{
		{
			name:           "missing arg1",
			args:           []string{"arg2", "2.0"},
			expected:       ErrMissingValue("arg1"),
			expectedValues: ArgStructBool{},
		},
		{
			name:           "invalid arg1",
			args:           []string{"arg1", "2.0"},
			expected:       ErrInvalidValue("arg1", fmt.Errorf("strconv.ParseBool: parsing \"2.0\": invalid syntax")),
			expectedValues: ArgStructBool{},
		},
		{
			name:     "valid bool 1",
			args:     []string{"arg1", "true", "arg2", "false"},
			expected: nil,
			expectedValues: ArgStructBool{
				Arg1: true,
				Arg2: false,
			},
		},
		{
			name:     "valid bool 2",
			args:     []string{"arg1", "1", "arg2", "0"},
			expected: nil,
			expectedValues: ArgStructBool{
				Arg1: true,
				Arg2: false,
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			dest := &ArgStructBool{}
			err := ParseArgv(dest, tc.args)
			if tc.expected != nil {
				if err == nil {
					t.Errorf("expected %v, got <nil>", tc.expected)
				} else if err.Error() != tc.expected.Error() {
					t.Errorf("expected %v, got %v", tc.expected, err)
				}
			}
			if tc.expected == nil {
				assert.Equal(t, &tc.expectedValues, dest)
			}
		})
	}
}

func TestParseArgvString(t *testing.T) {

	testCases := []struct {
		name           string
		args           []string
		expected       error
		expectedValues ArgStructString
	}{
		{
			name:           "missing arg1",
			args:           []string{"arg2", "xxx"},
			expected:       ErrMissingValue("arg1"),
			expectedValues: ArgStructString{},
		},
		{
			name:     "valid strings, empty list",
			args:     []string{"arg1", "xxx", "arg2", ""},
			expected: nil,
			expectedValues: ArgStructString{
				Arg1: "xxx",
				Arg2: []string{},
			},
		},
		{
			name:     "valid strings",
			args:     []string{"arg1", "xxx", "arg2", "value1,value2,value3"},
			expected: nil,
			expectedValues: ArgStructString{
				Arg1: "xxx",
				Arg2: []string{"value1", "value2", "value3"},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			dest := &ArgStructString{}
			err := ParseArgv(dest, tc.args)
			if tc.expected != nil {
				if err == nil {
					t.Errorf("expected %v, got <nil>", tc.expected)
				} else if err.Error() != tc.expected.Error() {
					t.Errorf("expected %v, got %v", tc.expected, err)
				}
			}
			if tc.expected == nil {
				assert.Equal(t, &tc.expectedValues, dest)
			}
		})
	}
}

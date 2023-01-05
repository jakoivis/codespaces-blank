package tz

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_IsValidUtcOffset_and_UtcOffsetSecods(t *testing.T) {
	assert := assert.New(t)

	var tests = []struct {
		isValid   bool
		utcOffset string
		seconds   int
	}{
		// ColonSyntax
		{true, "+02:02", 7320},
		{true, "-02:02", -7320},
		{true, "UTC+02:02", 7320},
		{true, "utc+02:02", 7320},

		{false, "asd+02:02", 0},
		{false, "+2:02", 0},
		{false, "+2:2", 0},
		{false, "-2:02", 0},
		{false, "-2:2", 0},
		{false, "02:02", 0},
		{false, "+022:02", 0},
		{false, "+02:022", 0},
		{false, "+02:a2", 0},
		{false, "+a2:02", 0},

		// HourMinuteSyntax
		{true, "+0202", 7320},
		{true, "-0202", -7320},
		{true, "UTC+0202", 7320},
		{true, "utc+0202", 7320},

		{false, "asd+0202", 0},
		{false, "+02020", 0},
		{false, "+a202", 0},
		{false, "+02a2", 0},

		// HourSyntax
		{true, "+02", 7200},
		{true, "-02", -7200},
		{true, "UTC+02", 7200},
		{true, "utc+02", 7200},

		{false, "asd+02", 0},
		{false, "+2", 0},
		{false, "+022", 0},
		{false, "+a2", 0},
		{false, "02", 0},
	}

	for _, tt := range tests {

		t.Run(tt.utcOffset, func(t *testing.T) {
			isValid := IsValidUtcOffset(tt.utcOffset)
			assert.Equal(tt.isValid, isValid, "wrong validity value retured while testing "+tt.utcOffset)

			seconds, err := UtcOffsetSeconds(tt.utcOffset)
			assert.Equal(tt.seconds, seconds, "testing "+tt.utcOffset)

			hasError := false
			isNotValid := !tt.isValid

			if err != nil {
				hasError = true
			}

			assert.Equal(isNotValid, hasError, "wrong error value retured while testing "+tt.utcOffset)
		})
	}
}

func Test_parseColonSyntax_shouldConvertToSeconds(t *testing.T) {
	offsetSeconds, err := parseColonSyntax("UTC+02:02")
	assert.Equal(t, 7320, offsetSeconds)
	assert.NoError(t, err)
}

func Test_parseColonSyntax_shouldReturnErrorWhenHoursIsNotANumber(t *testing.T) {
	offsetSeconds, err := parseColonSyntax("UTC+AA:02")
	assert.Equal(t, 0, offsetSeconds)
	assert.Error(t, err)
	assert.Equal(t, "invalid UTC offset UTC+AA:02, strconv.Atoi: parsing \"AA\": invalid syntax", err.Error())
}

func Test_parseHourMinuteSyntax_shouldConvertToSeconds(t *testing.T) {
	offsetSeconds, err := parseHourMinuteSyntax("UTC+0202")
	assert.Equal(t, 7320, offsetSeconds)
	assert.NoError(t, err)
}

func Test_parseHourMinuteSyntax_shouldReturnErrorWhenHoursIsNotANumber(t *testing.T) {
	offsetSeconds, err := parseHourMinuteSyntax("UTC+AA02")
	assert.Equal(t, 0, offsetSeconds)
	assert.Error(t, err)
	assert.Equal(t, "invalid UTC offset UTC+AA02, strconv.Atoi: parsing \"AA\": invalid syntax", err.Error())
}

func Test_parseHourSyntax_shouldConvertToSeconds(t *testing.T) {
	offsetSeconds, err := parseHourSyntax("UTC+02")
	assert.Equal(t, 7200, offsetSeconds)
	assert.NoError(t, err)
}

func Test_parseHourSyntax_shouldReturnErrorWhenHoursIsNotANumber(t *testing.T) {
	offsetSeconds, err := parseHourSyntax("UTC+AA")
	assert.Equal(t, 0, offsetSeconds)
	assert.Error(t, err)
	assert.Equal(t, "invalid UTC offset UTC+AA, strconv.Atoi: parsing \"AA\": invalid syntax", err.Error())
}

func Test_hoursAndMinutesToSecods_shouldConvert(t *testing.T) {
	seconds, err := hoursAndMinutesToSecods("02", "02")
	assert.Equal(t, 7320, seconds)
	assert.NoError(t, err)
}

func Test_hoursAndMinutesToSecods_shouldReturnErrorWhenHoursCannotBeConverted(t *testing.T) {
	seconds, err := hoursAndMinutesToSecods("0A", "02")
	assert.Equal(t, 0, seconds)
	assert.Error(t, err)
}

func Test_hoursAndMinutesToSecods_shouldReturnErrorWhenMinutesCannotBeConverted(t *testing.T) {
	seconds, err := hoursAndMinutesToSecods("02", "0A")
	assert.Equal(t, 0, seconds)
	assert.Error(t, err)
}

func Test_stripUtcPrefix_shouldRemoveUTCPrefix(t *testing.T) {
	utcOffset := stripUtcPrefix("UTC-02")
	assert.Equal(t, "-02", utcOffset)

	utcOffset = stripUtcPrefix("utc-02")
	assert.Equal(t, "-02", utcOffset)

	utcOffset = stripUtcPrefix("utcc-02")
	assert.Equal(t, "c-02", utcOffset)
}

func Test_stripUtcPrefix_shouldNotRemoveAnything(t *testing.T) {
	utcOffset := stripUtcPrefix("UTC-02")
	assert.Equal(t, "-02", utcOffset)

	utcOffset = stripUtcPrefix("ASD-02")
	assert.Equal(t, "ASD-02", utcOffset)

	utcOffset = stripUtcPrefix("UUTC-02")
	assert.Equal(t, "UUTC-02", utcOffset)
}

func Test_stripSignPrefix_shouldRemoveTheSignAndReturnIt(t *testing.T) {
	utcOffset, direction := stripSignPrefix("-02")
	assert.Equal(t, "02", utcOffset)
	assert.Equal(t, -1, direction)

	utcOffset, direction = stripSignPrefix("+02:02")
	assert.Equal(t, "02:02", utcOffset)
	assert.Equal(t, 1, direction)
}

func Test_stripSignPrefix_shouldNotRemoveTheSignAndReturnZero(t *testing.T) {
	utcOffset, direction := stripSignPrefix("?02")
	assert.Equal(t, "?02", utcOffset)
	assert.Equal(t, 0, direction)
}

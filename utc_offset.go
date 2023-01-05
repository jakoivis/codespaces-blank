package utcoffset

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// matches UTC+02:02 +02:02
var ColonSyntax = regexp.MustCompile(`(?i)^(UTC)?[+-](\d\d:\d\d)$`)

// matches UTC+0202, +0202
var HourMinuteSyntax = regexp.MustCompile(`(?i)^(UTC)?[+-](\d{4})$`)

// matches UTC+02, +02
var HourSyntax = regexp.MustCompile(`(?i)^(UTC)?[+-](\d\d)$`)

// return boolean indicating whether the given UTC offset has valid format
// valid formats are ±[hh]:[mm], ±[hh][mm], or ±[hh] e.g. "+02:34", "+0234", "+02"
// Allows using case insensitive UTC prefix e.g. "UTC+02:34", "UTC+0234", "UTC+02"
func IsValidUtcOffset(utcOffset string) bool {

	return ColonSyntax.MatchString(utcOffset) ||
		HourMinuteSyntax.MatchString(utcOffset) ||
		HourSyntax.MatchString(utcOffset)
}

// returns utc offset in seconds
// utcOffset format ±[hh]:[mm], ±[hh][mm], or ±[hh] e.g. "+02:34", "+0234", "+02"
// Allows using case insensitive UTC prefix e.g. "UTC+02:34", "UTC+0234", "UTC+02"
func UtcOffsetSeconds(utcOffset string) (int, error) {

	if ColonSyntax.MatchString(utcOffset) {

		return parseColonSyntax(utcOffset)

	} else if HourMinuteSyntax.MatchString(utcOffset) {

		return parseHourMinuteSyntax(utcOffset)

	} else if HourSyntax.MatchString(utcOffset) {

		return parseHourSyntax(utcOffset)
	}

	return 0, fmt.Errorf("invalid UTC offset %v", utcOffset)
}

// format is ±[hh]:[mm] or UTC±[hh]:[mm]
func parseColonSyntax(utcOffset string) (int, error) {
	offset := stripUtcPrefix(utcOffset)
	offset, direction := stripSignPrefix(offset)
	parts := strings.Split(offset, ":")
	hours := parts[0]
	minutes := parts[1]
	seconds, err := hoursAndMinutesToSecods(hours, minutes)

	if err != nil {

		return 0, fmt.Errorf("invalid UTC offset %v, %v", utcOffset, err)
	}

	return seconds * direction, nil
}

// format is ±[hh][mm] or UTC±[hh][mm]
func parseHourMinuteSyntax(utcOffset string) (int, error) {
	offset := stripUtcPrefix(utcOffset)
	offset, direction := stripSignPrefix(offset)
	hours := offset[:2]
	minutes := offset[2:]
	seconds, err := hoursAndMinutesToSecods(hours, minutes)

	if err != nil {

		return 0, fmt.Errorf("invalid UTC offset %v, %v", utcOffset, err)
	}

	return seconds * direction, nil
}

// format is ±[hh] or UTC±[hh]
func parseHourSyntax(utcOffset string) (int, error) {
	offset := stripUtcPrefix(utcOffset)
	offset, direction := stripSignPrefix(offset)
	seconds, err := hoursAndMinutesToSecods(offset, "0")

	if err != nil {

		return 0, fmt.Errorf("invalid UTC offset %v, %v", utcOffset, err)
	}

	return seconds * direction, nil
}

func hoursAndMinutesToSecods(hoursStr string, minutesStr string) (int, error) {
	hours, err := strconv.Atoi(hoursStr)

	if err != nil {

		return 0, err
	}

	minutes, err := strconv.Atoi(minutesStr)

	if err != nil {

		return 0, err
	}

	return hours*60*60 + minutes*60, nil
}

// returns the UTC offset without the "UTC" word e.g. UTC+02 will become +02
// This package allows having the "UTC" prefix in the UTC offset syntax
// it is not needed anywhere so it can be just removed
func stripUtcPrefix(utcOffset string) string {

	if strings.ToLower(utcOffset[0:3]) == "utc" {

		return utcOffset[3:]
	}

	return utcOffset
}

// returns the UTC offset without the sign and returns the sign as 1 or -1 indicating the offset direction
// when calling this, it is assumed that the sign is the first character of the string
func stripSignPrefix(utcOffset string) (offset string, direction int) {
	switch utcOffset[0] {

	case '+':
		offset = utcOffset[1:]
		direction = 1

	case '-':
		offset = utcOffset[1:]
		direction = -1

	default:
		offset = utcOffset
		direction = 0
	}

	return
}

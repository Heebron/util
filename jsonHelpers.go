// Package util provides utility functions and types for common operations.
package util

import (
	"encoding/json"
	"strconv"
	"time"
)

// UnixTimeFromIntString is a custom type for unmarshaling JSON string values containing Unix timestamps
// into Go time.Time objects. It implements the json.Unmarshaler interface.
type UnixTimeFromIntString time.Time

// UnmarshalJSON implements the json.Unmarshaler interface for UnixTimeFromIntString.
// It parses a JSON string containing a Unix timestamp (seconds since epoch) and converts it to a time.Time.
//
// Parameters:
//   - data: The JSON data to unmarshal
//
// Returns:
//   - An error if unmarshaling or parsing fails, nil otherwise
func (ut *UnixTimeFromIntString) UnmarshalJSON(data []byte) error {
	var timestampStr string
	if err := json.Unmarshal(data, &timestampStr); err != nil {
		return err
	}

	timestamp, err := strconv.ParseInt(timestampStr, 10, 64)
	if err != nil {
		return err
	}

	*ut = UnixTimeFromIntString(time.Unix(timestamp, 0))

	return nil
}

// Value returns the underlying time.Time value from the UnixTimeFromIntString.
//
// Returns:
//   - The time.Time value represented by this UnixTimeFromIntString
func (ut *UnixTimeFromIntString) Value() time.Time {
	return time.Time(*ut)
}

// UnixTimeFromInt is a custom type for unmarshaling JSON numeric values containing Unix timestamps
// into Go time.Time objects. It implements the json.Unmarshaler interface.
type UnixTimeFromInt time.Time

// UnmarshalJSON implements the json.Unmarshaler interface for UnixTimeFromInt.
// It parses a JSON number containing a Unix timestamp (seconds since epoch) and converts it to a time.Time.
//
// Parameters:
//   - data: The JSON data to unmarshal
//
// Returns:
//   - An error if unmarshaling fails, nil otherwise
func (ut *UnixTimeFromInt) UnmarshalJSON(data []byte) error {
	var timestamp int64
	if err := json.Unmarshal(data, &timestamp); err != nil {
		return err
	}

	*ut = UnixTimeFromInt(time.Unix(timestamp, 0))

	return nil
}

// Value returns the underlying time.Time value from the UnixTimeFromInt.
//
// Returns:
//   - The time.Time value represented by this UnixTimeFromInt
func (ut *UnixTimeFromInt) Value() time.Time {
	return time.Time(*ut)
}

// Uint64FromString is a custom type for unmarshaling JSON string values containing
// unsigned 64-bit integers into Go uint64 values. It implements the json.Unmarshaler interface.
type Uint64FromString uint64

// UnmarshalJSON implements the json.Unmarshaler interface for Uint64FromString.
// It parses a JSON string containing a numeric value and converts it to a uint64.
//
// Parameters:
//   - data: The JSON data to unmarshal
//
// Returns:
//   - An error if unmarshaling or parsing fails, nil otherwise
func (ut *Uint64FromString) UnmarshalJSON(data []byte) error {
	var value string
	if err := json.Unmarshal(data, &value); err != nil {
		return err
	}

	if uintValue, err := strconv.ParseUint(value, 10, 64); err != nil {
		return err
	} else {
		*ut = Uint64FromString(uintValue)
	}

	return nil
}

// Value returns the underlying uint64 value from the Uint64FromString.
//
// Returns:
//   - The uint64 value represented by this Uint64FromString
func (ut *Uint64FromString) Value() uint64 {
	return uint64(*ut)
}

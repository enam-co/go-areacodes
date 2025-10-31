package areacodes

import (
	_ "embed"
	"encoding/json"
	"errors"
	"regexp"
	"strings"
	"sync"
)

// ErrInvalidPhone is returned when a phone number cannot be normalized into a valid NANPA number.
var ErrInvalidPhone = errors.New("phone number not valid")

// ErrNotFound is returned when an area code is not present in the embedded database.
var ErrNotFound = errors.New("not found")

// Location describes the centroid the original project associated with an area code.
type Location struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
}

// Entry stores metadata about a NANPA area code.
type Entry struct {
	Type      string   `json:"type"`
	City      string   `json:"city"`
	State     string   `json:"state"`
	StateCode string   `json:"stateCode"`
	Location  Location `json:"location"`
}

//go:embed data/data.json
var rawData []byte

var (
	once    sync.Once
	data    map[string]Entry
	loadErr error
)

var nonDigits = regexp.MustCompile(`[^0-9]`)

// Lookup normalizes the supplied phone number and returns the matching area code entry.
// The function accepts the same wide range of formats as the original Node.js implementation.
func Lookup(phone string) (Entry, error) {
	if err := ensureLoaded(); err != nil {
		return Entry{}, err
	}

	normalized, err := normalizePhone(phone)
	if err != nil {
		return Entry{}, err
	}

	entry, ok := data[normalized[:3]]
	if !ok {
		return Entry{}, ErrNotFound
	}

	return entry, nil
}

// All returns a defensive copy of the area code database.
func All() (map[string]Entry, error) {
	if err := ensureLoaded(); err != nil {
		return nil, err
	}

	cloned := make(map[string]Entry, len(data))
	for code, entry := range data {
		cloned[code] = entry
	}

	return cloned, nil
}

func ensureLoaded() error {
	once.Do(func() {
		var parsed map[string]Entry
		if err := json.Unmarshal(rawData, &parsed); err != nil {
			loadErr = err
			return
		}
		data = parsed
	})

	return loadErr
}

func normalizePhone(phone string) (string, error) {
	if phone == "" {
		return "", ErrInvalidPhone
	}

	digits := nonDigits.ReplaceAllString(phone, "")
	if len(digits) == 0 {
		return "", ErrInvalidPhone
	}

	if len(digits) > 10 && strings.HasPrefix(digits, "1") {
		digits = digits[1:]
	}

	if len(digits) != 10 {
		return "", ErrInvalidPhone
	}

	return digits, nil
}

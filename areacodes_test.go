package areacodes

import (
	"errors"
	"fmt"
	"testing"
)

func TestLookupSuccessFormats(t *testing.T) {
	db, err := All()
	if err != nil {
		t.Fatalf("All: %v", err)
	}

	for code, expected := range db {
		code := code
		expected := expected
		t.Run("code_"+code, func(t *testing.T) {
			phones := []string{
				fmt.Sprintf("+1-%s-123-4567", code),
				fmt.Sprintf("1-%s-123-4567", code),
				fmt.Sprintf("1 (%s) 123-4567", code),
				fmt.Sprintf("(%s) 123-4567", code),
				fmt.Sprintf(" (%s) 123-4567", code),
				fmt.Sprintf("%s 123-4567", code),
				fmt.Sprintf("%s123-4567", code),
				fmt.Sprintf("%s1234567", code),
				fmt.Sprintf("%s.123.4567", code),
				fmt.Sprintf("1.%s.123.4567", code),
				"+-" + code + "-123-4567",
			}

			for idx, phone := range phones {
				entry, err := Lookup(phone)
				if err != nil {
					t.Fatalf("format #%d (%q) returned error: %v", idx+1, phone, err)
				}
				if entry != expected {
					t.Fatalf("format #%d (%q) returned %+v, want %+v", idx+1, phone, entry, expected)
				}
			}
		})
	}
}

func TestLookupInvalidFormats(t *testing.T) {
	db, err := All()
	if err != nil {
		t.Fatalf("All: %v", err)
	}

	for code := range db {
		code := code
		t.Run("invalid_"+code, func(t *testing.T) {
			type expectation struct {
				phone string
				err   error
			}

			expectations := []expectation{
				{fmt.Sprintf("+1-%s-123-456", code), ErrNotFound},
				{fmt.Sprintf("1-%s-123-456", code), ErrNotFound},
				{fmt.Sprintf("0-%s-123-456", code), ErrNotFound},
				{fmt.Sprintf("0-%s-123-4567", code), ErrInvalidPhone},
				{fmt.Sprintf("0-%s-12-4567", code), ErrNotFound},
				{fmt.Sprintf("0-%s-12-4567", code), ErrNotFound},
			}

			for idx, exp := range expectations {
				_, err := Lookup(exp.phone)
				if !errors.Is(err, exp.err) {
					t.Fatalf("format #%d (%q) returned error %v, want %v", idx+12, exp.phone, err, exp.err)
				}
			}
		})
	}
}

func TestLookupNonexistentAreaCodes(t *testing.T) {
	db, err := All()
	if err != nil {
		t.Fatalf("All: %v", err)
	}

	known := make(map[string]struct{}, len(db))
	for code := range db {
		known[code] = struct{}{}
	}

	for i := 0; i < 1000; i++ {
		code := fmt.Sprintf("%03d", i)
		if _, exists := known[code]; exists {
			continue
		}

		t.Run("nonexistent_"+code, func(t *testing.T) {
			phones := []string{
				fmt.Sprintf("+1-%s-123-4567", code),
				fmt.Sprintf("1-%s-123-4567", code),
				fmt.Sprintf("1 (%s) 123-4567", code),
				fmt.Sprintf("(%s) 123-4567", code),
				fmt.Sprintf(" (%s) 123-4567", code),
				fmt.Sprintf("%s 123-4567", code),
				fmt.Sprintf("%s123-4567", code),
				fmt.Sprintf("%s1234567", code),
				fmt.Sprintf("%s.123.4567", code),
				fmt.Sprintf("1.%s.123.4567", code),
				"+-" + code + "-123-4567",
			}

			for idx, phone := range phones {
				_, err := Lookup(phone)
				if !errors.Is(err, ErrNotFound) {
					t.Fatalf("format #%d (%q) returned error %v, want %v", idx+1, phone, err, ErrNotFound)
				}
			}
		})
	}
}

func TestAllCopy(t *testing.T) {
	db, err := All()
	if err != nil {
		t.Fatalf("All: %v", err)
	}

	if len(db) == 0 {
		t.Fatal("expected non-empty database")
	}

	const firstCode = "201"

	// mutate returned map and ensure subsequent calls remain unchanged
	db[firstCode] = Entry{}

	next, err := All()
	if err != nil {
		t.Fatalf("All: %v", err)
	}

	if entry, ok := next[firstCode]; !ok || entry == (Entry{}) {
		t.Fatalf("database not defensively copied; got %+v", entry)
	}
}

func TestNormalizePhone(t *testing.T) {
	tests := []struct {
		input    string
		expected string
		err      error
	}{
		{"+1-303-123-4567", "3031234567", nil},
		{"(720) 555-1212", "7205551212", nil},
		{"", "", ErrInvalidPhone},
		{"000", "", ErrInvalidPhone},
		{"2-1-303-123-4567", "", ErrInvalidPhone},
	}

	for _, tt := range tests {
		actual, err := normalizePhone(tt.input)
		if tt.err == nil {
			if err != nil {
				t.Fatalf("normalizePhone(%q) returned error %v", tt.input, err)
			}
			if actual != tt.expected {
				t.Fatalf("normalizePhone(%q) = %q, want %q", tt.input, actual, tt.expected)
			}
		} else {
			if !errors.Is(err, tt.err) {
				t.Fatalf("normalizePhone(%q) error %v, want %v", tt.input, err, tt.err)
			}
		}
	}
}

func ExampleLookup() {
	entry, err := Lookup("+1-303-123-4567")
	if err != nil {
		panic(err)
	}

	fmt.Printf("%s, %s (%s) — %s", entry.City, entry.State, entry.StateCode, entry.Type)
	// Output: Aurora, Colorado (CO) — local
}

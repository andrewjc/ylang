package common

import "testing"

func TestIsDigitUnit(t *testing.T) {
	testCases := []struct {
		name     string
		ch       rune
		expected bool
	}{
		{"ASCII Digit 0", '0', true},
		{"ASCII Digit 5", '5', true},
		{"ASCII Digit 9", '9', true},
		{"ASCII Letter a", 'a', false},
		{"ASCII Letter Z", 'Z', false},
		{"Symbol %", '%', false},
		{"Symbol -", '-', false},
		{"Space", ' ', false},
		{"Newline", '\n', false},
		{"Tab", '\t', false},
		{"Underscore", '_', false},
		{"Unicode Digit ١", '١', true},   // Arabic-Indic Digit One
		{"Unicode Digit ০", '০', true},   // Bengali Digit Zero
		{"Unicode Letter α", 'α', false}, // Greek Alpha
		{"Zero Rune", 0, false},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if got := IsDigit(tc.ch); got != tc.expected {
				t.Errorf("IsDigit('%c') = %v, want %v", tc.ch, got, tc.expected)
			}
		})
	}
}

func TestIsLetterUnit(t *testing.T) {
	testCases := []struct {
		name     string
		ch       rune
		expected bool
	}{
		{"ASCII Lower a", 'a', true},
		{"ASCII Lower z", 'z', true},
		{"ASCII Upper A", 'A', true},
		{"ASCII Upper Z", 'Z', true},
		{"Underscore", '_', true},
		{"ASCII Digit 0", '0', false},
		{"ASCII Digit 9", '9', false},
		{"Symbol %", '%', false},
		{"Symbol -", '-', false},
		{"Space", ' ', false},
		{"Newline", '\n', false},
		{"Tab", '\t', false},
		{"Unicode Letter α", 'α', true},    // Greek Alpha (Lowercase)
		{"Unicode Letter Ω", 'Ω', true},    // Greek Omega (Uppercase)
		{"Unicode Letter é", 'é', true},    // Latin Small Letter E with Acute
		{"Unicode Letter 国", '国', true},    // CJK Unified Ideograph (Han character)
		{"Unicode Digit ١", '١', false},    // Arabic-Indic Digit One
		{"Combining Mark ´", '´', false},   // Acute Accent (Spacing Modifier != Letter)
		{"Punctuation Mark ?", '?', false}, // Punctuation
		{"Zero Rune", 0, false},
		// Boundary cases - different Unicode categories that might be adjacent
		{"Modifier Letter ʻ", 'ʻ', true},   // Modifier Letter Turned Comma (Category Lm)
		{"Letter Number  रोमन", '௰', true}, // Tamil Number Ten (Category Nl) - Unicode considers letter numbers as letters
		{"Other Symbol  Ohm Ω", 'Ω', true}, // Ohm Sign (Category So - Other Symbol, but often used like a letter) - unicode.IsLetter considers it true
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			if got := IsLetter(tc.ch); got != tc.expected {
				t.Errorf("IsLetter('%c') = %v, want %v", tc.ch, got, tc.expected)
			}
		})
	}
}

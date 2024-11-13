package cmd

import (
	crand "crypto/rand"
	"errors"
	"fmt"
	"math/big"
	"math/rand"
	"os"
	"strings"
	"unicode"

	"github.com/spf13/cobra"
)

var (
	length     uint8
	digits     bool
	symbols    bool
	minDigits  uint8
	minSymbols uint8
	ambiguous  bool
)

//goland:noinspection ALL
var (
	upperCharsAll    = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	upperCharsNoAmb  = "ABCDEFGHJKLMNPQRSTUVWXYZ" // exclude I O
	lowerCharsAll    = "abcdefghijklmnopqrstuvwxyz"
	lowerCharsNoAmb  = "abcdefghijkmnpqrstuvwxyz" // exclude l o
	digitsCharsAll   = "0123456789"
	digitsCharsNoAmb = "23456789" // exclude 0 1
	symbolsChars     = "!@#$%^&*"
)

var rootCmd = &cobra.Command{
	Use:   "pass-gen",
	Short: "A random password generator",
	Long: `A powerful random password generator that creates secure passwords based on specified rules.
Controls password length, character types, and minimum digits of digits and special characters.`,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		return validateFlags()
	},
	Run: func(cmd *cobra.Command, args []string) {
		upperChars := upperCharsNoAmb
		lowerChars := lowerCharsNoAmb
		digitsChars := digitsCharsNoAmb
		if ambiguous {
			upperChars = upperCharsAll
			lowerChars = lowerCharsAll
			digitsChars = digitsCharsAll
		}

		var password strings.Builder
		var chars string

		if digits {
			var init uint8 = 0
			for ; init < minDigits; init++ {
				password.WriteByte(digitsChars[secureRandomInt(len(digitsChars))])
			}
		}
		if symbols {
			var init uint8 = 0
			for ; init < minSymbols; init++ {
				password.WriteByte(symbolsChars[secureRandomInt(len(symbolsChars))])
			}
		}

		chars += upperChars
		chars += lowerChars

		if digits {
			chars += digitsChars
		}
		if symbols {
			chars += symbolsChars
		}

		var init uint8 = 0
		for ; init < length-(minDigits+minSymbols); init++ {
			password.WriteByte(chars[secureRandomInt(len(chars))])
		}

		pwdRunes := []rune(password.String())
		for {
			rand.Shuffle(
				len(pwdRunes), func(i, j int) {
					pwdRunes[i], pwdRunes[j] = pwdRunes[j], pwdRunes[i]
				},
			)
			if !strings.Contains(symbolsChars, string(pwdRunes[0])) && !unicode.IsDigit(pwdRunes[0]) {
				break
			}
		}
		fmt.Println(string(pwdRunes))
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.Flags().Uint8VarP(&length, "length", "l", 16, "password length")
	rootCmd.Flags().BoolVarP(&digits, "digits", "d", false, "include digits")
	rootCmd.Flags().BoolVarP(&symbols, "symbols", "s", false, "include symbols")
	rootCmd.Flags().Uint8VarP(&minDigits, "min-digits", "D", 3, "minimum digits of digits")
	rootCmd.Flags().Uint8VarP(&minSymbols, "min-symbols", "S", 2, "minimum digits of symbols")
	rootCmd.Flags().BoolVarP(&ambiguous, "ambiguous", "A", false, "include ambiguous characters")
}

func secureRandomInt(max int) int {
	num, err := crand.Int(crand.Reader, big.NewInt(int64(max)))
	if err != nil {
		panic(err)
	}
	return int(num.Int64())
}

func validateFlags() error {
	if length < 8 {
		return errors.New("password length must be at least 8 characters")
	}

	if minDigits+minSymbols > length/2 {
		return errors.New("minimum digits and symbols must be less than half of password length")
	}
	return nil
}

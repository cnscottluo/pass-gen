package cmd

import (
	crand "crypto/rand"
	"fmt"
	"math/big"
	"math/rand"
	"os"
	"strings"
	"unicode"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "pass-gen",
	Short: "A random password generator",
	Long: `A powerful random password generator that creates secure passwords based on specified rules.
Controls password length, character types, and minimum number of digits and special characters.`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := validateFlags(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}

		var (
			capitalCharsAll   = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
			capitalCharsNoAmb = "ABCDEFGHJKLMNPQRSTUVWXYZ" // 排除了 I 和 O
			smallCharsAll     = "abcdefghijklmnopqrstuvwxyz"
			smallCharsNoAmb   = "abcdefghijkmnpqrstuvwxyz" // 排除了 l 和 o
			numberCharsAll    = "0123456789"
			numberCharsNoAmb  = "23456789" // 排除了 0 和 1
			symbolChars       = "!@#$%^&*"
		)

		capitalChars := capitalCharsAll
		smallChars := smallCharsAll
		numberChars := numberCharsAll
		if avoidAmbiguous {
			capitalChars = capitalCharsNoAmb
			smallChars = smallCharsNoAmb
			numberChars = numberCharsNoAmb
		}

		var password strings.Builder
		var chars string

		if number {
			for i := 0; i < minNumber; i++ {
				password.WriteByte(numberChars[secureRandomInt(len(numberChars))])
			}
		}
		if symbol {
			for i := 0; i < minSymbol; i++ {
				password.WriteByte(symbolChars[secureRandomInt(len(symbolChars))])
			}
		}

		if capital {
			chars += capitalChars
		}
		if small {
			chars += smallChars
		}
		if number {
			chars += numberChars
		}
		if symbol {
			chars += symbolChars
		}

		for i := 0; i < length-(minNumber+minSymbol); i++ {
			password.WriteByte(chars[secureRandomInt(len(chars))])
		}

		pwdRunes := []rune(password.String())
		for {
			rand.Shuffle(len(pwdRunes), func(i, j int) {
				pwdRunes[i], pwdRunes[j] = pwdRunes[j], pwdRunes[i]
			})
			if !strings.Contains(symbolChars, string(pwdRunes[0])) && !unicode.IsDigit(pwdRunes[0]) {
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

var (
	// Password length
	length int
	// Include uppercase letters
	capital bool
	// Include lowercase letters
	small bool
	// Include numbers
	number bool
	// Include special characters
	symbol bool
	// Minimum number of digits
	minNumber int
	// Minimum number of special characters
	minSymbol int
	// Avoid ambiguous characters
	avoidAmbiguous bool
)

func init() {
	rootCmd.Flags().IntVar(&length, "length", 16, "Password length")
	rootCmd.Flags().BoolVar(&capital, "capital", true, "Include uppercase letters")
	rootCmd.Flags().BoolVar(&small, "small", true, "Include lowercase letters")
	rootCmd.Flags().BoolVar(&number, "number", true, "Include numbers")
	rootCmd.Flags().BoolVar(&symbol, "symbol", true, "Include special characters")
	rootCmd.Flags().IntVar(&minNumber, "min-number", 1, "Minimum number of digits")
	rootCmd.Flags().IntVar(&minSymbol, "min-symbol", 1, "Minimum number of special characters")
	rootCmd.Flags().BoolVar(&avoidAmbiguous, "avoid-ambiguous", true, "Avoid ambiguous characters")
}

func secureRandomInt(max int) int {
	num, err := crand.Int(crand.Reader, big.NewInt(int64(max)))
	if err != nil {
		panic(err)
	}
	return int(num.Int64())
}

func validateFlags() error {
	if length < 1 {
		return fmt.Errorf("password length must be positive")
	}
	if !capital && !small && !number && !symbol {
		return fmt.Errorf("at least one character type must be selected")
	}
	if minNumber+minSymbol > length {
		return fmt.Errorf("minimum requirements exceed password length")
	}
	return nil
}

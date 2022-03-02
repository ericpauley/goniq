package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"
	"unicode"

	"github.com/pborman/getopt/v2"
)

var countFlag = getopt.BoolLong("count", 'c', "prefix lines with occurence count")
var uniqueFlag = getopt.BoolLong("unique", 'u', "only print uniques")
var duplicateFlag = getopt.BoolLong("repeated", 'd', "only print duplicates")
var ignoreCaseFlag = getopt.BoolLong("ignore-case", 'i', "ignore case")
var skipFieldsFlag = getopt.Uint64Long("skip-fields", 'f', 0, "skip n fields")
var skipCharsFlag = getopt.Uint64Long("skip-chars", 's', 0, "skip characters (done after skipping fields)")
var checkCharsFlag = getopt.Uint64Long("check-chars", 'w', 0, "compare max N characters per line")
var helpFlag = getopt.BoolLong("help", 'h', "show help")

func canonicalize(in string) (result string) {
	result = in
	skipFields := *skipFieldsFlag

	if skipFields != 0 {
		skippingSpaces := true
		for i, r := range result {
			if skippingSpaces && !unicode.IsSpace(r) {
				// We found a non-space character, now consume non-spaces
				skippingSpaces = false
			} else if !skippingSpaces && unicode.IsSpace(r) {
				// We reached the end of a field, switch to consuming spaces or break if we're done skipping fields
				skipFields--
				skippingSpaces = true
				if skipFields == 0 {
					result = result[i:]
					break
				}

			}
		}
	}

	skipChars := *skipCharsFlag

	for i := range result {
		if skipChars == 0 {
			result = result[i:]
			break
		}
		skipChars--
	}

	checkChars := *checkCharsFlag
	if checkChars != 0 {
		for i := range result {
			if checkChars == 0 {
				result = result[:i]
				break
			}
			checkChars--
		}
	}

	if *ignoreCaseFlag {
		result = strings.ToLower(result)
	}

	return
}

func main() {
	getopt.CommandLine.SetParameters("")
	getopt.Parse()
	if *helpFlag {
		fmt.Println("goniq is like uniq but does not require lines be consecutive.")
		fmt.Println("Input/output is always stdin/stdout.")
		getopt.PrintUsage(os.Stdout)
		os.Exit(0)
	}
	if *uniqueFlag && *duplicateFlag {
		log.Fatal("Cannot use -u and -d together.")
	}
	// We can't stream if we need to keep counts or only output uniques.
	stream := !*countFlag && !*uniqueFlag
	counts := make(map[string]int)
	// Canonical holds the non-lowercase string from the first time we've seen it. Only used with -i
	canonical := make(map[string]string)

	scanner := bufio.NewScanner(os.Stdin)
	// We need to track the order if we aren't streaming the output.
	var orderedKeys []string
	for scanner.Scan() {
		// Orig always holds the actual seen string so we can output it
		orig := scanner.Text()
		text := canonicalize(orig)
		counts[text] += 1
		if stream {
			if *duplicateFlag {
				// Either output the original (first-seen) copy of the string or set that if it's the first time seen.
				if counts[text] == 2 {
					fmt.Println(canonical[text])
					delete(canonical, text)
				} else if counts[text] == 1 {
					canonical[text] = orig
				}
			} else if counts[text] == 1 {
				fmt.Println(orig)
			}
		} else {
			// Track the order we saw keys in. orderedKeys are always canonical (first-seen) strings.
			if counts[text] == 1 {
				orderedKeys = append(orderedKeys, orig)
			}
		}

	}
	if scanner.Err() != nil {
		log.Fatal("Failed to read input:", scanner.Err())
	}
	if !stream {
		for i := 0; i < len(orderedKeys); i++ {
			orig := orderedKeys[i]
			text := canonicalize(orig)
			if *uniqueFlag && counts[text] != 1 {
				continue
			}
			if *duplicateFlag && counts[text] == 1 {
				continue
			}
			// Print with or without counts, always using the original first-seen string.
			if *countFlag {
				fmt.Println(counts[text], orig)
			} else {
				fmt.Println(orig)
			}
		}
	}
}

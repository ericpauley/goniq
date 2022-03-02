package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"strings"

	"github.com/pborman/getopt/v2"
)

var countFlag = getopt.BoolLong("count", 'c', "prefix lines with occurence count")
var uniqueFlag = getopt.BoolLong("unique", 'u', "only print uniques")
var duplicateFlag = getopt.BoolLong("repeated", 'd', "only print duplicates")
var ignoreCase = getopt.BoolLong("ignore-case", 'i', "ignore case")
var helpFlag = getopt.BoolLong("help", 'h', "show help")

func main() {
	getopt.Parse()
	if *helpFlag {
		fmt.Println("goniq is like uniq but does not require lines be consecutive.")
		fmt.Println("Input/output is always stdin/stdout.")
		getopt.PrintUsage(os.Stdout)
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
		text := orig
		if *ignoreCase {
			text = strings.ToLower(text)
		}
		counts[text] += 1
		if stream {
			if *duplicateFlag {
				// Either output the canonical (first-seen) copy of the string or set that if it's the first time seen.
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
	if scanner.Err() != io.EOF {
		log.Fatal("Failed to read input:", scanner.Err())
	}
	if !stream {
		for i := 0; i < len(orderedKeys); i++ {
			orig := orderedKeys[i]
			text := orig
			if *ignoreCase {
				text = strings.ToLower(text)
			}
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

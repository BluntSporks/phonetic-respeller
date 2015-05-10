// Respell a text file phonetically by looking up available words.
package main

import (
	"bufio"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/BluntSporks/cmudict"
	"github.com/BluntSporks/natural-language"
)

func main() {
	// Parse flags.
	dictFile := flag.String("dict", cmudict.DefaultDictPath(), "Name of CMU-formatted file to modify")
	textFile := flag.String("text", "", "Name of text file to respell")
	flag.Parse()

	if len(*textFile) == 0 {
		log.Fatal("Missing -text argument")
	}

	// Load CMUDict file.
	cmuDict := cmudict.LoadDict(*dictFile)

	// Open dict file.
	hdl, err := os.Open(*textFile)
	if err != nil {
		log.Fatal(err)
	}
	defer hdl.Close()

	// Scan file line by line.
	scanner := bufio.NewScanner(hdl)
	for scanner.Scan() {
		line := scanner.Text()
		if err := scanner.Err(); err != nil {
			log.Fatal(err)
		}
		words := natlang.MatchAllWordsPlus(line)
		for _, word := range words {
			upword := strings.ToUpper(word)
			isUpper := word != strings.ToLower(word)
			pron := cmuDict[upword]
			if pron == "" {
				fmt.Print(word)
			} else {
				phonemes := cmudict.GetPhonemes(pron, true)
				fixed := cmudict.FixPhonemes(phonemes)
				for i, phoneme := range fixed {
					bare := cmudict.StripAccent(phoneme)
					spelling := Spellings[bare]
					out := phoneme
					if spelling != "" {
						out = spelling
					}

					// Capitalize the first letter of output if the original word was not
					// lowercase.
					if i == 0 && isUpper {
						cap := strings.ToUpper(out[0:1])
						if len(out) > 1 {
							cap += out[1:]
						}
						out = cap
					}
					fmt.Print(out)
				}
			}
		}
		fmt.Println()
	}
}

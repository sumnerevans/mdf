package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/araddon/dateparse"
	"github.com/fatih/color"
)

// Git diff regexes
var gitSummaryFileRe = regexp.MustCompile(`( .* \|\s+)(\d+) ([\+-]+).*`)
var gitSummaryRe = regexp.MustCompile(` (\d+) (files?) changed(?:, (\d+) (insertions?)\(\+\))(?:, (\d+) (deletions?)\(-\))`)
var diffRe = regexp.MustCompile(`diff --git`)
var diffDescRe = regexp.MustCompile(`@@ (-\d+,\d+ \+\d+,\d+) @@`)

var bold = color.New(color.Bold, color.FgWhite)
var boldGreen = color.New(color.Bold, color.FgGreen)
var boldRed = color.New(color.Bold, color.FgRed)

// Regular email regular expressions
var dateMatchRe = regexp.MustCompile(`Date: (.*)`)
var emailRe = regexp.MustCompile(`<mailto:(.*?)>`)
var urlRe = regexp.MustCompile(`https?:\/\/[^\s\)\]>]*`)
var tzNameRe = regexp.MustCompile(`\/.*\/([^/]+\/[^/]+)`)

func RunFilter(rootUri string) {
	scanner := bufio.NewScanner(os.Stdin)

	// Force color output
	color.NoColor = false

	// Git Diff state
	hitDiff := false
	hitDiffFooter := false

	for scanner.Scan() {
		line := scanner.Text()
		if diffRe.MatchString(line) {
			hitDiff = true
			bold.Println(line)
		} else if match := gitSummaryRe.FindStringSubmatch(line); match != nil {
			bold.Printf(" %s %s changed", match[1], match[2])
			if match[3] != "" && match[4] != "" {
				bold.Print(", ")
				boldGreen.Printf("%s %s(+)", match[3], match[4])
			}
			if match[5] != "" && match[6] != "" {
				bold.Print(", ")
				boldRed.Printf("%s %s(-)", match[5], match[6])
			}
			fmt.Println()
		} else if match := gitSummaryFileRe.FindStringSubmatch(line); match != nil {
			fmt.Print(match[1])
			bold.Printf("%s ", match[2])
			pluses := 0
			minuses := 0
			for _, c := range match[3] {
				if c == '+' {
					pluses++
				} else if c == '-' {
					minuses++
				}
			}
			boldGreen.Print(strings.Repeat("+", pluses))
			boldRed.Print(strings.Repeat("-", minuses))
			fmt.Println()
		} else if hitDiffFooter {
			color.Yellow(line)
		} else if hitDiff {
			if strings.HasPrefix(line, "--") && !strings.HasPrefix(line, "---") {
				hitDiffFooter = true
				color.Yellow(line)
			} else if strings.HasPrefix(line, "+") {
				color.Green(line)
			} else if strings.HasPrefix(line, "-") {
				color.Red(line)
			} else if match := diffDescRe.FindStringSubmatch(line); match != nil {
				color.Blue(line)
			} else if strings.HasPrefix(line, " ") {
				fmt.Println(line)
			} else {
				bold.Println(line)
			}
		} else if match := dateMatchRe.FindStringSubmatch(line); match != nil {
			parsed, err := dateparse.ParseAny(match[1])
			if err != nil {
				fmt.Printf("%s [couldn't parse date]", line)
			} else {
				tzPath, err := os.Readlink("/etc/localtime")
				if err != nil {
					fmt.Printf("%s [couldn't get /etc/localtime]", line)
				} else {
					tzName := tzNameRe.FindStringSubmatch(tzPath)
					if tzName == nil {
						fmt.Printf("%s [couldn't parse tz path %s]", line, tzPath)
					} else {
						loc, err := time.LoadLocation(tzName[1])
						if err != nil {
							fmt.Printf("%s [couldn't load location %s]", line, tzName)
						} else {
							fmt.Printf("Date: %s", parsed.In(loc).Format("Mon, 02 Jan 2006 15:04:05 MST (-07:00)"))
						}
					}
				}
			}
			fmt.Println()
		} else {
			for _, match := range emailRe.FindAllStringSubmatch(line, -1) {
				line = strings.ReplaceAll(line, fmt.Sprintf("%s<mailto:%s>", match[1], match[1]), match[1])
			}
			for _, match := range urlRe.FindAllStringSubmatch(line, -1) {
				if len(rootUri)+6 < len(match[0]) {
					resp, err := http.Post(rootUri+"new", "text/text", strings.NewReader(match[0]))
					if err == nil {
						id, err := ioutil.ReadAll(resp.Body)
						if err == nil {
							line = strings.ReplaceAll(line, match[0], rootUri+string(id))
						}
					}
				}
			}
			fmt.Println(line)
		}
	}
}

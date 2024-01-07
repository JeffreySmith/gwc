package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"strings"
	"unicode/utf8"
)

type Line struct {
	text                string
	words, bytes, chars int
}

func initLine(text string) Line {
	line := Line{text: text}
	line.words = len(strings.Fields(line.text))
	line.bytes = len(line.text)
	line.chars = utf8.RuneCountInString(line.text)
	return line
}
func parseLines(Lines []Line) (int, int, int) {

	wordCount := 0
	byteCount := 0
	charCount := 0

	for _, l := range Lines {
		wordCount += l.words
		byteCount += l.bytes
		charCount += l.chars
	}

	return wordCount, byteCount, charCount
}
func parseInput(src io.Reader) ([]Line, int, int) {
	var Lines []Line

	lineCount := 0
	longest := 0

	scanner := bufio.NewScanner(src)

	line := ""

	scanner.Split(bufio.ScanRunes)
	for scanner.Scan() {
		if scanner.Text() == "\n" {
			if len(strings.Trim(line, "\n")) > longest {
				longest = len(strings.Trim(line, "\n"))
			}
			newLine := initLine(line)
			Lines = append(Lines, newLine)
			line = ""
			lineCount++
		}
		line += scanner.Text()
	}
	//Deal with no new line at the end of input
	if len(Lines) == 0 || len(line) > 0 {
		if len(line) > longest {
			longest = len(line)
		}
		newLine := initLine(line)
		Lines = append(Lines, newLine)
	}
	return Lines, lineCount, longest
}

func main() {
	var Lines []Line
	var count int
	var longest int
	var w, b, c int
	word := flag.Bool("w", false, "Display the number of words")
	bytes := flag.Bool("c", false, "Display the number of bytes")
	lineCount := flag.Bool("l", false, "Display the number of lines")
	chars := flag.Bool("m", false, "Display the number of characters. UTF-8 aware")

	flag.Parse()

	fileNames := flag.Args()

	info, err := os.Stdin.Stat()
	if err != nil {
		panic(err)
	}

	//Per the manpages, -m supercedes -c
	if *bytes && *chars {
		*bytes = false
	}

	if !*word && !*lineCount && (!*bytes && !*chars) {
		fmt.Println("No options selected")

	}

	if len(fileNames) == 0 && info.Mode()&os.ModeCharDevice == 0 {
		//Process stdin stuff
		Lines, count, longest = parseInput(os.Stdin)

	} else if len(fileNames) == 0 {
		//Run until EOF hit ()
	} else {
		//Process all files passed in
	}

	//Process output here

	if len(fileNames) > 0 {

	} else {
		w, b, c = parseLines(Lines)
		fmt.Println(w, b, c)
	}

	fmt.Println(count, longest)
	fmt.Println(*word, *bytes, *lineCount, *chars)
}

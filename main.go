package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf8"
)

type Line struct {
	text                string
	words, bytes, chars int
}

type FileInput struct {
	lines   []Line
	name    string
	count   int
	longest int
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
		/*if err := scanner.Err(); err != nil {
			if err != io.EOF {
				panic(err)
			}
			break
		}*/
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
	var Files []FileInput
	var count int

	var w, b, c int
	longest := 0
	word := flag.Bool("w", false, "Display the number of words")
	byteOut := flag.Bool("c", false, "Display the number of bytes")
	lines := flag.Bool("l", false, "Display the number of lines")
	chars := flag.Bool("m", false, "Display the number of characters. UTF-8 aware")
	longestFlag := flag.Bool("L", false, "Display the number of bytes or characters (if -m is passed). If more than one file is provided, the longest line is reported in 'total'")
	flag.Parse()

	fileNames := flag.Args()

	info, err := os.Stdin.Stat()
	if err != nil {
		panic(err)
	}

	//Per the manpages, -m supercedes -c
	if *byteOut && *chars {
		*byteOut = false
	}

	if !*word && !*lines && (!*byteOut && !*chars) {
		*word = true
		*lines = true
		*byteOut = true
	}

	if len(fileNames) == 0 && info.Mode()&os.ModeCharDevice == 0 {
		//Process stdin stuff
		Lines, count, longest = parseInput(os.Stdin)

	} else if len(fileNames) == 0 {
		//Run until EOF hit
		reader := bufio.NewReader(os.Stdin)
		var byteArray []byte
		for {
			byte, err := reader.ReadByte()
			if err != nil {
				if err != io.EOF {
					panic(err)
				}
				break
			}
			byteArray = append(byteArray, byte)
		}
		inputReader := bufio.NewReader(bytes.NewBuffer(byteArray))
		Lines, count, longest = parseInput(inputReader)

	} else {
		//Process all files passed in
		for _, file := range fileNames {
			f, err := os.Open(file)
			defer f.Close()
			if err != nil {
				fmt.Printf("%s: %v: open: No such file or directory\n", filepath.Base(os.Args[0]), file)
			}

			Lines, count, longest = parseInput(f)
			newFile := FileInput{lines: Lines, count: count, name: file, longest: longest}
			Files = append(Files, newFile)

		}
	}

	//Process output here

	if len(fileNames) > 0 {
		totalLines := 0
		totalWords := 0
		totalBytes := 0
		totalChars := 0
		longest := 0
		for _, f := range Files {
			w, b, c = parseLines(f.lines)

			totalLines += f.count
			totalWords += w
			totalBytes += b
			totalChars += c

			if f.longest > longest {
				longest = f.longest
			}

			if *lines {
				fmt.Printf("%8v", f.count)
			}
			if *word {
				fmt.Printf("%8v", w)
			}
			if *byteOut {
				fmt.Printf("%8v", b)
			}
			if *chars {
				fmt.Printf("%8v", c)
			}
			if *longestFlag {
				fmt.Printf("%8v", f.longest)
			}
			fmt.Printf(" %v\n", f.name)
		}
		if len(fileNames) > 1 {
			if *lines {
				fmt.Printf("%8v", totalLines)
			}
			if *word {
				fmt.Printf("%8v", totalWords)
			}
			if *byteOut {
				fmt.Printf("%8v", totalBytes)
			}
			if *chars {
				fmt.Printf("%8v", totalChars)
			}
			if *longestFlag {
				fmt.Printf("%8v", longest)
			}
			fmt.Printf(" total\n")
		}
	} else {
		w, b, c = parseLines(Lines)

		if *lines {
			fmt.Printf("%8v", count)
		}
		if *word {
			fmt.Printf("%8v", w)
		}
		if *byteOut {
			fmt.Printf("%8v", b)
		}
		if *chars {
			fmt.Printf("%8v", c)
		}
		fmt.Println()
	}
}

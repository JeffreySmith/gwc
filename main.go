package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
	"unicode/utf8"

	"github.com/jessevdk/go-flags"
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

	var opts struct {
		WordCount   bool `short:"w" description:"Display the number of words"`
		LineCount   bool `short:"l" description:"Display the number of lines"`
		ByteCount   bool `short:"c" description:"Display the number of bytes. Supercedes -m"`
		CharCount   bool `short:"m" description:"Display the number of characters. UTF-8 aware"`
		LongestLine bool `short:"L" description:"Display the line with the most bytes or characters (with -m). When there is more than one file, the longest line will be shown as the value for 'total'."`
	}

	var Lines []Line
	var Files []FileInput
	var count int

	var w, b, c int
	longest := 0

	fileNames, err := flags.Parse(&opts)
	
	if err!=nil{
		flagsErr, ok := err.(*flags.Error)
		if ok{
			if flagsErr.Type == flags.ErrHelp {
				os.Exit(0)
				return
			}
			if flagsErr.Type == flags.ErrUnknownFlag{
				fmt.Printf("Use -h or --help for usage information\n")
				fmt.Printf("Usage: gwc [-Lclmw] [file ...]\n")

			}
		}
		os.Exit(1)
	}
	info, err := os.Stdin.Stat()
	if err != nil {

		os.Exit(1)
	}

	//Per the manpages, -c supercedes -m
	if opts.ByteCount && opts.CharCount {
		opts.CharCount = false
	}
	//Default, with no flags provided is to show the line count, word count, and the byte count
	if !opts.WordCount && !opts.LineCount && (!opts.ByteCount && !opts.CharCount) {
		opts.WordCount = true
		opts.LineCount = true
		opts.ByteCount = true
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
					os.Exit(1)
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

			if opts.LineCount {
				fmt.Printf("%8v", f.count)
			}
			if opts.WordCount {
				fmt.Printf("%8v", w)
			}
			if opts.ByteCount {
				fmt.Printf("%8v", b)
			}
			if opts.CharCount {
				fmt.Printf("%8v", c)
			}
			if opts.LongestLine {
				fmt.Printf("%8v", f.longest)
			}
			fmt.Printf(" %v\n", f.name)
		}
		if len(fileNames) > 1 {
			if opts.LineCount {
				fmt.Printf("%8v", totalLines)
			}
			if opts.WordCount {
				fmt.Printf("%8v", totalWords)
			}
			if opts.ByteCount {
				fmt.Printf("%8v", totalBytes)
			}
			if opts.CharCount {
				fmt.Printf("%8v", totalChars)
			}
			if opts.LongestLine {
				fmt.Printf("%8v", longest)
			}
			fmt.Printf(" total\n")
		}
	} else {
		w, b, c = parseLines(Lines)

		if opts.LineCount {
			fmt.Printf("%8v", count)
		}
		if opts.WordCount {
			fmt.Printf("%8v", w)
		}
		if opts.ByteCount {
			fmt.Printf("%8v", b)
		}
		if opts.CharCount {
			fmt.Printf("%8v", c)
		}
		fmt.Println()
	}
}

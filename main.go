package main

import (
	"flag"
	"fmt"
	"os"

)

func main(){
	word := flag.Bool("w",false,"Display the number of words")
	bytes := flag.Bool("c",false,"Display the number of bytes")
	lineCount := flag.Bool("l",false,"Display the number of lines")
	chars := flag.Bool("m",false,"Display the number of characters. UTF-8 aware")

	flag.Parse()

	fileNames := flag.Args()

	
	info,err := os.Stdin.Stat()
	if err != nil{
		panic(err)
	}

	//Per the manpages, -m supercedes -c
	if *bytes && *chars {
		*bytes = false
	}
	
	if len(fileNames) == 0 && info.Mode() & os.ModeCharDevice == 0{
		//Process stdin stuff
	} else if len(fileNames) == 0 {
		//Run until EOF hit
	} else {
		//Process all files passed in
	}
	
	fmt.Println(*word,*bytes,*lineCount,*chars)
}

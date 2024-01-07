# gwc
A golang version of the unix tool wc

This entirely go version of wc vaguely follows the options for the FreeBSD & MacOS version of wc. It does not, however, implement the --libxo option.

Usage:
```wc -Llwcm [file ...]```

-L : Writes the length of the line with the most bytes (default), or the one with the most chars (-m). If more than one file is provided, the longest line in all of the files will be provided in the final total.

-c : Prints out the number of bytes for each input to stdout. Cancels out the -m flag.

-m : Prints out the number of multibyte characters to stdout.

-w : Prints out the number of words to stdout


This is written entirely in go using only one library from outside the standard library: https://github.com/jessevdk/go-flags. This was because the default flag parser in go uses plan-9 style flag arguments which doesn't allow you to combine multiple options together (ie gwc -wcl). My aim was to make this as similar to the real tool as possible.

The biggest thing I learned here is how non-standard this standard unix tool is accross platforms. Each version across FreeBSD, Linux, NetBSD, and OpenBSD, all have different options or different ways of interpreting certain features (does -c or -m take precedence? This isn't the same across each platform). 

# stegify
[![GoDoc](https://godoc.org/github.com/DimitarPetrov/stegify?status.svg)](https://godoc.org/github.com/DimitarPetrov/stegify)
[![Go Report Card](https://goreportcard.com/badge/github.com/DimitarPetrov/stegify)](https://goreportcard.com/report/github.com/DimitarPetrov/stegify)

## Overview
`stegify` is a simple command line tool that is capable of fully transparent hiding files within an image.
The technique is known as [steganography](https://en.wikipedia.org/wiki/steganography). This tool uses
the technique known as LSB (Least Significant Bit) Steganography. 

## Install
```
go get github.com/DimitarPetrov/stegify
```

## Usage

```
stegify -op encode -carrier <file-name> -data <file-name> -result <file-name>
stegify -op decode -carrier <file-name> -result <file-name>
```
When encoding, the file with name given to flag `-data` is hidden inside file with name given to flag
`-carrier` and the resulting file is saved new file in the current working directory under the
name given to flag `-result`.

When decoding, given a file name of a carrier file with previously encoded data in it, the data is extracted
and saved in new file under the name given to flag `-result` in current working directory.

In both cases the flag `-result` could be omitted and it will be used the default file name: `result`
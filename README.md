# stegify
[![Build Status](https://travis-ci.org/DimitarPetrov/stegify.svg?branch=master)](https://travis-ci.org/DimitarPetrov/stegify)
[![Coverage Status](https://coveralls.io/repos/github/DimitarPetrov/stegify/badge.svg?branch=master)](https://coveralls.io/github/DimitarPetrov/stegify?branch=master)
[![GoDoc](https://godoc.org/github.com/DimitarPetrov/stegify?status.svg)](https://godoc.org/github.com/DimitarPetrov/stegify)
[![Go Report Card](https://goreportcard.com/badge/github.com/DimitarPetrov/stegify)](https://goreportcard.com/report/github.com/DimitarPetrov/stegify)
[![Mentioned in Awesome Go](https://awesome.re/mentioned-badge.svg)](https://github.com/avelino/awesome-go)  


## Overview
`stegify` is a simple command line tool capable of fully transparent hiding any file within an image or set of images.
This technique is known as LSB (Least Significant Bit) [steganography](https://en.wikipedia.org/wiki/steganography) 

## Demonstration

| Carrier                                | Data                                | Result                                               |
| ---------------------------------------| ------------------------------------|------------------------------------------------------|
| ![Original File](examples/street.jpeg) | ![Data file](examples/lake.jpeg)    | ![Encoded File](examples/test_decode.jpeg)           |

The `Result` file contains the `Data` file hidden in it. And as you can see it is fully transparent.

If multiple `Carrier` files are provided, the `Data` file will be split in pieces and every piece is encoded in the respective carrier.

| Carrier1                                     | Carrier2                                   | Data                                       | Result1                                                          | Result2                                                          |
| ---------------------------------------------|--------------------------------------------|--------------------------------------------|------------------------------------------------------------------|------------------------------------------------------------------|
| <img src="examples/street.jpeg" width="500"> | <img src="examples/lake.jpeg" width="500"> | <img src="examples/video.gif" width="500"> | <img src="examples/test_multi_carrier_decode1.jpeg" width="500"> | <img src="examples/test_multi_carrier_decode2.jpeg" width="500"> |
 
The `Result1` file contains one half of the `Data` file hidden in it and `Result2` the other. As always fully transparent.

## Installation

#### Installing from Source
```
go get -u github.com/DimitarPetrov/stegify
```

#### Installing via Homebrew (macOS)
```
brew tap DimitarPetrov/stegify
brew install stegify
```

Or you can download a binary for your system [here](https://github.com/DimitarPetrov/stegify/releases).

## Usage

### As a command line tool

#### Single carrier encoding/decoding
```
stegify encode --carrier <file-name> --data <file-name> --result <file-name>

stegify decode --carrier <file-name> --result <file-name>
```
When encoding, the file with name given to flag `--data` is hidden inside the file with name given to flag
`--carrier` and the resulting file is saved in new file in the current working directory under the
name given to flag `--result`.

> **_NOTE:_** The result file won't have any file extension and therefore it should be specified explicitly in `--result` flag. 

When decoding, given a file name of a carrier file with previously encoded data in it, the data is extracted
and saved in new file in the current working directory under the name given to flag `--result`.

> **_NOTE:_** The result file won't have any file extension and therefore it should be specified explicitly in `--result` flag.

In both cases the flag `--result` could be omitted and default values will be used.

#### Multiple carriers encoding/decoding

```
stegify encode --carriers "<file-names...>" --data <file-name> --results "<file-names...>"
OR
stegify encode --carrier <file-name> --carrier <file-name> ... --data <file-name> --result <file-name> --result <file-name> ...

stegify decode --carriers "<file-names...>" --result <file-name>
OR
stegify decode --carrier <file-name> --carrier <file-name> ... --result <file-name>
```
When encoding a data file in more than one carriers, the data file is split in *N* chunks, where *N* is number of provided carriers.
Each of the chunks is then encoded in the respective carrier.

> **_NOTE:_** When decoding, carriers should be provided in the **exact** same order for result to be properly extracted. 

This kind of encoding provides one more layer of security and more flexibility regarding size limitations.

In both cases the flag `--result/--results` could be omitted and default values will be used.

> **_NOTE:_** When encoding the number of the result files (if provided) should be equal to the number of carrier files. When decoding, exactly one result is expected. 

When multiple carriers are provided with mixed kinds of flags, the names provided through `carrier` flag are taken first and with `carriers/c` flags second.
Same goes for the `result/results` flag.


### Programmatically in your code

`stegify` can be used programmatically too and it provides easy to use functions working with file names
or raw Readers and Writers. You can visit [godoc](https://godoc.org/github.com/DimitarPetrov/stegify) under
`steg` package for details.

## Disclaimer

If carrier file is in jpeg or jpg format, after encoding the result file image will be png encoded (therefore it may be bigger in size)
despite of file extension specified in the result flag.

## Showcases

### 🚩 Codefest’19

`stegify` was used for one of the *Capture The Flag* challenges in [**Codefest’19**](https://www.hackerrank.com/codefest19-ctf).

Participants were given a photo of a bunch of "innocent" cats. Nothing suspicious right? Think again!

You can read more [here](https://medium.com/bugbountywriteup/codefest19-ctf-writeups-a8f4e9b45d1) and [here](https://medium.com/@markonsecurity/image-challenges-1-cats-are-innocent-right-69cd4220bc99). 

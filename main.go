package main

import (
	"encoding/base64"
	"flag"
	"fmt"
	"os"
	"time"

	"golang.design/x/clipboard"
)

var decodeCommand *flag.FlagSet
var decodeInputFilePtr *string
var encodeCommand *flag.FlagSet
var encodeOutputNamePtr *string

func init() {
	decodeCommand = flag.NewFlagSet("decode", flag.ExitOnError)
	encodeCommand = flag.NewFlagSet("encode", flag.ExitOnError)
	decodeInputFilePtr = decodeCommand.String("f", "", "Specifies file content should be used for input instead of terminal argument. Path required.")
	encodeOutputNamePtr = encodeCommand.String("f", "", "Writes base64-encoded file to a new file. Name can be passed optionally.")
}

func main() {
	defaultErrorMessage := "Subcommand 'encode' or 'decode' required."
	if len(os.Args) < 2 {
		fmt.Println(defaultErrorMessage)
		os.Exit(1)
	}

	switch os.Args[1] {
	case "encode":
		encodeCommand.Parse(os.Args[2:])
		runEncodeCommand()
	case "decode":
		decodeCommand.Parse(os.Args[2:])
		runDecodeCommand()
	default:
		fmt.Println(defaultErrorMessage)
		os.Exit(1)
	}
}

func runEncodeCommand() {
	if len(os.Args) < 3 {
		printSubcategoryMenuAndExit(*encodeCommand, "Writes base64-encoded file to clipboard. Path to file required.")
	}

	inputFilename := os.Args[len(os.Args)-1]
	_, err := os.Stat(inputFilename)
	exitOnError(err)

	fileContentBase64, err := encode(inputFilename)
	exitOnError(err)
	if *encodeOutputNamePtr != "" {
		var outputFilename string
		if *encodeOutputNamePtr == inputFilename {
			outputFilename = fmt.Sprintf("%s.b64", inputFilename)
		} else {
			outputFilename = *encodeOutputNamePtr
		}
		exitOnError(os.WriteFile(outputFilename, []byte(fileContentBase64), 0644))
		fmt.Printf("%s has been converted to base64 and saved to %s.", inputFilename, outputFilename)
	} else {
		clipboard.Write(clipboard.FmtText, []byte(fileContentBase64))
		fmt.Printf("%s has been converted to base64 and saved to clipboard.", inputFilename)
	}
}

func encode(filename string) (string, error) {
	fileContent, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(fileContent), nil
}

func runDecodeCommand() {
	if len(os.Args) < 3 {
		printSubcategoryMenuAndExit(*decodeCommand, "Decodes base64 string from terminal. Base64 string argument required.")
	}

	var decoded []byte
	var err error
	if *decodeInputFilePtr != "" {
		decoded, err = decodeFromFileContent()
	} else {
		decoded, err = decodeFromArgs()
	}
	exitOnError(err)

	outputName := fmt.Sprintf("b64_%s.bin", time.Now().UTC().Format("20060102150405"))
	err = os.WriteFile(outputName, decoded, 0644)
	exitOnError(err)
	fmt.Printf("Base64 string has successfully been decoded. Binary has been saved to %s.", outputName)
}

func decodeFromArgs() ([]byte, error) {
	base64String := os.Args[len(os.Args)-1]
	return base64.StdEncoding.DecodeString(base64String)
}

func decodeFromFileContent() ([]byte, error) {
	base64Bytes, err := os.ReadFile(os.Args[len(os.Args)-1])
	if err != nil {
		return nil, err
	}
	content := string(base64Bytes)
	return base64.StdEncoding.DecodeString(content)
}

func printSubcategoryMenuAndExit(flagSet flag.FlagSet, text string) {
	if text != "" {
		fmt.Println(text)
	}
	flagSet.PrintDefaults()
	os.Exit(1)
}

func exitOnError(err error) {
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

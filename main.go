package main

import (
	"flag"
	"fmt"
	"image/png"
	"os"
	"regexp"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
	"github.com/boombuler/barcode/twooffive"
	"github.com/pschlump/dbgo"
	"github.com/pschlump/filelib"
)

var TextToEncode = flag.String("text", "Hello World", "What to Encode in the barcode.")
var Output = flag.String("output", "out.png", "Send output data to file...")
var Height = flag.Int("height", 256, "Height of output image.")
var Width = flag.Int("width", 256, "Width of output image.")
var Format = flag.String("format", "QR-Code", "One of { QR-Code, 2-of-5, ... }")

// Specific to... 2-of-5
var Interleaved = flag.Bool("interleaved", false, "2-of-5 interlaeved mode")

/*
2-of-5
Aztec-Code
Codabar
Code-128
Code-39
Code-93
Datamatrix
EAN-13
EAN-8
PDF-417
QR-Code
*/
var validCodes []string = []string{"2-of-5", "Aztec-Code", "Codabar", "Code-128", "Code-39", "Code-93", "Datamatrix", "EAN-13", "EAN-8", "PDF-417", "QR-Code"}

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "gen-barcode: Usage: %s [flags]\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.Parse() // Parse CLI arguments to this, --cfg <name>.json

	fns := flag.Args()
	if len(fns) != 0 {
		fmt.Printf("Extra arguments are not supported [%s]\n", fns)
		os.Exit(1)
	}

	if !filelib.InArray(*Format, validCodes) {
		fmt.Printf("invalid --format value, must be one of: %s\n", validCodes)
	}

	switch *Format {
	case "QR-Code":

		// Create the barcode
		qrCode, err := qr.Encode(*TextToEncode, qr.M, qr.Auto)
		if err != nil {
			dbgo.Fprintf(os.Stderr, "Error: %s AT %(LF)\n", err)
			os.Exit(1)
		}

		// Scale the barcode to 200x200 pixels
		qrCode, err = barcode.Scale(qrCode, *Width, *Height)
		if err != nil {
			dbgo.Fprintf(os.Stderr, "Error: %s AT %(LF)\n", err)
			os.Exit(1)
		}

		// create the output file
		file, err := os.Create(*Output)
		if err != nil {
			dbgo.Fprintf(os.Stderr, "Error: %s AT %(LF)\n", err)
			os.Exit(1)
		}
		defer file.Close()

		// encode the barcode as png
		png.Encode(file, qrCode)

		return

	case "2-of-5":

		// xyzzy - should check that it is all digits!
		numeric := regexp.MustCompile(`\d`).MatchString(*TextToEncode)
		if !numeric {
			fmt.Fprintf(os.Stderr, "2-of-5 can only encode digits, 0...9\n")
			os.Exit(1)
		}

		// Enven umber of digits if --interleaved
		if *Interleaved {
			if len(*TextToEncode)%2 != 0 {
				fmt.Fprintf(os.Stderr, "An odd length of %d can not use --interleaved=%v\n", len(*TextToEncode), *Interleaved)
				os.Exit(1)
			}
		}

		// fmt.Printf("Mod %d Len %d *Interleaved %v\n", len(*TextToEncode)%2, len(*TextToEncode), *Interleaved)

		// barCode, err := twooffive.Encode(*TextToEncode, interleaved)
		barCode, err := twooffive.Encode(*TextToEncode, *Interleaved)
		if err != nil || barCode == nil {
			fmt.Fprintf(os.Stderr, "Ivalid Code:%s\n", err)
		} else {

			// Scale the barcode to 200x200 pixels
			barCode, err = barcode.Scale(barCode, *Width, *Height)
			if err != nil {
				dbgo.Fprintf(os.Stderr, "Error: %s AT %(LF)\n", err)
				os.Exit(1)
			}

			// create the output file
			file, err := os.Create(*Output)
			if err != nil {
				// xyzzy
				dbgo.Fprintf(os.Stderr, "Error: %s AT %(LF)\n", err)
				os.Exit(1)
			}
			defer file.Close()

			// encode the barcode as png
			png.Encode(file, barCode)
		}
		return

	default:
		fmt.Fprintf(os.Stderr, "Invalid %s, Not supported yet\n", *Format)

	}

}

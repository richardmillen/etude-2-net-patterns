package main

import (
	"bytes"
	"flag"
	"fmt"
	"strconv"
	"strings"

	"github.com/richardmillen/etude-2-net-patterns/src-go/uuid"
)

var count = flag.Int("count", 1000000, "the number of uuid's to run the test against.")
var codePoints = flag.String("find", "0,0", `the bytes (chars) to look for in each uuid.
	n.b. it's currently only possible to specify decimal code point values. 

	Examples:

	look for two null chars:

		uuid-check --find=0,0
	
	look for ": ":

		uuid-check --find=58,32
	
	look for LF (line feed):

		uuid-check --find=10`)

func main() {
	flag.Parse()

	b := getBytesFromChars(*codePoints)
	if len(b) == 0 {
		fmt.Println("no code point values specified. run --help for usage info.")
		return
	}

	fmt.Println("testing", getPrettyInt(*count), "UUIDs for byte sequence:", b)

	found := 0
	needNewLine := false

	for n := 0; n < *count; n++ {
		if n > 0 && n%100000 == 0 {
			needNewLine = true
			fmt.Print(".")
		}

		id := uuid.New()
		if find(id, b) {
			found++
		}
	}

	if needNewLine {
		fmt.Println()
	}
	if found > 0 {
		fmt.Print("*** FOUND *** ")
	}
	fmt.Println(found, "UUID(s) of", getPrettyInt(*count), "contained", b)
}

func getBytesFromChars(s string) []byte {
	chars := strings.Split(s, ",")
	b := make([]byte, len(chars))

	for n, char := range chars {
		v, _ := strconv.Atoi(char)
		b[n] = byte(v)
	}

	return b
}

// find is called to look for seq within id.
func find(id uuid.Bytes, seq []byte) bool {
	b := []byte(id)

	for n := 0; n < (len(b)-len(seq))+1; n++ {
		cmp := b[n : n+len(seq)]
		if bytes.Equal(cmp, seq) {
			return true
		}
	}
	return false
}

func getPrettyInt(i int) string {
	s := strconv.Itoa(i)

	commaPos := len(s) % 3
	if commaPos == 0 {
		commaPos = 3
	}

	var buf bytes.Buffer
	for n := 0; n < len(s); n++ {
		if n == commaPos {
			buf.WriteByte(',')
			commaPos += 3
		}
		buf.WriteByte(s[n])
	}
	return buf.String()
}

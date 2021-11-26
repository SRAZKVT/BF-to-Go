package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func main() {
	args := os.Args
	if len(args) != 2 {
		fmt.Println("Please input valid file name")
		return
	}
	fileInput := args[1]
	fileOutput := "output.go"

	finput, err := os.Open(fileInput)
	if err != nil {
		log.Fatal(err)
	}
	foutput, err := os.Create(fileOutput)
	if err != nil {
		log.Fatal(err)
	}

	scanner := bufio.NewScanner(finput)
	var line [] byte
	var character byte
	indent := 1

	isFmtUsed := false
	isMemUsed := false
	isPtrUsed := false

	for scanner.Scan() {
		line = scanner.Bytes()

		for _, character = range line {
			switch character {
			case '>', '<':
				isPtrUsed = true
			case '[', '+', '-':
				isMemUsed = true
				isPtrUsed = true
			case '.', ',':
				isPtrUsed = true
				isMemUsed = true
				isFmtUsed = true
			}
		}
	}
	finput.Close()

	finput, err = os.Open(fileInput)
	if err != nil {
		log.Fatal(err)
	}
	scanner = bufio.NewScanner(finput)

	w := bufio.NewWriter(foutput)

	w.WriteString("package main\n")
	if isFmtUsed {
		w.WriteString("import \"fmt\"\n")
	}
	w.WriteString("func main() {\n")
	if isMemUsed {
		w.WriteString("var mem [65536] uint8\n")
	}
	if isPtrUsed {
		w.WriteString("var ptr uint16 = 0\n")
	}


	for scanner.Scan() {
		line = scanner.Bytes()

		for _, character = range line {

			if character == ']' {
				indent--
			}

			for i := 0; i < indent; i++ {
				w.WriteString("\t")
			}

			switch character {
			case '>':
				w.WriteString("ptr++\n")
			case '<':
				w.WriteString("ptr--\n")
			case '[':
				w.WriteString("for mem[ptr] != 0 {\n")
				indent++
			case ']':
				w.WriteString("}\n")
			case '+':
				w.WriteString("mem[ptr]++\n")
			case '-':
				w.WriteString("mem[ptr]--\n")
			case '.':
				w.WriteString("fmt.Printf(\"%c\", mem[ptr])\n")
			case ',':
				w.WriteString("fmt.Scan(mem[ptr])\n")
			}
			w.Flush()
		}
	}
	w.WriteString("}")
	w.Flush()
	foutput.Close()
	fmt.Println("Compilation completed.")
}

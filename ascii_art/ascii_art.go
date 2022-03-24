package ascii_art

import (
	"bufio"
	"errors"
	"os"
	"strings"
)

func AsciiArt(input, font string) (string, error) {

	// read the file with asci letters and put its contents into a string slice
	arrFile, err := readFile("ascii_art/" + font)

	if err != nil {
		return "", err
	}

	// split the line into words if there is a "\r\n" symbol
	arrWords := strings.Split(input, "\r\n")

	// put the recieved ascii-art words in the string variable
	asciWord, err := readLetters(arrWords, arrFile)

	if err != nil {
		return "", err
	}

	return asciWord, nil
}

// converts slice of asci words into a single string
func convertToString(asciWord [8]string) string {
	if len(asciWord[0]) == 0 {
		return ""
	}

	// adding each line of a slice, separating it with a new line
	outstring := ""
	for _, line := range asciWord {
		outstring += line + "\n"
	}
	return outstring
}

// read the file with asci letters and put its contents into a string slice
func readFile(args string) ([]string, error) {

	file, err := os.Open(args)

	if err != nil {
		return nil, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	arrFile := []string{}
	for scanner.Scan() {
		arrFile = append(arrFile, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return arrFile, nil
}

// put the recieved ascii-art words in the string variable
func readLetters(words, arrFile []string) (string, error) {

	res := ""
	// go over the words of slice
	for _, wrd := range words {

		if wrd == "" {
			res += "\n"
			continue
		}
		var asciWord [8]string

		// go over the letters of the words
		for _, let := range wrd {

			if let < ' ' || let > '~' {
				return "", errors.New("400, Bad Request")
			}

			// find the coordinate of the required letter in the asci slice
			position := (int(let) - 32) * 9

			// we go over the desired range of the slice with asci symbols
			for i, line := range arrFile[position+1 : position+9] {
				asciWord[i] += line
			}
		}
		res += convertToString(asciWord)
	}
	return res, nil
}

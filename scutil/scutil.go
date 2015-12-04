package scutil

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
)

var errEndOfBlock = errors.New("eob")

func parseDict(scanner *bufio.Scanner) (map[string]interface{}, error) {
	res := make(map[string]interface{})

	for scanner.Scan() {
		line := scanner.Text()
		if err := scanner.Err(); err != nil {
			return nil, err
		}
		key, value, err := parseLine(scanner, line)
		if err != nil {
			if err != errEndOfBlock {
				return nil, err
			}
			return res, nil
		}
		res[key] = value
	}
	return res, nil
}

func parseArray(scanner *bufio.Scanner) ([]interface{}, error) {
	res := make([]interface{}, 0)

	for scanner.Scan() {
		line := scanner.Text()
		if err := scanner.Err(); err != nil {
			return nil, err
		}
		_, value, err := parseLine(scanner, line)
		if err != nil {
			if err != errEndOfBlock {
				return nil, err
			}
			return res, nil
		}
		res = append(res, value)
	}

	return res, nil
}

func parseLine(scanner *bufio.Scanner, line string) (key string, value interface{}, err error) {
	parts := strings.SplitN(line, " : ", 2)
	if len(parts) == 1 {
		if line[len(line)-1] == '}' {
			return "", nil, errEndOfBlock
		}
		return "", nil, fmt.Errorf("do not know how to parse: %s", line)
	}

	switch p := parts[1]; p {
	case "<dictionary> {":
		value, err = parseDict(scanner)
	case "<array> {":
		value, err = parseArray(scanner)
	default:
		value = p
	}

	return strings.TrimSpace(parts[0]), value, err
}

// JSONEncode reads from r which contains a scutil --nc formatted data
// and writes the equivalent JSON structure in w.
func JSONEncode(r io.Reader, w io.Writer) error {
	br := bufio.NewReader(r)
	firstLine, err := br.ReadString('\n')
	if err != nil {
		return err
	}
	firstLine = strings.TrimRight(firstLine, "\n\r ")
	topKey := strings.TrimSuffix(firstLine, " <dictionary> {")
	if topKey == firstLine {
		return fmt.Errorf("first line should be of the form `key <dictionary> {` our expectation: %s", firstLine)
	}
	topMap := make(map[string]interface{})
	scanner := bufio.NewScanner(br)
	res, err := parseDict(scanner)
	if err != nil {
		return err
	}
	topMap[topKey] = res
	encoder := json.NewEncoder(w)
	if err := encoder.Encode(&topMap); err != nil {
		return fmt.Errorf("error encoding json: %s", err)
	}
	return nil
}

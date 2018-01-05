package gonx

import (
	"bufio"
	"bytes"
	"io"
	"strings"
	"testing"
)

func BenchmarkScannerReader(b *testing.B) {
	s := `89.234.89.123 [08/Nov/2013:13:39:18 +0000] "GET /api/foo/bar HTTP/1.1"`
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		file := strings.NewReader(s)
		scanner := bufio.NewScanner(file)
		scanner.Scan()
		scanner.Text()
		if err := scanner.Err(); err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkReaderReaderAppend(b *testing.B) {
	s := `89.234.89.123 [08/Nov/2013:13:39:18 +0000] "GET /api/foo/bar HTTP/1.1"`
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		file := strings.NewReader(s)
		reader := bufio.NewReader(file)
		_, err := readLineAppend(reader)
		if err != nil && err != io.EOF {
			b.Fatal(err)
		}
	}
}

func BenchmarkReaderReaderBuffer(b *testing.B) {
	s := `89.234.89.123 [08/Nov/2013:13:39:18 +0000] "GET /api/foo/bar HTTP/1.1"`
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		file := strings.NewReader(s)
		reader := bufio.NewReader(file)
		_, err := readLineBuffer(reader)
		if err != nil && err != io.EOF {
			b.Fatal(err)
		}
	}
}

func BenchmarkLongReaderReaderAppend(b *testing.B) {
	longStr := RandString(10 * 64 * 1024)
	s := `89.234.89.123 [08/Nov/2013:13:39:18 +0000] "GET ` + longStr + ` HTTP/1.1"`
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		file := strings.NewReader(s)
		reader := bufio.NewReader(file)
		_, err := readLineAppend(reader)
		if err != nil && err != io.EOF {
			b.Fatal(err)
		}
	}
}

func BenchmarkLongReaderReaderBuffer(b *testing.B) {
	longStr := RandString(10 * 64 * 1024)
	s := `89.234.89.123 [08/Nov/2013:13:39:18 +0000] "GET ` + longStr + ` HTTP/1.1"`
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		file := strings.NewReader(s)
		reader := bufio.NewReader(file)
		_, err := readLineBuffer(reader)
		if err != nil && err != io.EOF {
			b.Fatal(err)
		}
	}
}

func readLineAppend(reader *bufio.Reader) (string, error) {
	line, isPrefix, err := reader.ReadLine()
	if err != nil {
		return "", err
	}
	if !isPrefix {
		return string(line), nil
	}
	var ln []byte
	for isPrefix && err == nil {
		ln, isPrefix, err = reader.ReadLine()
		if err == nil {
			line = append(line, ln...)
		}
	}
	return string(line), err
}

func readLineBuffer(reader *bufio.Reader) (string, error) {
	line, isPrefix, err := reader.ReadLine()
	if err != nil {
		return "", err
	}
	if !isPrefix {
		return string(line), nil
	}
	var buffer bytes.Buffer
	_, err = buffer.Write(line)
	for isPrefix && err == nil {
		line, isPrefix, err = reader.ReadLine()
		if err == nil {
			_, err = buffer.Write(line)
		}
	}
	return buffer.String(), err
}

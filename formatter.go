package gostructfmt

import "fmt"

// Format parses and formats Go source using the struct literal formatter.
// The returned source is deterministic and contains no partial output when
// parsing or formatting fails.
func Format(filename string, source []byte) ([]byte, error) {
	file, err := Parse(filename, source)
	if err != nil {
		return nil, err
	}

	if err := FormatStructLiterals(file); err != nil {
		return nil, fmt.Errorf("format %q: %w", displayFilename(filename), err)
	}

	output, err := file.Print()
	if err != nil {
		return nil, fmt.Errorf("format %q: %w", displayFilename(filename), err)
	}

	return output, nil
}

func displayFilename(filename string) string {
	if filename == "" {
		return "<input>"
	}

	return filename
}

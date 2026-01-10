package reader

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/ledongthuc/pdf"
)

func ExtractCheckData(filename string) (*CheckData, error) {
	if !strings.HasSuffix(filename, ".pdf") {
		return nil, fmt.Errorf(
			"error in ExtractCheckData(%s): wrong file format (must be a pdf).",
			filename,
		)
	}

	f, r, err := pdf.Open(filename)
	if err != nil {
		return nil, fmt.Errorf(
			"error in ExtractCheckData(%s): couldn't open a file: %v",
			filename,
			err,
		)
	}
	defer f.Close()

	b, err := r.GetPlainText()
	if err != nil {
		return nil, fmt.Errorf(
			"error in ExtractCheckData(%s): couldn't extract plain text from a PDF: %v",
			filename,
			err,
		)
	}

	var buf bytes.Buffer
	_, err = buf.ReadFrom(b)
	if err != nil {
		return nil, fmt.Errorf(
			"error in ExtractCheckData(%s): couldn't read from a reader: %v",
			filename,
			err,
		)
	}

	cd, err := NewCheckData(buf.String())
	if err != nil {
		return nil, fmt.Errorf(
			"error in ExtractCheckData(%s): couldn't extract check data: %v",
			filename,
			err,
		)
	}

	return cd, nil
}

func ExtractAllCheckData(filenames ...string) ([]*CheckData, error) {
	var cdList []*CheckData
	if len(filenames) == 0 {
		return cdList, fmt.Errorf(
			"error in ExtractAllCheckData(%v): zero filenames were provided.",
			filenames,
		)
	}

	for i, filename := range filenames {
		cd, err := ExtractCheckData(filename)
		if err != nil {
			return cdList, fmt.Errorf(
				"error in ExtractAllCheckData(%v): couldn't extract data from a file with index '%d': %v.",
				filenames,
				i,
				err,
			)
		}

		cdList = append(cdList, cd)
	}

	return cdList, nil
}

package reader

import (
	"fmt"
	re "regexp"
	"strconv"
	"strings"
	"time"
)

type CheckData struct {
	rawText string

	TimeOfCreation time.Time

	Items []CheckItem

	Total float64

	OrgName string
}

func (a *CheckData) IsEqual(b *CheckData) bool {
	if a.OrgName != b.OrgName || a.Total != b.Total {
		return false
	}

	if !a.TimeOfCreation.Equal(b.TimeOfCreation) {
		return false
	}

	for i, itemA := range a.Items {
		if !itemA.IsEqual(&b.Items[i]) {
			return false
		}
	}

	return true
}

func NewCheckData(rawCheckText string) (*CheckData, error) {
	rawCheckText = trimAllLeadingNonsense(rawCheckText)
	r := re.MustCompile(`В т\.ч\. НДС \d+%`)
	lines := r.Split(rawCheckText, -1)
	if len(lines) == 0 {
		return nil, fmt.Errorf(
			"error in reader.NewCheckData(%s): couldn't split the text by 'В т\\.ч\\. НДС \\d+%%'-pattern",
			rawCheckText,
		)
	}

	timeOfCreation, err := extractTime(lines[0])
	if err != nil {
		return nil, fmt.Errorf(
			"error in reader.NewCheckData(%s): couldn't extract time: %v",
			rawCheckText,
			err,
		)
	}

	normalizedItems, err := normalizeItems(lines)
	if err != nil {
		return nil, fmt.Errorf(
			"error in reader.NewCheckData(%s): couldn't normalize items: %v",
			rawCheckText,
			err,
		)
	}

	extractedItems, err := extractItems(normalizedItems)
	if err != nil {
		return nil, fmt.Errorf(
			"error in reader.NewCheckData(%s): couldn't extract items: %v",
			rawCheckText,
			err,
		)
	}

	orgName, err := extractOrgName(normalizedItems[len(lines)-1])
	if err != nil {
		return nil, fmt.Errorf(
			"error in reader.NewCheckData(%s): couldn't extract org name: %v",
			rawCheckText,
			err,
		)
	}

	total, err := extractTotal(normalizedItems[len(lines)-1])
	if err != nil {
		return nil, fmt.Errorf(
			"error in reader.NewCheckData(%s): couldn't extract total: %v",
			rawCheckText,
			err,
		)
	}

	return &CheckData{
		rawText:        rawCheckText,
		TimeOfCreation: timeOfCreation,
		Items:          extractedItems,
		Total:          total,
		OrgName:        orgName,
	}, nil
}

const (
	TBANK_TIME_FORMAT = "02.01.2006   15:04:05"
)

func extractTime(rawLine string) (time.Time, error) {
	idx := strings.Index(rawLine, "\n")
	if idx == -1 {
		return time.Time{}, fmt.Errorf("error in reader.extractTime(%s): couldn't find '\\n'-symbol", rawLine)
	}

	t, err := time.Parse(TBANK_TIME_FORMAT, rawLine[:idx])
	if err != nil {
		return t, fmt.Errorf("error in reader.extractTime(%s): couldn't parse string '%s': %v", rawLine, rawLine[:idx], err)
	}

	return t, err
}

func trimAllLeadingNonsense(s string) string {
	for {
		s = strings.TrimSpace(s)
		if !strings.HasPrefix(s, "\n") {
			break
		}
		s = strings.TrimPrefix(s, "\n")
	}
	return s
}

func normalizeItems(rawLines []string) ([]string, error) {
	var lines []string
	idx := strings.Index(rawLines[0], "\n")
	if idx == -1 {
		return lines, fmt.Errorf("error in reader.normalizeItems(%v): couldn't normalize first line (can't find '\\n'-symbol); line: '%s'", lines, rawLines[0])
	}
	rawLines[0] = rawLines[0][idx:]

	for _, rawLine := range rawLines {
		lines = append(
			lines, strings.TrimSpace(
				strings.ReplaceAll(rawLine, "\n", " "),
			),
		)
	}

	return lines, nil
}

func extractItems(lines []string) ([]CheckItem, error) {
	var checkItems []CheckItem

	for _, line := range lines[:len(lines)-1] {
		item, err := NewCheckItem(line)
		if err != nil {
			return nil, fmt.Errorf(
				"error in reader.extractItems(%v): couldn't extract an item from line '%s': %v",
				lines,
				line,
				err,
			)
		}
		checkItems = append(checkItems, item)
	}

	return checkItems, nil
}

func extractOrgName(rawText string) (string, error) {
	//rawText = trimAllLeadingNonsense(rawText)
	idx := strings.Index(rawText, "РН ККТ: ") + len("РН ККТ: ")
	if idx == -1 {
		return "", fmt.Errorf("error in reader.extractOrgName(%s): couldn't find substr 'РН ККТ: '", rawText)
	}
	rawText = rawText[idx:]

	startIdx := strings.Index(rawText, " ") + 1
	if startIdx == -1 {
		return "", fmt.Errorf("error in reader.extractOrgName(%s): couldn't find ' '-symbol for startIdx", rawText)
	}

	endIdx := strings.Index(rawText, "ИНН")
	if endIdx == -1 {
		return "", fmt.Errorf("error in reader.extractOrgName(%s): couldn't find substr 'ИНН' for endIdx", rawText)
	}

	return strings.TrimSpace(rawText[startIdx:endIdx]), nil
}

func extractTotal(rawText string) (float64, error) {
	rawText = trimAllLeadingNonsense(rawText)
	startIdx := strings.Index(rawText, "Электронный платеж ") + len("Электронный платеж ")
	if startIdx == -1 {
		return -1, fmt.Errorf("error in reader.extractTotal(%s): couldn't find substr 'Электронный платеж ' for startIdx", rawText)
	}

	endIdx := strings.Index(rawText, "Цена")
	if startIdx == -1 {
		return -1, fmt.Errorf("error in reader.extractTotal(%s): couldn't find substr 'Цена' for endIdx", rawText)
	}

	normalizedTotal := normalizePriceStr(rawText[startIdx:endIdx])
	total, err := strconv.ParseFloat(normalizedTotal, 64)
	if err != nil {
		return -1, fmt.Errorf(
			"error in reader.extractTotal(%s): couldn't parse the total from line '%s': %v",
			rawText,
			normalizedTotal,
			err,
		)
	}

	return total, nil
}

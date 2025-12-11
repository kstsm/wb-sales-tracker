package export

import (
	"encoding/csv"
	"errors"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"strings"
	"time"
	"unicode"

	"github.com/google/uuid"
)

func GetStructColumnNames(data any) ([]string, error) {
	v := reflect.ValueOf(data)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return nil, errors.New("data must be a struct or pointer to struct")
	}

	t := v.Type()
	cols := make([]string, 0, t.NumField())

	for i := range t.NumField() {
		field := t.Field(i)
		jsonTag := field.Tag.Get("json")
		if jsonTag != "" && jsonTag != "-" {
			if idx := strings.Index(jsonTag, ","); idx != -1 {
				jsonTag = jsonTag[:idx]
			}
			cols = append(cols, formatColumnName(jsonTag))
		} else {
			cols = append(cols, formatColumnName(field.Name))
		}
	}

	return cols, nil
}

func formatColumnName(name string) string {
	if name == "" {
		return name
	}
	runes := []rune(name)
	if len(runes) > 0 {
		runes[0] = unicode.ToLower(runes[0])
	}

	return string(runes)
}

func ConvertStructToCSV(data any) ([]string, error) {
	v := reflect.ValueOf(data)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return nil, errors.New("data must be a struct or pointer to struct")
	}

	t := v.Type()
	record := make([]string, 0, t.NumField())

	for i := range v.NumField() {
		field := v.Field(i)
		fieldValue := field.Interface()

		var strValue string
		switch val := fieldValue.(type) {
		case uuid.UUID:
			strValue = val.String()
		case time.Time:
			strValue = val.Format(time.RFC3339)
		case float64:
			strValue = strconv.FormatFloat(val, 'f', 2, 64)
		case string:
			strValue = val
		default:
			strValue = fmt.Sprintf("%v", val)
		}

		record = append(record, strValue)
	}

	return record, nil
}

func getSampleData(v reflect.Value) any {
	if v.Len() > 0 {
		return v.Index(0).Interface()
	}

	sliceType := v.Type()
	elemType := sliceType.Elem()

	if elemType.Kind() == reflect.Ptr {
		elemType = elemType.Elem()
	}

	return reflect.New(elemType).Interface()
}

func WriteItemsCSV(w io.Writer, cols []string, items any) error {
	writer := csv.NewWriter(w)
	defer writer.Flush()

	v := reflect.ValueOf(items)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Slice {
		return errors.New("items must be a slice")
	}

	if len(cols) == 0 {
		sampleData := getSampleData(v)
		var err error
		cols, err = GetStructColumnNames(sampleData)
		if err != nil {
			return fmt.Errorf("GetStructColumnNames: %w", err)
		}
	}

	if err := writer.Write(cols); err != nil {
		return fmt.Errorf("writer.Write: %w", err)
	}

	for i := range v.Len() {
		item := v.Index(i).Interface()
		record, err := ConvertStructToCSV(item)
		if err != nil {
			return fmt.Errorf("ConvertStructToCSV: %w", err)
		}

		if writeErr := writer.Write(record); writeErr != nil {
			return fmt.Errorf("writer.Write: %w", writeErr)
		}
	}

	return nil
}

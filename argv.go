package argv

import (
	"reflect"
	"strconv"
	"strings"
	"time"
)

const (
	annotationTag = "argv"
)

// field mapper
type FieldParser func(in string) (any, error)

// struct types that should not be recursively parsed
var (
	reservedFieldTypes = []string{"time.Time"}
	fieldParser        = make(map[string]FieldParser, 0)
)

// add a custom field parser
func AddParser(in string, fn FieldParser) {
	fieldParser[in] = fn
}

func AddReservedType(t string) {
	reservedFieldTypes = append(reservedFieldTypes, t)
}

// Check if field type name is reserved
func isReserved(t string) bool {
	for _, v := range reservedFieldTypes {
		if v == t {
			return true
		}
	}
	return false
}

func ParseNames(dest any) ([]string, error) {
	t := reflect.TypeOf(dest)
	v := reflect.ValueOf(dest)
	if t.Kind() != reflect.Ptr {
		return nil, ErrInvalidDest
	}
	t = t.Elem()
	v = v.Elem()
	if t.Kind() != reflect.Struct {
		return nil, ErrInvalidDestType
	}
	result := make([]string, 0)

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		kind := v.Field(i).Kind()
		reserved := isReserved(field.Type.String())
		if kind == reflect.Struct && !reserved {
			vals, err := ParseNames(v.Field(i).Addr().Interface())
			if err != nil {
				return nil, err
			}
			result = append(result, vals...)
		} else {
			fieldName, _ := parseTag(t.Field(i).Tag.Get(annotationTag))
			if len(fieldName) == 0 {
				continue
			}
			if !v.Field(i).CanInterface() {
				continue
			}
			result = append(result, fieldName)
		}
	}
	return result, nil
}

func extractArgs(args []string) (map[string]string, error) {
	if len(args)%2 > 0 {
		return nil, ErrInvalidParameterCount
	}
	result := make(map[string]string, 0)
	i := 0
	for i < len(args) {
		argName := args[i]
		if strings.HasPrefix(argName, "--") {
			argName = argName[2:]
		} else if strings.HasPrefix(argName, "-") {
			argName = argName[1:]
		}
		result[argName] = args[i+1]
		i += 2
	}
	return result, nil
}

func parseTag(tag string) (string, bool) {
	if len(tag) == 0 {
		return "", false
	}
	toks := strings.Split(tag, ",")
	if len(toks) == 1 {
		return toks[0], false
	}
	return toks[0], toks[1] == "optional"
}

func ParseArgv(dest any, argv []string) error {
	if len(argv) == 0 {
		return ErrEmptyArgs
	}
	args, err := extractArgs(argv)
	if err != nil {
		return err
	}
	return parseArgv(dest, args)
}

func parseArgv(dest any, args map[string]string) error {
	if len(args) == 0 {
		return nil
	}

	t := reflect.TypeOf(dest)
	v := reflect.ValueOf(dest)
	if t.Kind() != reflect.Ptr {
		return ErrInvalidDest
	}
	t = t.Elem()
	v = v.Elem()
	if t.Kind() != reflect.Struct {
		return ErrInvalidDestType
	}

	for i := 0; i < t.NumField(); i++ {
		field := v.Field(i)
		kind := v.Field(i).Kind()
		reserved := isReserved(field.Type().String())
		if kind == reflect.Struct && !reserved {
			if err := parseArgv(v.Field(i).Addr().Interface(), args); err != nil {
				return err
			}
		} else {
			fieldName, optional := parseTag(t.Field(i).Tag.Get(annotationTag))
			if len(fieldName) == 0 {
				continue
			}
			if !v.Field(i).CanInterface() {
				continue
			}
			// field has a tag, but it is not settable
			if kind != reflect.Interface {
				if !field.CanSet() {
					return ErrReadOnly(fieldName)
				}
			}

			if fValue, ok := args[fieldName]; !ok {
				if !optional {
					return ErrMissingValue(fieldName)
				}
				continue
			} else {
				fType := field.Type().String()
				switch fType {
				case "time.Time":
					v, err := mapTime(fValue)
					if err != nil {
						return ErrInvalidValue(fieldName, err)
					}
					field.Set(reflect.ValueOf(v))

				case "bool":
					v, err := parseBool(fValue)
					if err != nil {
						return ErrInvalidValue(fieldName, err)
					}
					field.SetBool(v)

				case "byte", "uint8":
					v, err := parseUint(fValue, 8)
					if err != nil {
						return ErrInvalidValue(fieldName, err)
					}
					field.SetUint(v)

				case "int8":
					v, err := parseInt(fValue, 8)
					if err != nil {
						return ErrInvalidValue(fieldName, err)
					}
					field.SetInt(v)
				case "uint", "uint32":
					v, err := parseUint(fValue, 32)
					if err != nil {
						return ErrInvalidValue(fieldName, err)
					}
					field.SetUint(v)

				case "uint64":
					v, err := parseUint(fValue, 64)
					if err != nil {
						return ErrInvalidValue(fieldName, err)
					}
					field.SetUint(v)

				case "int", "int32":
					v, err := parseInt(fValue, 32)
					if err != nil {
						return ErrInvalidValue(fieldName, err)
					}
					field.SetInt(v)
				case "int64":
					v, err := parseInt(fValue, 64)
					if err != nil {
						return ErrInvalidValue(fieldName, err)
					}
					field.SetInt(v)
				case "float32":
					v, err := parseFloat(fValue, 32)
					if err != nil {
						return ErrInvalidValue(fieldName, err)
					}
					field.SetFloat(v)

				case "float64":
					v, err := parseFloat(fValue, 64)
					if err != nil {
						return ErrInvalidValue(fieldName, err)
					}
					field.SetFloat(v)

				case "string":
					field.SetString(fValue)

				case "[]string":
					field.Set(reflect.ValueOf(parseStringArray(fValue)))
				default:
					if fn, ok := fieldParser[fType]; ok {
						v, err := fn(fValue)
						if err != nil {
							return ErrInvalidValue(fieldName, err)
						}
						field.Set(reflect.ValueOf(v))
					} else {
						return ErrNotSupported(fieldName)
					}
				}
			}
		}
	}
	return nil
}

func parseBool(in string) (bool, error) {
	return strconv.ParseBool(in)
}

func parseUint(in string, size int) (uint64, error) {
	return strconv.ParseUint(in, 10, size)
}

func parseInt(in string, size int) (int64, error) {
	return strconv.ParseInt(in, 10, size)
}

func mapTime(in string) (time.Time, error) {
	return time.Parse(time.RFC3339, in)
}

func parseFloat(in string, size int) (float64, error) {
	return strconv.ParseFloat(in, size)
}

func parseStringArray(in string) []string {
	result := make([]string, 0)
	if len(in) == 0 {
		return result
	}
	for _, v := range strings.Split(in, ",") {
		result = append(result, strings.TrimSpace(v))
	}
	return result
}

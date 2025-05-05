package bencode

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"reflect"
	"strconv"
)

func Parse(v any, r io.Reader) error {
	val := reflect.ValueOf(v)
    if val.Kind() != reflect.Ptr {
        return errors.New("v must be a pointer.")
    }    
	reader := bufio.NewReader(r)

	return parseVal(val, reader)
}

func parseVal(v reflect.Value, r *bufio.Reader) error {
	val := reflect.Indirect(v)
	switch val.Kind() {
		case reflect.String:
			res, err := parseString(r)
			if err != nil {
				return fmt.Errorf("Parse error: %v", err)
			}
			if !val.CanSet() {
				return errors.New("Cant change v.")
			}
			val.SetString(res)
		case reflect.Int:
			res, err := parseInt(r)
			if err != nil {
				return fmt.Errorf("Parse error: %v", err)
			}
			if !val.CanSet() {
				return errors.New("Cant change v.")
			}
			val.SetInt(int64(res))
		case reflect.Slice:
			if err := buildSlice(v, r); err != nil {
				return fmt.Errorf("Parse error: %v", err)
			}
		case reflect.Struct:
			if err := buildDict(v, r); err != nil {
				return fmt.Errorf("Parse error: %v", err)
			}
			
		default:
			return errors.New("Unupported v type.")
	}

	return nil	

}

func isDigit(r rune) bool {
	return r >= '0' && r <= '9'
}

func parseString(reader *bufio.Reader) (string, error) {
	r, _, err := reader.ReadRune()
	reader.UnreadRune()
	if err != nil {
		return "", ReadError
	}
	if !isDigit(r) {
		return "", errors.New("String parse error. Not a string.")
	}

	lStr := ""
	for {
		rune, _, err := reader.ReadRune()
		if err != nil {
			return "", err
		}
		if !isDigit(rune) {
			break
		}
		lStr += string(rune)
	}
	
	l, err := strconv.Atoi(lStr)

	if err != nil {
		return "", err
	}
	
	s := make([]byte, l)
	n, err := io.ReadFull(reader, s)


	if n < l || err != nil {
		return "", ReadError
	}

	return string(s), nil

}

func parseInt(reader *bufio.Reader) (int, error) {
	r, _, err := reader.ReadRune()
	if err != nil {
		return 0, ReadError
	}
	if r != 'i' {
		reader.UnreadRune()
		return 0, errors.New("Not expected type parsed.")
	}

	numStr, err := reader.ReadString('e')
	
	if err != nil {
		return 0, ReadError
	}

	numStr = numStr[:len(numStr)-1]

	num, err := strconv.Atoi(numStr)

	if err != nil {
		return 0, err
	}

	return num, nil

}

func buildSlice(v reflect.Value, r *bufio.Reader) error {
	rune, _, err := r.ReadRune()

	if err != nil {
		return ReadError
	}
	if rune != 'l'{
		r.UnreadRune()
		return errors.New("List parse error. Not a list.")
	}

	sVal := reflect.Indirect(v)
	sType := sVal.Type()
	elemType := sType.Elem()

	newSlice := reflect.MakeSlice(sType, 0, 0)

	for {
		bytes, err := r.Peek(1)
		if err != nil {
			return ReadError
		}
		t := bytes[0]
		if t == 'e' {
			_, err := r.ReadByte()
			if err != nil {
				return ReadError
			}
			break
		}

		newElemValPtr := reflect.New(elemType)
		err = parseVal(newElemValPtr, r)
		if err != nil {
			return err
		}
		newElemVal := newElemValPtr.Elem()

		if !newElemVal.Type().AssignableTo(elemType) {
			return errors.New("Incompatible types.")
		}
		newSlice = reflect.Append(newSlice , newElemVal)
	}

	if !newSlice.Type().AssignableTo(sType) {
		return errors.New("Cant assign to slice.")
	}
	sVal.Set(newSlice)
	return nil
}


func buildDict(v reflect.Value, r *bufio.Reader) error {
	val := reflect.Indirect(v)
	
	rune, _, err := r.ReadRune()
	if err != nil {
		return ReadError
	}
	if rune != 'd' {
		r.UnreadRune()
		return errors.New("Dict parse error. Not a dict.")
	}

	for {
		key, err := parseString(r)
		if err != nil {
			bytes, err := r.Peek(1)
			if err != nil {
				return ReadError
			}
			if bytes[0] == 'e' {
				_, err := r.ReadByte()
				if err != nil {
					return ReadError
				}
			}
			break
		}
		fieldVal, err := getFieldWithMatchingTag(key, val)

		if err != nil {
			continue
		}
		parseVal(fieldVal, r)	
	}


	return nil
}

func getFieldWithMatchingTag(key string, st reflect.Value) (reflect.Value, error) {
	stType := st.Type()

	for i := range stType.NumField() {
		field := stType.Field(i)
		fVal := st.Field(i)
		
		tag, ok := field.Tag.Lookup("bencode")
		
		if !ok {
			continue
		}

		if tag == key {
			return fVal, nil
		}
	}

	return reflect.Value{}, fmt.Errorf("Tag %s not found.\n", key)

}
var ReadError = errors.New("Parse error. Reading error.")

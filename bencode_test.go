package bencode

import (
	"slices"
	"strings"
	"testing"
)

func TestUnmarshalString(t *testing.T) {
	input := "7:testval"
	reader := strings.NewReader(input)
	var res string

	err := Unmarshal(&res, reader)

	if err != nil {
		t.Fatal(err)
	}

	if res != "testval" {
		t.Errorf("'%s' != 'testval'", res)
	}
}

func TestUnmarshalStringIllformed(t *testing.T) {
	input := "7:testva"
	reader := strings.NewReader(input)
	var res string

	err := Unmarshal(&res, reader)

	if err == nil {
		t.Fatalf("Parsing ill-formed bencode must return error.")
	}
}

func TestUnmarshalInt(t *testing.T) {
	input := "i1337e"
	reader := strings.NewReader(input)
	var res int

	err := Unmarshal(&res, reader)
	if err != nil {
		t.Fatal(err)
	}

	if res != 1337 {
		t.Fatalf("%d != 1337", res)
	}
}

func TestUnmarshalIntIllformed(t *testing.T) {

	input := "i1notdigit1e"
	reader := strings.NewReader(input)
	var res int

	err := Unmarshal(&res, reader)

	if err == nil {
		t.Fatalf("Parsing ill-formed bencode must return error.")
	}
}

func TestUnmarshalListInt(t *testing.T) {
	input := "li1ei2ei3ee"
	reader := strings.NewReader(input)
	var res []int
	expected := []int{1, 2, 3}

	err := Unmarshal(&res, reader)
	if err != nil {
		t.Fatal(err)
	}

	if !slices.Equal(res, expected) {
		t.Errorf("%v != %v", res, expected)
	}
}

func TestUnmarshalListIntIllformed1(t *testing.T) {
	input := "li1ei2ei3e"
	reader := strings.NewReader(input)
	var res []int

	err := Unmarshal(&res, reader)
	if err == nil {
		t.Fatalf("Parsing ill-formed bencode must return error.")
	}
}

func TestUnmarshalListIntIllformed2(t *testing.T) {
	input := "li1ei2e2:abe"
	reader := strings.NewReader(input)
	var res []int

	err := Unmarshal(&res, reader)
	if err == nil {
		t.Fatalf("Parsing ill-formed bencode must return error.")
	}
}

func TestUnmarshalListString(t *testing.T) {
	input := "l2:aa2:bb2:cce"
	reader := strings.NewReader(input)
	var res []string
	expected := []string{"aa", "bb", "cc"}

	err := Unmarshal(&res, reader)
	if err != nil {
		t.Fatal(err)
	}

	if !slices.Equal(res, expected) {
		t.Errorf("%v != %v", res, expected)
	}
}

func TestUnmarshalListEmpty(t *testing.T) {
	input := "le"
	reader := strings.NewReader(input)
	var res []string
	expected := []string{}

	err := Unmarshal(&res, reader)
	if err != nil {
		t.Fatal(err)
	}

	if !slices.Equal(res, expected) {
		t.Errorf("%v != %v", res, expected)
	}
}

func TestUnmarshalListOfListOfStrings(t *testing.T) {
	input := "ll2:aa2:bb2:ccel2:aa2:bb2:ccel2:aa2:bb2:ccee"
	reader := strings.NewReader(input)
	var res [][]string
	expected := [][]string{{"aa", "bb", "cc"}, {"aa", "bb", "cc"}, {"aa", "bb", "cc"}}

	err := Unmarshal(&res, reader)
	if err != nil {
		t.Fatal(err)
	}

	if len(res) == 0 {
		t.Fatal("res len is 0")
	}

	for i, s := range res {
		if !slices.Equal(s, expected[i]) {
			t.Fatalf("%v != %v", res, expected)
		}
	}
}

func TestUnmarshalListOfListOfInt(t *testing.T) {
	input := "lli1ei2ei3eeli1ei2ei3eeli1ei2ei3eee"
	reader := strings.NewReader(input)
	var res [][]int
	expected := [][]int{{1, 2, 3}, {1, 2, 3}, {1, 2, 3}}

	err := Unmarshal(&res, reader)
	if err != nil {
		t.Fatal(err)
	}

	if len(res) == 0 {
		t.Fatal("res len is 0")
	}

	for i, s := range res {
		if !slices.Equal(s, expected[i]) {
			t.Fatalf("%v != %v", res, expected)
		}
	}
}

func TestUnmarshalListOfDicts(t *testing.T) {
	type testStruct struct {
		Field1 int    `bencode:"sf1"`
		Field2 string `bencode:"sf2"`
	}

	input := "ld3:sf1i1e3:sf25:test1ed3:sf1i2e3:sf25:test2ee"
	reader := strings.NewReader(input)

	res := []testStruct{}
	expected := []testStruct{
		{Field1: 1, Field2: "test1"},
		{Field1: 2, Field2: "test2"},
	}

	err := Unmarshal(&res, reader)
	if err != nil {
		t.Fatal(err)
	}

	if !slices.Equal(res, expected) {
		t.Errorf("%v != %v", res, expected)
	}
}

func TestUnmarshalDictFlat(t *testing.T) {
	input := "d6:field1i1e6:field2i2e6:field34:teste"
	reader := strings.NewReader(input)

	res := struct {
		Field1 int    `bencode:"field1"`
		Field2 int    `bencode:"field2"`
		Field3 string `bencode:"field3"`
	}{}
	expected := struct {
		Field1 int    `bencode:"field1"`
		Field2 int    `bencode:"field2"`
		Field3 string `bencode:"field3"`
	}{1, 2, "test"}

	err := Unmarshal(&res, reader)
	if err != nil {
		t.Fatal(err)
	}

	if res != expected {
		t.Errorf("%v != %v", res, expected)
	}
}

func TestUnmarshalDictWithListField(t *testing.T) {
	input := "d6:field1i1e6:field2li1ei2ei3ee6:field34:teste"
	reader := strings.NewReader(input)

	type testStruct struct {
		Field1 int    `bencode:"field1"`
		Field2 []int  `bencode:"field2"`
		Field3 string `bencode:"field3"`
	}
	res := testStruct{}

	expected := testStruct{Field1: 1, Field2: []int{1, 2, 3}, Field3: "test"}

	err := Unmarshal(&res, reader)
	if err != nil {
		t.Fatal(err)
	}

	isEqual := true
	if res.Field1 != expected.Field1 {
		isEqual = false
	}
	if res.Field3 != expected.Field3 {
		isEqual = false
	}
	if !slices.Equal(res.Field2, expected.Field2) {
		isEqual = false
	}
	if !isEqual {
		t.Fatalf("%v != %v", res, expected)
	}
}

func TestUnmarshalDictOnlyRequiredFields1(t *testing.T) {
	input := "d6:field1i1e6:field2i2e6:field34:teste"
	reader := strings.NewReader(input)

	res := struct {
		Field1 int    `bencode:"field1"`
		Field3 string `bencode:"field3"`
	}{}
	expected := struct {
		Field1 int    `bencode:"field1"`
		Field3 string `bencode:"field3"`
	}{1, "test"}

	err := Unmarshal(&res, reader)
	if err != nil {
		t.Fatal(err)
	}

	if res != expected {
		t.Errorf("%v != %v", res, expected)
	}
}

func TestUnmarshalDictOnlyRequiredFields2(t *testing.T) {
	input := "d6:field1i1e6:field2i2e6:field2li1ei2ei3ee6:field34:teste"
	reader := strings.NewReader(input)

	res := struct {
		Field1 int    `bencode:"field1"`
		Field3 string `bencode:"field3"`
	}{}
	expected := struct {
		Field1 int    `bencode:"field1"`
		Field3 string `bencode:"field3"`
	}{1, "test"}

	err := Unmarshal(&res, reader)
	if err != nil {
		t.Fatal(err)
	}

	if res != expected {
		t.Errorf("%v != %v", res, expected)
	}
}

func TestUnmarshalDictOnlyRequiredFields3(t *testing.T) {
	input := "d6:field1i1e6:field7d2:f1l2:xx2:yy:2:zze2:f2d2:f1i1e2:f2i2eee6:field2i2e6:field2li1ei2ei3ee6:field34:teste"
	reader := strings.NewReader(input)

	res := struct {
		Field1 int    `bencode:"field1"`
		Field3 string `bencode:"field3"`
	}{}
	expected := struct {
		Field1 int    `bencode:"field1"`
		Field3 string `bencode:"field3"`
	}{1, "test"}

	err := Unmarshal(&res, reader)
	if err != nil {
		t.Fatal(err)
	}

	if res != expected {
		t.Errorf("%v != %v", res, expected)
	}
}

func TestUnmarshalDictWithDictField(t *testing.T) {
	input := "d6:field1i1e6:field2d3:sf1i1e3:sf27:subteste6:field34:teste"
	reader := strings.NewReader(input)

	type subStruct struct {
		SubField1 int    `bencode:"sf1"`
		SubField2 string `bencode:"sf2"`
	}

	type testStruct struct {
		Field1 int       `bencode:"field1"`
		Field2 subStruct `bencode:"field2"`
		Field3 string    `bencode:"field3"`
	}
	res := testStruct{}

	expected := testStruct{Field1: 1, Field2: subStruct{SubField1: 1, SubField2: "subtest"}, Field3: "test"}

	err := Unmarshal(&res, reader)
	if err != nil {
		t.Fatal(err)
	}

	if res != expected {
		t.Fatalf("%v != %v", res, expected)
	}
}

func TestMarshalInt(t *testing.T) {
	input := 123
	expected := "i123e"

	res, err := Marshal(input)
	if err != nil {
		t.Fatal(err)
	}

	resString := string(res)

	if resString != expected {
		t.Fatalf("%v != %v", resString, expected)
	}
}

func TestMarshalString(t *testing.T) {
	input := "test"
	expected := "4:test"

	res, err := Marshal(input)
	if err != nil {
		t.Fatal(err)
	}

	resString := string(res)

	if resString != expected {
		t.Fatalf("'%s' != '%s'", resString, expected)
	}
}

func TestMarshalListOfStrings(t *testing.T) {
	input := []string{"str1", "str2", "str3"}
	expected := "l4:str14:str24:str3e"

	res, err := Marshal(input)
	if err != nil {
		t.Fatal(err)
	}

	resString := string(res)

	if resString != expected {
		t.Fatalf("'%s' != '%s'", resString, expected)
	}
}

func TestMarshalListOfInt(t *testing.T) {
	input := []int{11, 22, 33}
	expected := "li11ei22ei33ee"

	res, err := Marshal(input)
	if err != nil {
		t.Fatal(err)
	}

	resString := string(res)

	if resString != expected {
		t.Fatalf("'%s' != '%s'", resString, expected)
	}
}

func TestMarshalListOfLists(t *testing.T) {
	input := [][]int{{11, 22, 33}, {11, 22, 33}, {11, 22, 33}}
	expected := "lli11ei22ei33eeli11ei22ei33eeli11ei22ei33eee"

	res, err := Marshal(input)
	if err != nil {
		t.Fatal(err)
	}

	resString := string(res)

	if resString != expected {
		t.Fatalf("'%s' != '%s'", resString, expected)
	}
}

func TestMarshalStructFlat(t *testing.T) {
	input := struct {
		Field2 int    `bencode:"field2"`
		Field1 int    `bencode:"field1"`
		Field3 string `bencode:"field3"`
	}{2, 1, "test"}
	expected := "d6:field1i1e6:field2i2e6:field34:teste"

	res, err := Marshal(input)
	if err != nil {
		t.Fatal(err)
	}

	resString := string(res)

	if resString != expected {
		t.Fatalf("'%s' != '%s'", resString, expected)
	}
}

func TestMarshalStructNested(t *testing.T) {
	type subStruct struct {
		SubField2 int    `bencode:"sf2"`
		SubField1 string `bencode:"sf1"`
		SubField3 []int  `bencode:"sf3"`
	}

	type testStruct struct {
		Field2 int       `bencode:"field2"`
		Field1 subStruct `bencode:"field1"`
		Field3 string    `bencode:"field3"`
	}

	input := testStruct{
		Field2: 345,
		Field1: subStruct{SubField2: 123, SubField1: "subtest", SubField3: []int{1, 2, 3}},
		Field3: "test",
	}
	expected := "d6:field1d3:sf17:subtest3:sf2i123e3:sf3li1ei2ei3eee6:field2i345e6:field34:teste"

	res, err := Marshal(input)
	if err != nil {
		t.Fatal(err)
	}

	resString := string(res)

	if resString != expected {
		t.Fatalf("'%s' != '%s'", resString, expected)
	}
}

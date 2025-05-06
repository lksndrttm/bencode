package bencode

import (
	"slices"
	"strings"
	"testing"
)

func TestParseString(t *testing.T) {
	input := "7:testval"
	reader := strings.NewReader(input)
	var res string

	err := Parse(&res, reader)

	if err != nil {
		t.Fatal(err)
	}

	if res != "testval" {
		t.Errorf("'%s' != 'testval'", res)
	}
}

func TestParseStringIllformed(t *testing.T) {
	input := "7:testva"
	reader := strings.NewReader(input)
	var res string

	err := Parse(&res, reader)

	if err == nil {
		t.Fatalf("Parsing ill-formed bencode must return error.")
	}
}

func TestParseInt(t *testing.T) {
	input := "i1337e"
	reader := strings.NewReader(input)
	var res int

	err := Parse(&res, reader)
	if err != nil {
		t.Fatal(err)
	}

	if res != 1337 {
		t.Fatalf("%d != 1337", res)
	}
}

func TestParseIntIllformed(t *testing.T) {

	input := "i1notdigit1e"
	reader := strings.NewReader(input)
	var res int

	err := Parse(&res, reader)

	if err == nil {
		t.Fatalf("Parsing ill-formed bencode must return error.")
	}
}

func TestParseListInt(t *testing.T) {
	input := "li1ei2ei3ee"
	reader := strings.NewReader(input)
	var res []int
	expected := []int{1, 2, 3}

	err := Parse(&res, reader)
	if err != nil {
		t.Fatal(err)
	}

	if !slices.Equal(res, expected) {
		t.Errorf("%v != %v", res, expected)
	}
}

func TestParseListIntIllformed1(t *testing.T) {
	input := "li1ei2ei3e"
	reader := strings.NewReader(input)
	var res []int

	err := Parse(&res, reader)
	if err == nil {
		t.Fatalf("Parsing ill-formed bencode must return error.")
	}
}

func TestParseListIntIllformed2(t *testing.T) {
	input := "li1ei2e2:abe"
	reader := strings.NewReader(input)
	var res []int

	err := Parse(&res, reader)
	if err == nil {
		t.Fatalf("Parsing ill-formed bencode must return error.")
	}
}

func TestParseListString(t *testing.T) {
	input := "l2:aa2:bb2:cce"
	reader := strings.NewReader(input)
	var res []string
	expected := []string{"aa", "bb", "cc"}

	err := Parse(&res, reader)
	if err != nil {
		t.Fatal(err)
	}

	if !slices.Equal(res, expected) {
		t.Errorf("%v != %v", res, expected)
	}
}

func TestParseListEmpty(t *testing.T) {
	input := "le"
	reader := strings.NewReader(input)
	var res []string
	expected := []string{}

	err := Parse(&res, reader)
	if err != nil {
		t.Fatal(err)
	}

	if !slices.Equal(res, expected) {
		t.Errorf("%v != %v", res, expected)
	}
}

func TestParseListOfListOfStrings(t *testing.T) {
	input := "ll2:aa2:bb2:ccel2:aa2:bb2:ccel2:aa2:bb2:ccee"
	reader := strings.NewReader(input)
	var res [][]string
	expected := [][]string{{"aa", "bb", "cc"}, {"aa", "bb", "cc"}, {"aa", "bb", "cc"}}

	err := Parse(&res, reader)
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

func TestParseListOfListOfInt(t *testing.T) {
	input := "lli1ei2ei3eeli1ei2ei3eeli1ei2ei3eee"
	reader := strings.NewReader(input)
	var res [][]int
	expected := [][]int{{1, 2, 3}, {1, 2, 3}, {1, 2, 3}}

	err := Parse(&res, reader)
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

func TestParseListOfDicts(t *testing.T) {
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

	err := Parse(&res, reader)
	if err != nil {
		t.Fatal(err)
	}

	if !slices.Equal(res, expected) {
		t.Errorf("%v != %v", res, expected)
	}
}

func TestParseDictFlat(t *testing.T) {
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

	err := Parse(&res, reader)
	if err != nil {
		t.Fatal(err)
	}

	if res != expected {
		t.Errorf("%v != %v", res, expected)
	}
}

func TestParseDictWithListField(t *testing.T) {
	input := "d6:field1i1e6:field2li1ei2ei3ee6:field34:teste"
	reader := strings.NewReader(input)

	type testStruct struct {
		Field1 int    `bencode:"field1"`
		Field2 []int  `bencode:"field2"`
		Field3 string `bencode:"field3"`
	}
	res := testStruct{}

	expected := testStruct{Field1: 1, Field2: []int{1, 2, 3}, Field3: "test"}

	err := Parse(&res, reader)
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

func TestParseDictOnlyRequiredFields1(t *testing.T) {
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

	err := Parse(&res, reader)
	if err != nil {
		t.Fatal(err)
	}

	if res != expected {
		t.Errorf("%v != %v", res, expected)
	}
}

func TestParseDictOnlyRequiredFields2(t *testing.T) {
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

	err := Parse(&res, reader)
	if err != nil {
		t.Fatal(err)
	}

	if res != expected {
		t.Errorf("%v != %v", res, expected)
	}
}

func TestParseDictOnlyRequiredFields3(t *testing.T) {
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

	err := Parse(&res, reader)
	if err != nil {
		t.Fatal(err)
	}

	if res != expected {
		t.Errorf("%v != %v", res, expected)
	}
}



func TestParseDictWithDictField(t *testing.T) {
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

	err := Parse(&res, reader)
	if err != nil {
		t.Fatal(err)
	}

	if res != expected {
		t.Fatalf("%v != %v", res, expected)
	}
}

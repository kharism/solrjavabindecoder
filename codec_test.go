package solrjavacodec

import (
	"io/ioutil"
	"strconv"
	"testing"
)

func TestRead1(t *testing.T) {
	t.Skip()
	bytes, _ := ioutil.ReadFile("select_javabin.bin")
	//fd,_:=os.Open("select_javabin.bin")
	data := map[string]interface{}{}
	UnmarshalByte(&data, bytes)
}
func TestReadString(t *testing.T) {
	inputBytes := []byte{0x02, 0xe0, 0x2e, 0x72, 0x65, 0x73, 0x70, 0x6f, 0x6e, 0x73, 0x65, 0x48, 0x65, 0x61, 0x64, 0x65, 0x72}
	//reader := bytes.NewBuffer(inputBytes)
	var kk string
	kk = ""
	err := UnmarshalByte(&kk, inputBytes)
	if err != nil {
		t.Error(err.Error())
	}
	t.Log(kk)
	if kk != "responseHeader" {
		t.Fail()
	}
}
func TestReadMapStringString(t *testing.T) {
	inputBytes := []byte{0x02, 0xa2, 0xe0, 0x24, 0x6b, 0x65, 0x79, 0x31, 0x26, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x31, 0xe0, 0x24, 0x6b, 0x65, 0x79, 0x32, 0x26, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x32}
	var kk map[string]interface{}
	//kk = ""
	err := UnmarshalByte(&kk, inputBytes)
	if err != nil {
		t.Error(err.Error())
	}
	t.Log("Done Unmarshaling")
	t.Log(kk)
	keys := []string{"key1", "key2"}
	values := []string{"value1", "value2"}
	for i, k := range keys {
		if _, ok := kk[k]; !ok {
			t.Log("Key Not Found")
			t.Fail()
			return
		}
		if kk[k] != values[i] {
			t.Log("Value Not same")
			t.Fail()
			return
		}
	}
}
func hexStrToByte(inputString string) ([]byte, error) {
	inputBytes := []byte{}
	for i := 0; i < len(inputString); i += 2 {
		str := inputString[i : i+2]
		uu, err := strconv.ParseInt(str, 16, 9)
		if err != nil {
			//t.Log(err.Error())
			return inputBytes, err
		}
		//t.Logf("%s %02x", str, uu)
		inputBytes = append(inputBytes, byte(uu))
	}
	return inputBytes, nil
}
func TestReadVIntBools(t *testing.T) {
	inputString := "02a3e0246b65793150c002e0246b65793201e0246b65793302"
	inputBytes, _ := hexStrToByte(inputString) //[]byte{}
	var kk map[string]interface{}
	//kk = ""
	err := UnmarshalByte(&kk, inputBytes)
	if err != nil {
		t.Error(err.Error())
	}
	keys := []string{"key1", "key2", "key3"}
	values := []interface{}{int(5120), true, false}
	t.Log("Done Unmarshaling")
	for i, k := range keys {
		if _, ok := kk[k]; !ok {
			t.Log("Key Not Found")
			t.Fail()
			return
		}
		if kk[k] != values[i] {
			t.Log("Value Not same")
			t.Fail()
			return
		}
	}
}
func TestReadFloatBool(t *testing.T) {
	inputString := "02a3e0246b65793108424ccccde0246b65793201e0246b65793302"
	inputBytes, _ := hexStrToByte(inputString) //[]byte{}
	var kk map[string]interface{}
	//kk = ""
	err := UnmarshalByte(&kk, inputBytes)
	if err != nil {
		t.Error(err.Error())
	}
	keys := []string{"key1", "key2", "key3"}
	values := []interface{}{float32(51.20), true, false}
	t.Log("Done Unmarshaling")
	for i, k := range keys {
		if _, ok := kk[k]; !ok {
			t.Log("Key Not Found")
			t.Fail()
			return
		}
		if kk[k] != values[i] {
			t.Log("Value Not same")
			t.Fail()
			return
		}
	}
}
func TestReadIntDoule(t *testing.T) {
	inputString := "02a2e0246b6579314ce0246b657932054028cccccccccccd"
	inputBytes, _ := hexStrToByte(inputString) //[]byte{}
	var kk map[string]interface{}
	//kk = ""
	err := UnmarshalByte(&kk, inputBytes)
	if err != nil {
		t.Error(err.Error())
	}
	keys := []string{"key1", "key2"}
	values := []interface{}{int(12), float64(12.4)}
	t.Log("Done Unmarshaling")
	for i, k := range keys {
		if _, ok := kk[k]; !ok {
			t.Log("Key Not Found")
			t.Fail()
			return
		}
		if kk[k] != values[i] {
			t.Log("Value Not same")
			t.Fail()
			return
		}
	}
}
func TestReadMapStringLongString(t *testing.T) {
	inputString := "02a2e0246b6579312676616c756531e0246b6579323fbd0176616c7565327665656565727979796c6f6f6f6f6f6f6f6f6f6f6f6f6f6f6f6f6f6f6e6764616161616161616161747474747474746161616161616161616161616161616161616161616161616161616161616161616161616161616161616161616161616161616161616161616161616161616161616161616161616161616161616161616161616161616161616161616161616161616161616161616161616161616161616161616161616161616161616161616161616161616161616161616161616161616161616161616161616161616161616161616161"
	inputBytes, _ := hexStrToByte(inputString) //[]byte{}
	// for i := 0; i < len(inputString); i += 2 {
	// 	str := inputString[i : i+2]
	// 	uu, err := strconv.ParseInt(str, 16, 9)
	// 	if err != nil {
	// 		t.Log(err.Error())
	// 		return
	// 	}
	// 	//t.Logf("%s %02x", str, uu)
	// 	inputBytes = append(inputBytes, byte(uu))
	// }
	var kk map[string]interface{}
	//kk = ""
	err := UnmarshalByte(&kk, inputBytes)
	if err != nil {
		t.Error(err.Error())
	}
	t.Log("Done Unmarshaling")
	t.Log(kk)
	keys := []string{"key1", "key2"}
	values := []string{"value1", "value2veeeeryyyloooooooooooooooooongdaaaaaaaaatttttttaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"}
	for i, k := range keys {
		if _, ok := kk[k]; !ok {
			t.Log("Key Not Found")
			t.Fail()
			return
		}
		if kk[k] != values[i] {
			t.Log("Value Not same")
			t.Fail()
			return
		}
	}
}

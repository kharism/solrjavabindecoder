package solrjavacodec

import (
	"bytes"
	"fmt"
	"io"
	"math"
	"reflect"

	tk "github.com/eaciit/toolkit"

	"github.com/eaciit/errorlib"
)

const (
	//basic data type

	NULL       = 0
	BOOL_TRUE  = 1
	BOOL_FALSE = 2
	BYTE       = 3
	SHORT      = 4
	DOUBLE     = 5
	INT        = 6
	LONG       = 7
	FLOAT      = 8
	DATE       = 9
	MAP        = 10
	SOLRDOC    = 11
	SOLRDOCLST = 12
	BYTEARR    = 13
	ITERATOR   = 14

	STR           = byte(1 << 5)
	SINT          = byte(2 << 5)
	SLONG         = byte(3 << 5)
	ARR           = byte(4 << 5) //
	ORDERED_MAP   = byte(5 << 5) // SimpleOrderedMap (a NamedList subclass, and more common)
	NAMED_LST     = byte(6 << 5) // NamedList
	EXTERN_STRING = byte(7 << 5)
)

var packageName = "solrjavacodec"
var modCursorDecode = "Decoder"
var readStringAsCharSeq = false

func UnmarshalByte(m interface{}, data []byte) error {
	if !tk.IsPointer(m) {
		return errorlib.Error(packageName, modCursorDecode, "Fetch", "Model object should be pointer")
	}
	buff := bytes.NewBuffer(data)
	buff.ReadByte() //ambil versi, gak dipakai dulu
	fmt.Println("LLLL")
	return readVal(m, buff)
}
func UnmarshalStream(m interface{}, dataStream *bytes.Buffer) error {
	e := readVal(m, dataStream)
	if e == io.EOF {
		return nil
	} else {
		return e
	}
}

func readVal(m interface{}, reader *bytes.Buffer) error {
	return readObject(m, reader)
}
func readObject(m interface{}, dis *bytes.Buffer) error {
	tagByte, err := dis.ReadByte()
	fmt.Printf("TAG %x\n", tagByte)
	if err != nil {
		return err
	}
	checkType := tagByte >> 5
	switch checkType {
	case STR >> 5:
		err = readStr(dis, readStringAsCharSeq, &tagByte, m)
		if err != nil {
			return err
		}
	case SINT >> 5:
		intRes, err := readSmallInt(dis, tagByte)
		if err != nil {
			return err
		}
		setValue(m, intRes)
	case EXTERN_STRING >> 5:
		fmt.Println("Read ExternString")
		err = readExternString(m, dis, &tagByte)
	case ORDERED_MAP >> 5:
		err = readOrderedMap(m, dis, tagByte)
		return err
	}
	switch tagByte {
	case NULL:
		return nil
	case BOOL_FALSE:
		setValue(m, false)
		return nil
	case BOOL_TRUE:
		setValue(m, true)
		return nil
	case FLOAT:
		flt, err := readFloat(dis, tagByte)
		if err != nil {
			return err
		}
		setValue(m, flt)
		return nil
	case DOUBLE:
		dbl, err := readDouble(dis, tagByte)
		if err != nil {
			return err
		}
		setValue(m, dbl)
		return nil
	}
	return nil
}
func readExternString(m interface{}, dis *bytes.Buffer, tagByte *byte) error {
	idx, err := readSize(dis, *tagByte)
	if err != nil {
		return err
	}
	if idx != 0 {
		//do something later
	} else {
		//only do this at the moment
		//idx == 0 means it has a string value
		*tagByte, err = dis.ReadByte()
		if err != nil {
			return err
		}
		readStr(dis, false, tagByte, m)
	}
	return nil
}
func readOrderedMap(m interface{}, dis *bytes.Buffer, tagByte byte) error {
	fmt.Println("Read Ordered Map")
	sz, err := readSize(dis, tagByte)
	if err != nil {
		return err
	}
	newMap := map[string]interface{}{}
	for ii := 0; ii < sz; ii++ {
		key1 := ""
		err := readVal(&key1, dis)
		if err != nil {
			return err
		}
		//fmt.Println("FoundKey", key1)
		var data interface{}
		err = readVal(&data, dis)
		if err != nil {
			return err
		}
		//fmt.Println("FoundData", data)
		newMap[key1] = data
	}
	*m.(*map[string]interface{}) = newMap
	return nil
}
func readSize(dis *bytes.Buffer, tagByte byte) (int, error) {
	var sz int
	sz = int(tagByte) & 0x1f
	fmt.Printf("ReadSize %x %x\n", sz, tagByte)
	if sz == int(0x1f) {
		fmt.Println("Ada tambahan")
		szAddition, err := readVInt(dis, tagByte)
		if err != nil {
			return 0, err
		}
		sz += szAddition
	}
	fmt.Println("Total read size", sz)
	return sz, nil
}
func readFloat(dis *bytes.Buffer, tagByte byte) (float32, error) {
	res, err := readInt(dis, tagByte)
	if err != nil {
		return 0, err
	}

	//int64 bit to float64
	//b := make([]byte, 8)
	//binary.LittleEndian.PutUint64(b, res)
	fmt.Println(res)
	resFloat := math.Float32frombits(res)
	return resFloat, nil
}
func readDouble(dis *bytes.Buffer, tagByte byte) (float64, error) {
	res, err := readLong(dis, tagByte)
	if err != nil {
		return 0, err
	}

	//int64 bit to float64
	//b := make([]byte, 8)
	//binary.LittleEndian.PutUint64(b, res)
	fmt.Println(res)
	resDouble := math.Float64frombits(res)
	return resDouble, nil
}
func readInt(dis *bytes.Buffer, tagByte byte) (uint32, error) {
	var hasil uint32
	//var idx uint
	for idxC := 24; idxC >= 0; idxC -= 8 {
		tempByte, err := dis.ReadByte()
		if err != nil {
			//fmt.Println("CurIDX", idx)
			return 0, err
		}
		hasil |= uint32(tempByte) << uint(idxC)
	}
	return hasil, nil
}
func readLong(dis *bytes.Buffer, tagByte byte) (uint64, error) {
	var hasil uint64
	var idx uint
	for idxC := 56; idxC >= 0; idxC -= 8 {
		tempByte, err := dis.ReadByte()
		if err != nil {
			fmt.Println("CurIDX", idx)
			return 0, err
		}
		hasil |= uint64(tempByte) << uint(idxC)
	}
	return hasil, nil
	// return  (((long)dis.ReadByte()) << 56)
	//         | (((long)readUnsignedByte()) << 48)
	//         | (((long)readUnsignedByte()) << 40)
	//         | (((long)readUnsignedByte()) << 32)
	//         | (((long)readUnsignedByte()) << 24)
	//         | (readUnsignedByte() << 16)
	//         | (readUnsignedByte() << 8)
	//         | (readUnsignedByte());
}
func readSmallInt(dis *bytes.Buffer, tagByte byte) (int, error) {
	var v int
	v = int(tagByte) & 0x0F
	if (tagByte & 0x10) != 0 {
		aa, err := readVInt(dis, tagByte)
		if err != nil {
			return 0, err
		}
		v = (aa << 4) | v
	}
	return v, nil
}

// this function is used to read any data that have length more than 255, maybe, not so sure
// about that
func readVInt(dis *bytes.Buffer, tagByte byte) (int, error) {
	b, err := dis.ReadByte()
	if err != nil {
		return 0, err
	}
	fmt.Printf("Read VINT %x\n", b)
	var i uint
	i = uint(b) & 0x7F
	var shift uint
	for shift = 7; (b & 0x80) != 0; shift += 7 {
		b, err = dis.ReadByte()
		if err != nil {
			return 0, err
		}
		i |= (uint(b) & 0x7F) << uint(shift)
	}
	fmt.Println("result read VINT", i)
	return int(i), nil
}

// need option to read as characte sequence, but for the moment force read as UTF8/UTF16 character
func readUtf8(dis *bytes.Buffer) error {

	return nil
}
func _readStr(dis *bytes.Buffer, sz int, m interface{}) error {
	tempBuff := make([]byte, sz)
	_, err := dis.Read(tempBuff)
	if err != nil {
		return err
	}
	newString := string(tempBuff)
	//fmt.Println("newString", newString)
	//*m.(*string) = newString
	setValue(m, newString)
	return nil
}
func setValue(m interface{}, value interface{}) {
	//v := reflect.TypeOf(m).Elem().Elem()
	//iv := reflect.New(v).Interface()
	reflect.ValueOf(m).Elem().Set(reflect.ValueOf(value))
}
func readStr(dis *bytes.Buffer, readStringAsCharSeq bool, tagByte *byte, m interface{}) error {
	//return readUtf8(dis)
	fmt.Println("Read STR")
	sz, err := readSize(dis, *tagByte)
	if err != nil {
		return err
	}
	return _readStr(dis, sz, m)
}

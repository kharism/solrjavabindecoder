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

	NULL             = 0
	BOOL_TRUE        = 1
	BOOL_FALSE       = 2
	BYTE             = 3
	SHORT            = 4
	DOUBLE           = 5
	INT              = 6
	LONG             = 7
	FLOAT            = 8
	DATE             = 9
	MAP              = 10
	SOLRDOC          = 11
	SOLRDOCLST       = 12
	BYTEARR          = 13
	ITERATOR         = 14
	END              = 15
	SOLRINPUTDOC     = 16
	MAP_ENTRY_ITER   = 17
	ENUM_FIELD_VALUE = 18
	MAP_ENTRY        = 19

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
	dataStream := bytes.NewBuffer(data)
	return UnmarshalStream(m, dataStream)
}
func UnmarshalStream(m interface{}, dataStream *bytes.Buffer) error {
	if !tk.IsPointer(m) {
		return errorlib.Error(packageName, modCursorDecode, "Fetch", "Model object should be pointer")
	}
	//buff := bytes.NewBuffer([]])
	dataStream.ReadByte() //ambil versi, gak dipakai dulu
	// fmt.Println("LLLL")
	stringCache := []string{}
	e := readVal(m, dataStream, &stringCache)
	if e == io.EOF {
		return nil
	} else {
		return e
	}
}

func readVal(m interface{}, reader *bytes.Buffer, stringCache *[]string) error {
	return readObject(m, reader, stringCache)
}

var end_obj_struct struct {
	ISEND bool
}

const END_OBJ = "HELLOSOLDER"

func readObject(m interface{}, dis *bytes.Buffer, stringCache *[]string) error {
	tagByte, err := dis.ReadByte()
	fmt.Printf("TAG %x\n", tagByte)
	if err != nil {
		return err
	}
	checkType := tagByte >> 5
	switch checkType {
	case STR >> 5:
		err = readStr(dis, readStringAsCharSeq, &tagByte, m)
		fmt.Println("Read String", *(m.(*interface{})))
		if err != nil {
			return err
		}
	case SINT >> 5:
		intRes, err := readSmallInt(dis, tagByte)
		if err != nil {
			return err
		}
		setValue(m, intRes)
	case SLONG >> 5:
		// fmt.Println("Read SLONG")
		longRes, err := readSmallLong(dis, tagByte)
		if err != nil {
			return err
		}
		setValue(m, longRes)
	case ARR >> 5:
		// fmt.Println("Read ARR")
		err := readArray(m, dis, &tagByte, stringCache)
		if err != nil {
			return err
		}
	case EXTERN_STRING >> 5:
		fmt.Println("Read ExternString")
		err = readExternString(m, dis, &tagByte, stringCache)
		if _, ok := (m.(*string)); ok {
			fmt.Println("Read ExternString", *(m.(*string)))
		}

		if err != nil {
			return err
		}
	case ORDERED_MAP >> 5:
		err = readOrderedMap(m, dis, tagByte, stringCache)
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
	case INT:
		ii, err := readInt(dis, tagByte)
		if err != nil {
			return err
		}
		setValue(m, ii)
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
		// case MAP:
		// 	readMap(m, dis, tagByte)
	case LONG:
		lg, err := readLong(dis, tagByte)
		if err != nil {
			return err
		}
		setValue(m, lg)
	case SOLRDOCLST:
		fmt.Println("ReadSolrDocList")
		err := readSolrDocumentList(m, dis, tagByte, stringCache)
		if err != nil {
			return err
		}
	case ITERATOR:
		fmt.Println("ReadIterator")
		err := readIterator(m, dis, tagByte, stringCache)
		if err != nil {
			return err
		}
	case END:
		setValue(m, END_OBJ)
		return nil
	case SOLRDOC:
		fmt.Println("ReadSolrDoc")
		err := readSolrDocument(m, dis, tagByte, stringCache)
		if err != nil {
			return err
		}
	case MAP:
		fmt.Println("found MAP")
		err := readMap(m, dis, tagByte, stringCache)
		if err != nil {
			return err
		}
	}

	return nil
}

//type SolrDocument map[string]interface{}
type SolrDocumentList struct {
	NumFound int64
	Start    int64
	MaxScore float32
	Docs     []interface{}
}

func readIterator(m interface{}, dis *bytes.Buffer, tagByte byte, stringCache *[]string) error {
	results := []interface{}{}
	for true {
		var newItem interface{}
		err := readVal(&newItem, dis, stringCache)
		if err != nil {
			return err
		}
		if newItem.(string) == END_OBJ {
			break
		}
		results = append(results, newItem)
	}
	setValue(m, results)
	return nil
}
func readSolrDocument(m interface{}, dis *bytes.Buffer, tagByte byte, stringCache *[]string) error {
	tagByte, err := dis.ReadByte()
	if err != nil {
		return err
	}
	size, err := readSize(dis, tagByte)
	if err != nil {
		return err
	}
	results := map[string]interface{}{}
	for i := 0; i < size; i++ {
		var keyStr string
		var keyIface interface{}
		var obj interface{}
		err := readVal(&keyIface, dis, stringCache)
		if err != nil {
			return err
		}

		if _, ok := keyIface.(string); ok {

			keyStr = keyIface.(string)
			err := readVal(&obj, dis, stringCache)
			if err != nil {
				return err
			}
			results[keyStr] = obj
		} else {
			fmt.Println("Found Non string", keyIface)
		}
	}
	setValue(m, results)
	return nil
}
func readSolrDocumentList(m interface{}, dis *bytes.Buffer, tagByte byte, stringCache *[]string) error {
	lists := []interface{}{}
	newList := SolrDocumentList{}
	err := readVal(&lists, dis, stringCache)
	if err != nil {
		return err
	}
	newList.NumFound = lists[0].(int64)
	newList.Start = lists[1].(int64)
	if lists[2] != nil {
		newList.MaxScore = lists[2].(float32)
	}

	newList.Docs = []interface{}{}
	err = readVal(&(newList.Docs), dis, stringCache)
	if err != nil {
		return err
	}
	setValue(m, newList)
	return nil
}

func readArray(m interface{}, dis *bytes.Buffer, tagByte *byte, stringCache *[]string) error {
	sz, err := readSize(dis, *tagByte)
	if err != nil {
		return err
	}
	l := []interface{}{}
	for i := 0; i < sz; i++ {
		var data interface{}
		readVal(&data, dis, stringCache)
		//fmt.Println("DATA", data)
		l = append(l, data)
	}
	//fmt.Println("ARRAY IS", l)
	setValue(m, l)
	return nil
}
func readExternString(m interface{}, dis *bytes.Buffer, tagByte *byte, stringCache *[]string) error {
	idx, err := readSize(dis, *tagByte)
	if err != nil {
		return err
	}
	if idx != 0 {
		//do something later
		setValue(m, (*stringCache)[idx-1])
		return nil
	} else {
		//only do this at the moment
		//idx == 0 means it has a string value
		*tagByte, err = dis.ReadByte()
		if err != nil {
			return err
		}
		err := readStr(dis, false, tagByte, m)
		if err != nil {
			return err
		}
		if _, ok := m.(*string); ok {
			ll := m.(*string)
			*stringCache = append(*stringCache, *ll)
		} else {
			ll := (*(m.(*interface{}))).(string)
			*stringCache = append(*stringCache, ll)
			//fmt.Println("Indirect type is:", reflect.Indirect(reflect.ValueOf(m)).Elem().Type())
		}

	}
	return nil
}

func readMap(m interface{}, dis *bytes.Buffer, tagByte byte, stringCache *[]string) error {
	sz, err := readVInt(dis, tagByte)
	if err != nil {
		return err
	}
	newMap := map[interface{}]interface{}{}
	for ii := 0; ii < sz; ii++ {
		var key1 interface{}
		err = readVal(&key1, dis, stringCache)
		if err != nil {
			return err
		}
		var data interface{}
		err = readVal(&data, dis, stringCache)
		if err != nil {
			return err
		}
		//fmt.Println("FoundData", data)
		newMap[key1] = data
	}
	setValue(m, newMap)
	return nil

}
func readOrderedMap(m interface{}, dis *bytes.Buffer, tagByte byte, stringCache *[]string) error {
	// fmt.Println("Read Ordered Map")
	sz, err := readSize(dis, tagByte)
	if err != nil {
		return err
	}
	newMap := map[string]interface{}{}
	for ii := 0; ii < sz; ii++ {
		key1 := ""
		err := readVal(&key1, dis, stringCache)
		if err != nil {
			return err
		}
		//fmt.Println("FoundKey", key1)
		var data interface{}
		err = readVal(&data, dis, stringCache)
		if err != nil {
			return err
		}
		//fmt.Println("FoundData", data)
		newMap[key1] = data
	}
	//*m.(*map[string]interface{}) = newMap
	setValue(m, newMap)
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
	//fmt.Println("Total read size", sz)
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
	//fmt.Println(res)
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
	//fmt.Println(res)
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
	//var idx uint
	for idxC := 56; idxC >= 0; idxC -= 8 {
		tempByte, err := dis.ReadByte()
		if err != nil {
			//fmt.Println("CurIDX", idx)
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
func readSmallLong(dis *bytes.Buffer, tagByte byte) (int64, error) {
	var v int64
	v = int64(tagByte) & 0x0F
	if (tagByte & 0x10) != 0 {
		aa, err := readVLong(dis, tagByte)
		if err != nil {
			return 0, err
		}
		v = (aa << 4) | v
	}
	return v, nil
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
	//fmt.Printf("Read VINT %x\n", b)
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
	// fmt.Println("result read VINT", i)
	return int(i), nil
}
func readVLong(dis *bytes.Buffer, tagByte byte) (int64, error) {
	b, err := dis.ReadByte()
	if err != nil {
		return 0, err
	}
	//fmt.Printf("Read VINT %x\n", b)
	var i uint64
	i = uint64(b) & 0x7F
	var shift uint
	for shift = 7; (b & 0x80) != 0; shift += 7 {
		b, err = dis.ReadByte()
		if err != nil {
			return 0, err
		}
		i |= (uint64(b) & 0x7F) << uint(shift)
	}
	// fmt.Println("result read VINT", i)
	return int64(i), nil
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
	// fmt.Println("Read STR")
	sz, err := readSize(dis, *tagByte)
	if err != nil {
		return err
	}
	return _readStr(dis, sz, m)
}

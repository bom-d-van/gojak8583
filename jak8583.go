package gojak8583

import (
	"fmt"
	"sort"
	"strconv"
)

type Builder struct{}

func Build() {}

type Message struct {
	MTI       string
	BitMap    string
	RawBitMap string
	Fields    []int
	Data      map[int]Data
}

type Data struct {
	string
	original interface{}
}

func (d Data) String() string {
	return d.string
}

func (d Data) Number() (int64, error) {
	return strconv.ParseInt(d.string, 10, 64)
}

func (d Data) Binary() (r string, err error) {
	hex, err := strconv.ParseInt(d.string, 16, 64)
	if err != nil {
		return
	}
	r = strconv.FormatInt(hex, 2)
	return
}

func Parse(raw string) (msg Message, err error) {
	msg.Data = map[int]Data{}
	msg.MTI = raw[0:4]
	raw = raw[4:]

	if raw, err = msg.parseBitMap(raw); err != nil {
		return
	}

	if err = msg.parseData(raw); err != nil {
		return
	}

	return
}

func (m *Message) parseBitMap(raw string) (r string, err error) {
	if m.BitMap, err = toBitMap(raw[0:16]); err != nil {
		return
	}
	m.RawBitMap = raw[0:16]
	raw = raw[16:]

	if m.BitMap[0] == '1' {
		var secondary string
		if secondary, err = toBitMap(raw[0:16]); err != nil {
			return
		}
		m.BitMap += secondary
		m.RawBitMap += raw[0:16]
		raw = raw[16:]
	}

	for i, _ := range m.BitMap {
		b, _ := strconv.Atoi(string(m.BitMap[i]))
		if i > 0 && b == 1 {
			m.Fields = append(m.Fields, i+1)
		}
	}

	r = raw

	return
}

func toBitMap(src string) (dst string, err error) {
	for _, code := range src {
		n, err := strconv.ParseInt(string(code), 16, 64)
		if err != nil {
			return "", fmt.Errorf("strconv.ParseInt(string(%s), 16, 64): %s", code, err)
		}
		code := strconv.FormatInt(n, 2)
		for len(code) < 4 {
			code = "0" + code
		}
		dst += code
	}

	for len(dst) < 64 {
		dst = "0" + dst
	}

	return
}

func (m *Message) parseData(raw string) (err error) {
	for _, field := range m.Fields {
		elem := dataElem[field]
		var data string
		if elem.typ == "b" {
			data = raw[:elem.len/4]
			raw = raw[elem.len/4:]
			goto setData
		}

		if elem.fixed {
			data = raw[:elem.len]
			raw = raw[elem.len:]
			goto setData
		}

		{
			elen := 2
			if elem.len >= 100 {
				elen = 3
			}
			var l int64
			if l, err = strconv.ParseInt(raw[:elen], 10, 64); err != nil {
				return
			}
			raw = raw[elen:]
			data = raw[:l]
			raw = raw[l:]
		}

	setData:
		m.Data[field] = Data{string: data}
	}
	return
}

func (m *Message) AddData(field int, value interface{}) (err error) {
	if m.Data == nil {
		m.Data = map[int]Data{}
	}
	data := Data{original: value}
	switch dataElem[field].typ {
	case "n":
		if i, ok := value.(int); ok {
			data.string = strconv.Itoa(i)
		} else {
			return fmt.Errorf("value type is %T; want int", value)
		}
	case "b":
		var binary int64
		if s, ok := value.(string); ok {
			if binary, err = strconv.ParseInt(s, 2, 64); err != nil {
				return
			}
		} else {
			return fmt.Errorf("value type is %T; want string", value)
		}
		data.string = strconv.FormatInt(binary, 16)
	default:
		if s, ok := value.(string); ok {
			data.string = s
		} else {
			return fmt.Errorf("value type is %T; want string", value)
		}
	}
	m.Data[field] = data
	return
}

func (m *Message) Build() (r string) {
	m.genBitMap()
	r += m.MTI
	r += m.RawBitMap
	r += m.packDataElem()
	return
}

func (m *Message) genBitMap() {
	var bitmap [128]int
	m.Fields = []int{}
	for field, _ := range m.Data {
		m.Fields = append(m.Fields, field)
		bitmap[field-1] = 1
		if field > 64 {
			bitmap[0] = 1
		}
	}
	sort.Sort(sort.IntSlice(m.Fields))
	m.BitMap = ""
	for i, bit := range bitmap {
		if i > 64 && bitmap[0] != 1 {
			break
		}
		m.BitMap += strconv.Itoa(bit)
	}
	m.RawBitMap = ""
	for i := 1; i <= len(m.BitMap)/4; i++ {
		data, _ := strconv.ParseInt(m.BitMap[(i-1)*4:i*4], 2, 64)
		m.RawBitMap += strconv.FormatInt(data, 16)
	}
}

func (m *Message) packDataElem() (r string) {
	for _, field := range m.Fields {
		elem := dataElem[field]
		data := m.Data[field]
		if elem.fixed {
			switch elem.typ {
			case "n", "b":
				r += fmt.Sprintf("%0"+fmt.Sprint(elem.len)+"s", data.string)
				// println(field, fmt.Sprintf("%0"+fmt.Sprint(elem.len)+"s", data.string))
			// case "b":
			default:
				r += fmt.Sprintf("% "+fmt.Sprint(elem.len)+"s", data.string)
				// println(field, fmt.Sprintf("% "+fmt.Sprint(elem.len)+"s", data.string))
			}
		} else {
			elen := "2"
			if elem.len > 99 {
				elen = "3"
			}
			r += fmt.Sprintf("%0"+elen+"d", len(data.string)) + data.string
			// println(field, fmt.Sprintf("%0"+elen+"d", len(data.string))+data.string)
		}
	}
	return
}

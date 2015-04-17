package gojak8583

import (
	"fmt"
	"reflect"
	"testing"
)

func TestParse(t *testing.T) {
	input := "0210f22800000a100808000000000000010012255777630950003000000000100000030310370713333720152245273363800007SUCCESS000820150303000002071Txn ID : 224527336380 for Tzs 1000 cash out has being confirmed by ipay007CASHOUT"
	got, err := Parse(input)
	if err != nil {
		t.Error(err)
	}
	want := Message{
		MTI:       "0210",
		BitMap:    "11110010001010000000000000000000000010100001000000001000000010000000000000000000000000000000000000000000000000000000000100000000",
		RawBitMap: "f22800000a1008080000000000000100",
		Fields:    []int{2, 3, 4, 7, 11, 13, 37, 39, 44, 53, 61, 120},
		Data: map[int]Data{
			2:   Data{string: "255777630950"},
			3:   Data{string: "003000"},
			4:   Data{string: "000000100000"}, // "100000"
			7:   Data{string: "0303103707"},   //"303103707"
			11:  Data{string: "133337"},
			13:  Data{string: "2015"},
			37:  Data{string: "224527336380"},
			39:  Data{string: "00"},
			44:  Data{string: "SUCCESS"},
			53:  Data{string: "000820150303000002"}, // "8.20150303E+14"
			61:  Data{string: "Txn ID : 224527336380 for Tzs 1000 cash out has being confirmed by ipay"},
			120: Data{string: "CASHOUT"},
		},
	}
	if !reflect.DeepEqual(got, want) {
		fmt.Printf("got.MTI  = %s\nwant.MTI = %s\n", got.MTI, want.MTI)
		fmt.Printf("got.BitMap  = %s\nwant.BitMap = %s\n", got.BitMap, want.BitMap)
		fmt.Printf("got.RawBitMap  = %s\nwant.RawBitMap = %s\n", got.RawBitMap, want.RawBitMap)
		for _, f := range got.Fields {
			println(f, got.Data[f].string)
		}
		for _, f := range want.Fields {
			println(f, want.Data[f].string)
		}
	}
}

func TestBuild(t *testing.T) {
	var msg Message
	msg.MTI = "0210"
	msg.AddData(2, 255777630950)
	msg.AddData(3, 3000)
	msg.AddData(4, 100000)
	msg.AddData(7, 303103707)
	msg.AddData(11, 133337)
	msg.AddData(13, 2015)
	msg.AddData(37, "224527336380")
	msg.AddData(39, "00")
	msg.AddData(44, "SUCCESS")
	msg.AddData(53, 820150303000002)
	msg.AddData(61, "Txn ID : 224527336380 for Tzs 1000 cash out has being confirmed by ipay")
	msg.AddData(120, "CASHOUT")

	got := msg.Build()
	want := "0210f22800000a100808000000000000010012255777630950003000000000100000030310370713333720152245273363800007SUCCESS000820150303000002071Txn ID : 224527336380 for Tzs 1000 cash out has being confirmed by ipay007CASHOUT"
	if got != want {
		t.Errorf("Message.Build = %s; want %s", got, want)
	}
}

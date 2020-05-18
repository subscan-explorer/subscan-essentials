package types

import (
	"encoding/binary"
	"fmt"
	"reflect"
	"strings"
	"subscan-end/utiles"
)

type IScaleDecoder interface {
	Init(data ScaleBytes, subType string, arg ...interface{})
	Process()
}

type ScaleDecoder struct {
	IScaleDecoder
	TypeString       string            `json:"-"`
	Data             ScaleBytes        `json:"-"`
	RawValue         string            `json:"-"`
	Value            interface{}       `json:"-"`
	SubType          string            `json:"-"`
	TypeMapping      map[string]string `json:"-"`
	StructOrderField []string          `json:"-"`
}

// arg argument
func (s *ScaleDecoder) Init(data ScaleBytes, subType string, arg ...interface{}) {
	s.SubType = subType
	if s.TypeMapping == nil && s.TypeString != "" {
		s.buildTypeMapping()
	}
	s.Data = data
	s.RawValue = ""
	s.Value = nil
}

func (s *ScaleDecoder) Process() {

}

func (s *ScaleDecoder) buildTypeMapping() {
	if s.TypeString != "" && string(s.TypeString[0]) == "(" && string(s.TypeString[len(s.TypeString)-1:]) == ")" {
		s.TypeMapping = make(map[string]string)
		for k, v := range strings.Split(s.TypeString[1:len(s.TypeString)-1], ",") {
			s.TypeMapping[fmt.Sprintf("col%d", k+1)] = strings.TrimSpace(v)
			s.StructOrderField = append(s.StructOrderField, fmt.Sprintf("col%d", k+1))
		}
	}
}

func (s *ScaleDecoder) GetNextBytes(length int) []byte {
	data := s.Data.GetNextBytes(length)
	s.RawValue += utiles.BytesToHex(data)
	return data
}

func (s *ScaleDecoder) GetNextU8() int {
	b := s.GetNextBytes(1)
	data := make([]byte, len(s.Data.Data))
	copy(data, s.Data.Data)
	bs := make([]byte, 4-len(b))
	bs = append(b[:], bs...)
	s.Data.Data = data
	return int(binary.LittleEndian.Uint32(bs))
}

func (s *ScaleDecoder) getNextBool() bool {
	data := s.GetNextBytes(1)
	return utiles.BytesToHex(data) == "01"
}

func (s *ScaleDecoder) String() string {
	if s.Value != nil {
		if reflect.ValueOf(s.Value).Kind() == reflect.String {
			return s.Value.(string)
		}
	}
	return ""
}

func (s *ScaleDecoder) UpdateData(v reflect.Value) {
	s.Data.Offset = int(v.Elem().FieldByName("Data").FieldByName("Offset").Int())
	s.Data.Data = v.Elem().FieldByName("Data").FieldByName("Data").Bytes()
}

func (s *ScaleDecoder) ProcessAndUpdateData(typeString string, args ...string) interface{} {
	v := s.ProcessType(typeString, args...)
	v.MethodByName("Process").Call(nil)
	s.UpdateData(v)
	return v.Elem().FieldByName("Value").Interface()
}

func (s *ScaleDecoder) ProcessType(typeString string, valueList ...string) reflect.Value {
	r := RuntimeType{}
	c, rcvr, subType := r.reg().getDecoderClass(typeString)
	if c == nil {
		panic(fmt.Sprintf("not found decoder class %s", typeString))
	}
	method, _ := c.MethodByName("Init")
	method.Func.Call([]reflect.Value{rcvr, reflect.ValueOf(s.Data), reflect.ValueOf(subType), reflect.ValueOf(valueList)})
	return rcvr
}

func (s *ScaleDecoder) ResetData() {
	s.Data.Data = []byte{}
	s.Data.Offset = 0
}

type ScaleType struct {
	ScaleDecoder
}

type MetadataModuleCall struct {
	ScaleType
	Name string              `json:"name"`
	Args []map[string]string `json:"args"`
	Docs []string            `json:"docs"`
}

func (m *MetadataModuleCall) Process() {
	m.Name = m.ProcessAndUpdateData("Bytes").(string)
	argsValue := m.ProcessAndUpdateData("Vec<MetadataModuleCallArgument>").([]interface{})
	var args []map[string]string
	for _, v := range argsValue {
		args = append(args, v.(map[string]string))
	}
	m.Args = args
	docs := m.ProcessAndUpdateData("Vec<Bytes>").([]interface{})
	for _, v := range docs {
		m.Docs = append(m.Docs, v.(string))
	}
	m.Value = MetadataModuleCall{
		Name: m.Name,
		Args: m.Args,
		Docs: m.Docs,
	}
}

type MetadataModuleCallArgument struct {
	ScaleType
	Name string `json:"name"`
	Type string `json:"type"`
}

func (m *MetadataModuleCallArgument) Process() {
	m.Name = m.ProcessAndUpdateData("Bytes").(string)
	m.Type = ConvertType(m.ProcessAndUpdateData("Bytes").(string))
	m.Value = map[string]string{
		"name": m.Name,
		"type": m.Type,
	}
	CheckCodecType(m.Type)
}

type MetadataModuleEvent struct {
	ScaleType
	Name string   `json:"name"`
	Docs []string `json:"docs"`
	Args []string `json:"args"`
}

func (m *MetadataModuleEvent) Process() {
	m.Name = m.ProcessAndUpdateData("Bytes").(string)
	args := m.ProcessAndUpdateData("Vec<Bytes>").([]interface{})
	for _, v := range args {
		m.Args = append(m.Args, v.(string))
	}
	docs := m.ProcessAndUpdateData("Vec<Bytes>").([]interface{})
	for _, v := range docs {
		m.Docs = append(m.Docs, v.(string))
	}
	m.Value = MetadataEvents{
		Name: m.Name,
		Args: m.Args,
		Docs: m.Docs,
	}
}

// Code generated by github.com/actgardner/gogen-avro/v10. DO NOT EDIT.
/*
 * SOURCE:
 *     event.avsc
 */
 package avro

 import (
	 "encoding/json"
	 "fmt"
	 "io"
 
	 "github.com/actgardner/gogen-avro/v10/compiler"
	 "github.com/actgardner/gogen-avro/v10/vm"
	 "github.com/actgardner/gogen-avro/v10/vm/types"
 )
 
 var _ = fmt.Printf
 
 // Schema to event message.
 type Event struct {
	 // The string is a unicode character sequence.
	 Message string `json:"message"`
 }
 
 const EventAvroCRC64Fingerprint = ";r\x1a\xed\xab+\xca\xe6"
 
 func NewEvent() Event {
	 r := Event{}
	 return r
 }
 
 func DeserializeEvent(r io.Reader) (Event, error) {
	 t := NewEvent()
	 deser, err := compiler.CompileSchemaBytes([]byte(t.Schema()), []byte(t.Schema()))
	 if err != nil {
		 return t, err
	 }
 
	 err = vm.Eval(r, deser, &t)
	 return t, err
 }
 
 func DeserializeEventFromSchema(r io.Reader, schema string) (Event, error) {
	 t := NewEvent()
 
	 deser, err := compiler.CompileSchemaBytes([]byte(schema), []byte(t.Schema()))
	 if err != nil {
		 return t, err
	 }
 
	 err = vm.Eval(r, deser, &t)
	 return t, err
 }
 
 func writeEvent(r Event, w io.Writer) error {
	 var err error
	 err = vm.WriteString(r.Message, w)
	 if err != nil {
		 return err
	 }
	 return err
 }
 
 func (r Event) Serialize(w io.Writer) error {
	 return writeEvent(r, w)
 }
 
 func (r Event) Schema() string {
	 return "{\"doc\":\"Schema to event message.\",\"fields\":[{\"doc\":\"The string is a unicode character sequence.\",\"name\":\"message\",\"type\":\"string\"}],\"name\":\"event.namespace.event\",\"type\":\"record\"}"
 }
 
 func (r Event) SchemaName() string {
	 return "event.namespace.event"
 }
 
 func (_ Event) SetBoolean(v bool)    { panic("Unsupported operation") }
 func (_ Event) SetInt(v int32)       { panic("Unsupported operation") }
 func (_ Event) SetLong(v int64)      { panic("Unsupported operation") }
 func (_ Event) SetFloat(v float32)   { panic("Unsupported operation") }
 func (_ Event) SetDouble(v float64)  { panic("Unsupported operation") }
 func (_ Event) SetBytes(v []byte)    { panic("Unsupported operation") }
 func (_ Event) SetString(v string)   { panic("Unsupported operation") }
 func (_ Event) SetUnionElem(v int64) { panic("Unsupported operation") }
 
 func (r *Event) Get(i int) types.Field {
	 switch i {
	 case 0:
		 w := types.String{Target: &r.Message}
 
		 return w
 
	 }
	 panic("Unknown field index")
 }
 
 func (r *Event) SetDefault(i int) {
	 switch i {
	 }
	 panic("Unknown field index")
 }
 
 func (r *Event) NullField(i int) {
	 switch i {
	 }
	 panic("Not a nullable field index")
 }
 
 func (_ Event) AppendMap(key string) types.Field { panic("Unsupported operation") }
 func (_ Event) AppendArray() types.Field         { panic("Unsupported operation") }
 func (_ Event) HintSize(int)                     { panic("Unsupported operation") }
 func (_ Event) Finalize()                        {}
 
 func (_ Event) AvroCRC64Fingerprint() []byte {
	 return []byte(EventAvroCRC64Fingerprint)
 }
 
 func (r Event) MarshalJSON() ([]byte, error) {
	 var err error
	 output := make(map[string]json.RawMessage)
	 output["message"], err = json.Marshal(r.Message)
	 if err != nil {
		 return nil, err
	 }
	 return json.Marshal(output)
 }
 
 func (r *Event) UnmarshalJSON(data []byte) error {
	 var fields map[string]json.RawMessage
	 if err := json.Unmarshal(data, &fields); err != nil {
		 return err
	 }
 
	 var val json.RawMessage
	 val = func() json.RawMessage {
		 if v, ok := fields["message"]; ok {
			 return v
		 }
		 return nil
	 }()
 
	 if val != nil {
		 if err := json.Unmarshal([]byte(val), &r.Message); err != nil {
			 return err
		 }
	 } else {
		 return fmt.Errorf("no value specified for message")
	 }
	 return nil
 }
 
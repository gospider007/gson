package gson

import (
	"errors"

	"github.com/gospider007/tools"
	jsoniter "github.com/json-iterator/go"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
	"go.mongodb.org/mongo-driver/bson"
)

var jsonConfig = jsoniter.Config{
	CaseSensitive: true,

	EscapeHTML:             false,
	SortMapKeys:            true,
	ValidateJsonRawMessage: true,
}.Froze()

type Client struct {
	g gjson.Result
}
type Map map[string]*Client

func (obj *Client) Map() map[string]*Client {
	result := map[string]*Client{}
	for kk, vv := range obj.g.Map() {
		result[kk] = &Client{g: vv}
	}
	return result
}
func (obj *Client) RawMap() map[string]any {
	result := map[string]any{}
	for kk, vv := range obj.g.Map() {
		result[kk] = vv.Value()
	}
	return result
}

func (obj *Client) Array() []*Client {
	lls := obj.g.Array()
	result := make([]*Client, len(lls))
	for i := 0; i < len(lls); i++ {
		result[i] = &Client{g: lls[i]}
	}
	return result
}
func (obj *Client) Value() any {
	return obj.g.Value()
}
func (obj *Client) String() string {
	return obj.g.String()
}
func (obj *Client) Bytes() []byte {
	return tools.StringToBytes(obj.String())
}
func (obj *Client) Int() int64 {
	return obj.g.Int()
}
func (obj *Client) Raw() string {
	return obj.g.Raw
}
func (obj *Client) Bool() bool {
	return obj.g.Bool()
}
func (obj *Client) Exists() bool {
	return obj.g.Exists()
}
func (obj *Client) IsArray() bool {
	return obj.g.IsArray()
}
func (obj *Client) IsObject() bool {
	return obj.g.IsObject()
}
func (obj *Client) Float() float64 {
	return obj.g.Float()
}
func (obj *Client) Uint() uint64 {
	return obj.g.Uint()
}
func (obj *Client) Get(path string) *Client {
	return &Client{g: obj.g.Get(path)}
}
func (obj *Client) Set(path string, value any) (err error) {
	obj.g.Raw, err = sjson.Set(obj.g.Raw, path, value)
	return
}
func (obj *Client) Delete(path string) (err error) {
	obj.g.Raw, err = sjson.Delete(obj.g.Raw, path)
	return
}
func (obj *Client) Find(path string) (result *Client) {
	obj.ForEach(func(value *Client) bool {
		result = value.Get(path)
		return !result.Exists()
	})
	return
}
func (obj *Client) Finds(path string) []*Client {
	results := []*Client{}
	obj.ForEach(func(value *Client) bool {
		if result := value.Get(path); result.Exists() {
			results = append(results, result)
		}
		return true
	})
	return results
}
func (obj *Client) ForEach(iterator func(value *Client) bool) bool {
	if obj.IsArray() {
		if !iterator(obj) {
			return false
		}
		for _, value := range obj.Array() {
			if !iterator(value) {
				return false
			}
			if !value.ForEach(iterator) {
				return false
			}
		}
	} else if obj.IsObject() {
		if !iterator(obj) {
			return false
		}
		for _, value := range obj.Map() {
			if !iterator(value) {
				return false
			}
			if !value.ForEach(iterator) {
				return false
			}
		}
	}
	return true
}
func (obj *Client) MarshalJSON() ([]byte, error) {
	return Encode(obj.g.Value())
}
func (obj *Client) MarshalBSON() ([]byte, error) {
	return bson.Marshal(obj.g.Value())
}
func Encode(data any) ([]byte, error) {
	switch value := data.(type) {
	case []byte:
		return value, nil
	}
	return jsonConfig.Marshal(data)
}

// 转成struct
func Decode(data any, strus ...any) (client *Client, err error) {
	var ga gjson.Result
	switch val := data.(type) {
	case gjson.Result:
		ga = val
	case *Client:
		ga = val.g
	case []byte:
		ga = gjson.ParseBytes(val)
	case string:
		ga = gjson.Parse(val)
	default:
		var con []byte
		con, err = Encode(data)
		if err == nil {
			ga = gjson.ParseBytes(con)
		}
	}
	client = &Client{g: ga}
	if err != nil {
		return
	} else if !client.IsObject() && !client.IsArray() {
		err = errors.New("不是一个json对象")
		return
	}
	if len(strus) > 0 {
		err = jsonConfig.Unmarshal(client.Bytes(), strus[0])
	}
	return
}

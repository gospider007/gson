package gson

import (
	"errors"

	"gitee.com/baixudong/tools"
	jsoniter "github.com/json-iterator/go"
	"github.com/tidwall/gjson"
	"go.mongodb.org/mongo-driver/bson"
)

var jsonConfig = jsoniter.Config{
	CaseSensitive: true,
}.Froze()

type Client struct {
	g gjson.Result
}

func (obj *Client) Map() map[string]*Client {
	result := map[string]*Client{}
	for kk, vv := range obj.g.Map() {
		result[kk] = &Client{g: vv}
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
func (obj *Client) MarshalJSON() ([]byte, error) {
	return Encode(obj.g.Value())
}
func (obj *Client) MarshalBSON() ([]byte, error) {
	return bson.Marshal(obj.g.Value())
}
func Encode(data any) ([]byte, error) {
	return jsonConfig.Marshal(data)
}

// 转成struct
func Decode(data any, strus ...any) (client *Client, err error) {
	var ga gjson.Result
	switch val := data.(type) {
	case gjson.Result:
		ga = val
	case []byte:
		if len(strus) > 0 {
			err = jsonConfig.Unmarshal(val, strus[0])
		}
		ga = gjson.ParseBytes(val)
	case string:
		if len(strus) > 0 {
			err = jsonConfig.Unmarshal(tools.StringToBytes(val), strus[0])
		}
		ga = gjson.Parse(val)
	default:
		con, err := Encode(data)
		if err != nil {
			return nil, err
		}
		if len(strus) > 0 {
			err = jsonConfig.Unmarshal(con, strus[0])
		}
		ga = gjson.ParseBytes(con)
	}
	if err == nil {
		if !ga.IsObject() && !ga.IsArray() {
			err = errors.New("不是一个json对象")
		}
	}
	client = &Client{g: ga}
	return
}

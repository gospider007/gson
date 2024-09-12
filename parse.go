package gson

import (
	"sort"
	"strings"

	"github.com/gospider007/kinds"
)

// 解析json 中的列表
func (obj *Client) ParseLis() []*Client {
	var lls []*Client
	var ok bool
	if obj.IsObject() {
		lls, ok = obj.parseLis()
		if ok {
			return lls
		}
	} else if obj.IsArray() {
		lls = obj.Array()
	} else {
		return lls
	}
	for len(lls) > 0 {
		zong := []*Client{}
		for _, ll_data := range lls {
			l_datas, ok := ll_data.parseLis()
			if ok {
				return l_datas
			} else {
				zong = append(zong, l_datas...)
			}
		}
		lls = zong
	}
	return lls
}
func (obj *Client) parseLis() ([]*Client, bool) {
	lls := []*Client{}
	lls2 := []*Client{}
	for _, val := range obj.Map() {
		if val.IsArray() {
			zzs := []*Client{}
			ok := true
			start_txt := ""
			for _, va := range val.Array() {
				if va.IsObject() {
					vks := kinds.NewSet[string]()
					for vk := range va.Map() {
						vks.Add(vk)
					}
					tmpData := vks.Array()
					sort.Strings(tmpData)
					now_txt := strings.Join(tmpData, "##")
					if start_txt == "" {
						start_txt = now_txt
					} else if start_txt != now_txt {
						ok = false
					}
					zzs = append(zzs, va)
				}
			}
			if len(zzs) > 1 && ok {
				if len(zzs) > len(lls2) {
					lls2 = zzs
				}
			} else {
				lls = append(lls, zzs...)
			}
		}
	}
	if len(lls2) > 0 {
		return lls2, true
	}
	return lls, false
}

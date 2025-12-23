package imgtag

import (
	"encoding/xml"
	"strings"

	"github.com/samber/lo"
	"github.com/spf13/cast"
)

type Tags struct {
	XMLName  xml.Name `xml:"Categories"`
	Children []*Tag   `xml:"Category"`
}

type Tag struct {
	Value    string `xml:",chardata"`
	Children []*Tag `xml:"Category"`
}

func UnmarshalHTags(d []byte) (*Tags, error) {
	tags := &Tags{}
	err := xml.Unmarshal(d, tags)
	if err != nil {
		return tags, err
	}
	return tags, nil
}

func (t *Tags) StringSlice() []string {
	tags := [][]string{}
	for _, c := range t.Children {
		ts := walkChildren(c)
		tags = append(tags, ts)
	}
	return lo.Uniq(lo.Flatten(tags))
}

func (t *Tags) IsEmpty() bool {
	return len(t.Children) == 0
}

func FlatTags(tags []string) []*Tag {
	t := []*Tag{}
	for _, f := range tags {
		t = append(t, NewFlatTag(f))
	}
	return t
}

func walkChildren(t *Tag) []string {
	if t == nil {
		return []string{}
	}
	tags := []string{t.Value}
	for _, c := range t.Children {
		tags = append(tags, walkChildren(c)...)
	}
	return tags
}

func NewFlatTag(val string) *Tag {
	tag := &Tag{}
	if l, ok := lo.Last(SplitTags(val)); ok {
		tag.Value = l
		return tag
	}
	tag.Value = val
	return tag
}

func cutTag(val string) (string, string, bool) {
	return strings.Cut(val, "/")
}

func splitHTag(v string) []string {
	return strings.Split(v, Separator)
}

func ParseExifTagValue(v any) []string {
	switch p := v.(type) {
	case string:
		return []string{p}
	default:
		return cast.ToStringSlice(p)
	}
}

func SplitTags(tag string) []string {
	return strings.FieldsFunc(tag, SplitFunc)
}

func parseTag(tag string) string {
	if strings.ContainsFunc(tag, SplitFunc) {
		spl := SplitTags(tag)
		return strings.Join(spl, Separator)
	}
	return tag
}

func HasChildren(tag string) bool {
	return strings.ContainsRune(tag, '/')
}

func SplitFunc(c rune) bool {
	return c == '|' || c == '>' || c == '/'
}

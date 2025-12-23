package img

import (
	"slices"
	"strings"
)

//go:generate stringer -type ExifField
type ExifField int

const (
	// title fields
	Title ExifField = iota
	Caption
	// creator fields
	Source
	Byline
	Credit
	Rights
	Copyright
	// caption fields
	UserComment
	Description
	Notes
	ImageDescription
	// tag fields
	CatalogSets
	HierarchicalSubject
	LastKeywordXMP
	Subject
	Keywords
	Categories
	TagsList
	// meta fields
	ImageHeight
	ImageWidth
	MIMEType
	FileTypeExtension
	WebP_Flags
	Duration
	FileName
	SourceFile
)

const (
	Separator = ` > `
	barSep    = `|`
	slashSep  = `/`
	ogTags    = `original`
	uniq      = `uniq`
	hTags     = `hierarchical`
	subject   = `subject`
)

var hTagFields = []ExifField{
	LastKeywordXMP,
	TagsList,
	HierarchicalSubject,
	CatalogSets,
	Keywords,
	Categories,
}

var flatTagFields = []ExifField{
	Subject,
}

var titleFields = []ExifField{
	Title,
	Caption,
}

var captionFields = []ExifField{
	UserComment,
	Description,
	Notes,
	ImageDescription,
}

var creatorFields = []ExifField{
	Byline,
	Credit,
	Copyright,
	Source,
	Rights,
}

var metaFields = []ExifField{
	ImageHeight,
	ImageWidth,
	MIMEType,
	FileTypeExtension,
	WebP_Flags,
	Duration,
	FileName,
	SourceFile,
}

func (f ExifField) Split(tag string) []string {
	return strings.FieldsFunc(tag, SplitFunc)
}

func (f ExifField) Join(t []string) string {
	return strings.Join(t, f.Sep())
}

func (f ExifField) IsHierarchical() bool {
	return slices.Contains(hTagFields, f)
}

func (f ExifField) Sep() string {
	switch f {
	case HierarchicalSubject, CatalogSets:
		return barSep
	case LastKeywordXMP:
		return slashSep
	default:
		return ""
	}
}

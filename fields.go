package imgtag

import "strings"

//go:generate stringer -type ExifField
type ExifField int

const (
	Title ExifField = iota
	// creator fields
	Source
	Byline
	Credit
	Rights
	Copyright
	// caption fields
	Caption
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

var captionFields = []ExifField{
	Caption,
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
	//return strings.Split(tag, f.Sep())
}

func (f ExifField) Join(t []string) string {
	return strings.Join(t, f.Sep())
}

func (f ExifField) IsHierarchical() bool {
	switch f {
	case HierarchicalSubject, CatalogSets, LastKeywordXMP:
		return true
	default:
		return false
	}
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

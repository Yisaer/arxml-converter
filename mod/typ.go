package mod

type DataType struct {
	ShorName string `json:"short_name"`
	Category string `json:"category"`
	*TypReference
	*Array
	*Structure
}

func NewBasicDataType(shortname, category string, ref string) *DataType {
	dt := &DataType{
		ShorName: shortname,
		Category: category,
	}
	trf := &TypReference{
		Ref: ref,
	}
	dt.TypReference = trf
	return dt
}

func NewStringDataType(shortname, category string, stringSize int64) *DataType {
	dt := &DataType{
		ShorName: shortname,
		Category: category,
	}
	trf := &TypReference{
		Ref:        "string",
		StringSize: stringSize,
	}
	dt.TypReference = trf
	return dt
}

type TypReference struct {
	Ref        string `json:"ref"`
	StringSize int64  `json:"string_size"`
}

type Array struct {
	ArraySize int64  `json:"array_size"`
	Inplace   bool   `json:"inplace"`
	RefType   string `json:"ref_type"`
}

type Structure struct {
	STRList []*StructureTypRef `json:"str_list"`
}

type StructureTypRef struct {
	InPlace  bool   `json:"in_place"`
	Ref      string `json:"ref"`
	ShorName string `json:"shor_name"`
}

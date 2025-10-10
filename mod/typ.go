package mod

type DataType struct {
	ShorName string `json:"short_name"`
	Category string `json:"category"`
	*TypReference
	Array *Array
	*Vector
	*Structure
}

func NewArrayDataType(shortname, category, arrayRef string, arraySize int64) *DataType {
	dt := &DataType{
		ShorName: shortname,
		Category: category,
		Array: &Array{
			ArraySize: arraySize,
			RefType:   arrayRef,
		},
	}
	return dt
}

func NewStructureDataType(shortname, category string, s *Structure) *DataType {
	dt := &DataType{
		ShorName: shortname,
		Category: category,
	}
	dt.Structure = s
	return dt
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

type Vector struct {
	RefType string `json:"ref_type"`
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

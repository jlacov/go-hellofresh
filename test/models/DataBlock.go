package models

type DataBlock struct {
	DateKey string
	XVal    float64
	YVal    uint32
	Count   int `csv:"-"`
}

var DataList []DataBlock

func Merge2(a DataBlock, b DataBlock) DataBlock {
	var merge DataBlock
	merge.Count = a.Count + b.Count
	merge.XVal = a.XVal + b.XVal
	merge.YVal = a.YVal + b.YVal
	return merge
}

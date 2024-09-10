package initializers

import "go-hellofresh/test/models"

func InitDataBlock() models.DataBlock {
	var block models.DataBlock
	block.Count = 0
	block.XVal = float64(0)
	block.YVal = uint32(0)
	return block
}

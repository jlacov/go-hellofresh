package controllers

import (
	"bytes"
	"errors"
	"fmt"
	"go-hellofresh/test/initializers"
	"go-hellofresh/test/models"
	"math"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/icza/gog"
)

var DataMapping sync.Map
var m sync.Mutex
var goodRun int = http.StatusAccepted

func PostEvents(c *gin.Context) {

	//Read in csv lines and split into array
	var b bytes.Buffer
	_, err := b.ReadFrom(c.Request.Body)

	if err != nil {
		c.AbortWithError(400, err)
	}

	var lines []string = regexp.MustCompile("\r?\n").Split(b.String(), -1)

	//Validate data in line and build data block
	for _, line := range lines {

		block, status, err := validateAndBuildBlock(line)

		if err != nil {
			if status >= 400 {
				c.AbortWithError(status, err)
			} else {
				fmt.Printf("Value out of range: %s\n", err)
				continue
			}
		}

		//If data in map for that key, update totals + counts otherwise just add
		m.Lock()
		value, ok := DataMapping.Load(block.DateKey)
		if ok {
			var merge models.DataBlock = value.(models.DataBlock)
			block.Count = block.Count + merge.Count
			block.XVal += merge.XVal
			block.YVal += merge.YVal
			DataMapping.Delete(block.DateKey)
		}

		DataMapping.Store(block.DateKey, block)
		m.Unlock()
	}

	c.Writer.WriteHeader(goodRun)
}

func GetStats(c *gin.Context) {
	//find the key to use
	nowKey := getNowKeyInt()
	start := nowKey - 100

	var buffer models.DataBlock = initializers.InitDataBlock()

	m.Lock()

	for i := start; i <= nowKey; i++ {
		value, ok := DataMapping.Load(strconv.FormatInt(i, 10))
		if ok {
			var temp models.DataBlock = value.(models.DataBlock)
			buffer = models.Merge2(buffer, temp)
		}
	}
	m.Unlock()

	divBy := gog.If(buffer.Count == 0, 1, buffer.Count)
	answer := fmt.Sprintf("%d,%.10f,%.10f,%d,%d", buffer.Count, buffer.XVal, buffer.XVal/float64(divBy), int(buffer.YVal), uint64(buffer.YVal)/uint64(divBy))

	c.Writer.WriteHeader(goodRun)
	_, err := c.Writer.WriteString(answer)

	if err != nil {
		c.AbortWithError(http.StatusInternalServerError, err)
	}

}

func validateAndBuildBlock(line string) (models.DataBlock, int, error) {
	//break up line
	var parts []string = strings.Split(line, ",")
	var maxY int64 = math.MaxInt32
	var minY int64 = int64(maxY / 2)

	var block models.DataBlock

	//Try to convert the date stamp from string
	timestamp, err := convStrInt64(parts[0])
	if err != nil {
		return block, http.StatusBadRequest, errors.New("date : Date conversion error")
	}

	//Try to convert the x val to a float
	xVal, err := strconv.ParseFloat(strings.TrimSpace(parts[1]), 64)
	if err != nil {
		return block, http.StatusBadRequest, errors.New("xVal : Unable to convert data")
	}

	//Try to convert the y val to a int64
	yVal, err := convStrInt64(parts[2])
	if err != nil {
		return block, http.StatusBadRequest, errors.New("yVal : Unable to convert data")
	}

	//Build key based on timestamp
	key := buildKey(timestamp)

	//If too old, toss it
	iNow := getNowKeyInt()
	iKey, _ := convStrInt64(key)

	if iKey < iNow-200 {
		return block, http.StatusContinue, errors.New("date : Data too old")
	}

	//If not between 0 - 1, toss it
	if 0 > xVal || 1 < xVal {
		return block, http.StatusContinue, errors.New("x : Data must be between 0 and 1")
	}

	//if not between MaxValue of integer and unsigned 1/2 of Max value, toss it.
	if int64(minY) > yVal || int64(maxY) < yVal {
		return block, http.StatusContinue, errors.New("y : Data must be between " + strconv.FormatInt(minY, 10) + " and " + strconv.FormatInt(maxY, 10))
	}

	//Set block values.
	block.Count = 1
	block.DateKey = key
	block.XVal = xVal
	block.YVal = uint32(yVal)

	return block, http.StatusAccepted, nil
}

func buildKey(milli int64) string {
	//back to nano
	nano := milli * 1e6
	//format YYYYMMDDHHmmss
	theDate := time.Unix(0, nano).Format("20060102150405")
	return theDate
}

func getNowKey() string {
	return buildKey(getNowUnixMilli())
}

func getNowKeyInt() int64 {
	val, _ := convStrInt64(getNowKey())
	return val
}

func getNowUnixMilli() int64 {
	return time.Now().UnixNano() / 1e6
}

func convStrInt64(val string) (int64, error) {
	return strconv.ParseInt(strings.TrimSpace(val), 10, 64)
}

package controllers

import (
	"fmt"
	"go-hellofresh/test/models"
	"io"
	"math"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

var goodTime int64 = getNowUnixMilli() + (3600 * 1000)
var xVal_good float64 = 0.0123124909
var xVal_bad float64 = 1.0093091321
var yVal_good int64 = int64(math.MaxInt32) - int64(19720913)
var yVal_bad int64 = math.MinInt32
var lineFmt string = "%d,%.10f,%d"

var block models.DataBlock
var status int
var err error

func TestValidateAndBuildBlock_PASS(t *testing.T) {

	line := fmt.Sprintf(lineFmt, goodTime, xVal_good, yVal_good)
	block, status, _ = validateAndBuildBlock(line)

	wantStatus := 202

	if status != wantStatus {
		t.Errorf("Status should be %d, but got %d", wantStatus, status)
	}

	if block.XVal != xVal_good {
		t.Errorf("xVal expected as %f, but got %f", xVal_good, block.XVal)
	}

}

func TestValidateAndBuildBlock_FAIL(t *testing.T) {
	line := fmt.Sprintf(lineFmt, goodTime-500000000, xVal_bad, yVal_bad)

	block, status, err = validateAndBuildBlock(line)

	wantStatus := http.StatusContinue
	if status != wantStatus {
		t.Errorf("Status should be %d, but got %d", wantStatus, status)
	}

	if err == nil {
		t.Errorf("Error expected in response, got nil")
	}
}

func TestGetStats(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	GetStats(c)

	if w.Code != 202 {
		b, _ := io.ReadAll(w.Body)
		t.Error(w.Code, string(b))
	}

}

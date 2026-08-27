package credits

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalculateCreditRange(t *testing.T) {
	// 1. Classe Terza (3)
	res3_6 := CalculateCreditRange(3, 6.0, 7, 0, false)
	assert.Equal(t, 7, res3_6.BaseCreditRangeMin)
	assert.Equal(t, 8, res3_6.BaseCreditRangeMax)
	assert.Equal(t, 7, res3_6.SuggestedCredit)

	res3_75 := CalculateCreditRange(3, 7.6, 9, 40, true)
	assert.Equal(t, 9, res3_75.BaseCreditRangeMin)
	assert.Equal(t, 10, res3_75.BaseCreditRangeMax)
	assert.Equal(t, 10, res3_75.SuggestedCredit) // Max of band due to conduct 9 and decimal >= .5

	res3_92 := CalculateCreditRange(3, 9.2, 10, 50, true)
	assert.Equal(t, 11, res3_92.BaseCreditRangeMin)
	assert.Equal(t, 12, res3_92.BaseCreditRangeMax)
	assert.Equal(t, 12, res3_92.SuggestedCredit)

	// 2. Classe Quarta (4)
	res4_6 := CalculateCreditRange(4, 6.0, 7, 0, false)
	assert.Equal(t, 8, res4_6.BaseCreditRangeMin)
	assert.Equal(t, 9, res4_6.BaseCreditRangeMax)

	res4_85 := CalculateCreditRange(4, 8.5, 9, 30, true)
	assert.Equal(t, 11, res4_85.BaseCreditRangeMin)
	assert.Equal(t, 12, res4_85.BaseCreditRangeMax)
	assert.Equal(t, 12, res4_85.SuggestedCredit)

	// 3. Classe Quinta (5)
	res5_6 := CalculateCreditRange(5, 6.0, 6, 0, false)
	assert.Equal(t, 9, res5_6.BaseCreditRangeMin)
	assert.Equal(t, 10, res5_6.BaseCreditRangeMax)

	res5_95 := CalculateCreditRange(5, 9.8, 10, 100, true)
	assert.Equal(t, 14, res5_95.BaseCreditRangeMin)
	assert.Equal(t, 15, res5_95.BaseCreditRangeMax)
	assert.Equal(t, 15, res5_95.SuggestedCredit)

	// Below 6.0
	res_fail := CalculateCreditRange(3, 5.8, 6, 0, false)
	assert.Equal(t, 0, res_fail.BaseCreditRangeMin)
	assert.Equal(t, 0, res_fail.SuggestedCredit)
}

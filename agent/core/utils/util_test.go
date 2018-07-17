package utils

import (
	"os"
	"testing"
)

func TestLimit2Decimal(t *testing.T) {
	number1 := 1.3456789
	number2 := 1.3446789
	number3 := float64(12345)
	if Limit2Decimal(number1) != 1.35 {

		t.Errorf("Limit2Decimal test error,value is %v\n", Limit2Decimal(number1))
	}

	if Limit2Decimal(number2) != 1.34 {

		t.Errorf("Limit2Decimal test error,value is %v\n", Limit2Decimal(number2))
	}

	if Limit2Decimal(number3) != 12345 {

		t.Errorf("Limit2Decimal test error,value is %v\n", Limit2Decimal(number3))
	}
}

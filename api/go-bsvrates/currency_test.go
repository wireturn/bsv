package bsvrates

import (
	"fmt"
	"math"
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestFormatCommas will test the method FormatCommas()
func TestFormatCommas(t *testing.T) {
	t.Parallel()

	var tests = []struct {
		testCase       string
		integer        int
		expectedString string
	}{
		{"zero", 0, "0"},
		{"one", 1, "1"},
		{"no comma", 123, "123"},
		{"thousand", 1234, "1,234"},
		{"ten thousand", 12345, "12,345"},
		{"hundred thousand", 127127, "127,127"},
		{"trillion", 1271271271271, "1,271,271,271,271"},
	}
	for _, test := range tests {
		t.Run(test.testCase, func(t *testing.T) {
			output := FormatCommas(test.integer)
			assert.Equal(t, test.expectedString, output)
		})
	}
}

// BenchmarkFormatCommas benchmarks the method FormatCommas()
func BenchmarkFormatCommas(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = FormatCommas(10000)
	}
}

// ExampleFormatCommas example using FormatCommas()
func ExampleFormatCommas() {
	val := FormatCommas(1000000)
	fmt.Printf("%s", val)
	// Output:1,000,000
}

// TestConvertSatsToBSV will test the method ConvertSatsToBSV()
func TestConvertSatsToBSV(t *testing.T) {
	t.Parallel()

	var tests = []struct {
		testCase     string
		integer      int
		expectedSats float64
	}{
		{"negative 1", -1, -1e-08},
		{"one", 1, 0.00000001},
		{"hundred", 100, 0.00000100},
		{"thousand", 1000, 0.0000100},
		{"ten thousand", 10000, 0.000100},
		{"hundred thousand", 100000, 0.00100},
		{"million", 1000000, 0.0100},
		{"ten million", 10000000, 0.100},
		{"hundred million", SatoshisPerBitcoin, 1.0},
		{"billion", 1000000000, 10},
		{"random amount", 12345678912, 123.45678912},
	}
	for _, test := range tests {
		t.Run(test.testCase, func(t *testing.T) {
			output := ConvertSatsToBSV(test.integer)
			assert.Equal(t, test.expectedSats, output)
		})
	}
}

// BenchmarkConvertSatsToBSV benchmarks the method ConvertSatsToBSV()
func BenchmarkConvertSatsToBSV(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = ConvertSatsToBSV(1000)
	}
}

// ExampleConvertSatsToBSV example using ConvertSatsToBSV()
func ExampleConvertSatsToBSV() {
	val := ConvertSatsToBSV(1001)
	fmt.Printf("%f", val)
	// Output:0.000010
}

// TestConvertPriceToSatoshis will test the method ConvertPriceToSatoshis()
func TestConvertPriceToSatoshis(t *testing.T) {
	t.Parallel()

	t.Run("valid cases", func(t *testing.T) {
		var tests = []struct {
			testCase     string
			currentRate  float64
			amount       float64
			expectedSats int64
		}{
			{"penny", 157.93895102, 0.01, 6332},
			{"ten", 158.18989656, 10, 6321517},
			{"one", 158.18989656, 1, 632152},
			{"one (different rate)", 158.38610459, 1, 631369},
			{"penny (different rate)", 158.38610459, 0.01, 6314},
			{"rate 1 dollar", 100, 1, 1000000},
			{"rate 100 ten cents", 100, 0.10, 100000},
			{"rate 100 one penny", 100, 0.01, 10000},
			{"rate 100 one hundred fifty", 100, 150, 150000000},
			{"rate 100 one hundred", 100, 100, SatoshisPerBitcoin},
			{"rate 1000", 1000, 1, 100000},
			{"rate 10000", 10000, 1, 10000},
			{"rate 100000", 100000, 1, 1000},
			{"rate 1000000", 1000000, 1, 100},
			{"rate 159.703", 159.70371368849, 0.01, 6262},
			{"rate 159.6089", 159.6089, 0.01, 6266},
		}
		for _, test := range tests {
			t.Run(test.testCase, func(t *testing.T) {
				output, err := ConvertPriceToSatoshis(test.currentRate, test.amount)
				assert.NoError(t, err)
				assert.Equal(t, test.expectedSats, output)
			})
		}
	})

	t.Run("invalid cases", func(t *testing.T) {
		var tests = []struct {
			testCase     string
			currentRate  float64
			amount       float64
			expectedSats int64
		}{
			{"rate zero", 0, 1, 0},
			{"amount zero", 1, 0, 0},
			{"math inf", math.Inf(1), 0, 0},
		}
		for _, test := range tests {
			t.Run(test.testCase, func(t *testing.T) {
				output, err := ConvertPriceToSatoshis(test.currentRate, test.amount)
				assert.Error(t, err)
				assert.Equal(t, test.expectedSats, output)
			})
		}
	})
}

// BenchmarkConvertPriceToSatoshis benchmarks the method ConvertPriceToSatoshis()
func BenchmarkConvertPriceToSatoshis(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = ConvertPriceToSatoshis(150, 10)
	}
}

// ExampleConvertPriceToSatoshis example using ConvertPriceToSatoshis()
func ExampleConvertPriceToSatoshis() {
	val, _ := ConvertPriceToSatoshis(150, 1)
	fmt.Printf("%d", val)
	// Output:666667
}

// TestFormatCentsToDollars will test the method FormatCentsToDollars()
func TestFormatCentsToDollars(t *testing.T) {
	t.Parallel()

	var tests = []struct {
		testCase       string
		integer        int
		expectedString string
	}{
		{"zero", 0, "0.00"},
		{"negative", -1, "-0.01"},
		{"dollar twenty seven", 127, "1.27"},
		{"random number", 199276, "1992.76"},
	}
	for _, test := range tests {
		t.Run(test.testCase, func(t *testing.T) {
			output := FormatCentsToDollars(test.integer)
			assert.Equal(t, test.expectedString, output)
		})
	}
}

// BenchmarkFormatCentsToDollars benchmarks the method FormatCentsToDollars()
func BenchmarkFormatCentsToDollars(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = FormatCentsToDollars(1000)
	}
}

// ExampleFormatCentsToDollars example using FormatCentsToDollars()
func ExampleFormatCentsToDollars() {
	val := FormatCentsToDollars(1000)
	fmt.Printf("%s", val)
	// Output:10.00
}

// TestTransformCurrencyToInt will test the method TransformCurrencyToInt()
func TestTransformCurrencyToInt(t *testing.T) {
	t.Parallel()

	t.Run("valid cases", func(t *testing.T) {
		var tests = []struct {
			testCase    string
			decimal     float64
			currency    Currency
			expectedInt int64
		}{
			{"zero", 0, CurrencyDollars, 0},
			{"dollar and change", 1.27, CurrencyDollars, 127},
			{"leading zero", 01.27, CurrencyDollars, 127},
			{"three decimals", 199.272, CurrencyDollars, 19927},
			{"third decimal greater than 5", 199.276, CurrencyDollars, 19928},
			{"ten sats", 0.00000010, CurrencyBitcoin, 10},
			{"thousand sats", 0.000010, CurrencyBitcoin, 1000},
			{"hundred thousand sats", 0.0010, CurrencyBitcoin, 100000},
			{"ten million", 0.10, CurrencyBitcoin, 10000000},
			{"one bitcoin", 1, CurrencyBitcoin, SatoshisPerBitcoin},
		}
		for _, test := range tests {
			t.Run(test.testCase, func(t *testing.T) {
				output, err := TransformCurrencyToInt(test.decimal, test.currency)
				assert.NoError(t, err)
				assert.Equal(t, test.expectedInt, output)
			})
		}
	})

	t.Run("invalid cases", func(t *testing.T) {
		var tests = []struct {
			testCase    string
			decimal     float64
			currency    Currency
			expectedInt int64
		}{
			{"invalid currency", 0.00000010, 123, 0},
		}
		for _, test := range tests {
			t.Run(test.testCase, func(t *testing.T) {
				output, err := TransformCurrencyToInt(test.decimal, test.currency)
				assert.Error(t, err)
				assert.Equal(t, test.expectedInt, output)
			})
		}

		// todo: issue with negative floats (-1.27 = -126)
	})
}

// BenchmarkTransformCurrencyToInt benchmarks the method TransformCurrencyToInt()
func BenchmarkTransformCurrencyToInt(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = TransformCurrencyToInt(10.00, CurrencyDollars)
	}
}

// ExampleTransformCurrencyToInt example using TransformCurrencyToInt()
func ExampleTransformCurrencyToInt() {
	val, _ := TransformCurrencyToInt(10.00, CurrencyDollars)
	fmt.Printf("%d", val)
	// Output:1000
}

// TestConvertFloatToIntBSV will test the method ConvertFloatToIntBSV()
func TestConvertFloatToIntBSV(t *testing.T) {
	t.Parallel()

	var tests = []struct {
		testCase    string
		float       float64
		expectedInt int64
	}{
		{"zero", 0, 0},
		{"all decimal places", 1.123456789, 112345679},
		{"one sat", 0.00000001, 1},
		{"111 sats", 0.00000111, 111},
		{"negative sats", -0.00000111, -111},
		// {math.Inf(1), -111}, // This will produce a panic in decimal package
	}
	for _, test := range tests {
		t.Run(test.testCase, func(t *testing.T) {
			output := ConvertFloatToIntBSV(test.float)
			assert.Equal(t, test.expectedInt, output)
		})
	}
}

// BenchmarkConvertFloatToIntBSV benchmarks the method ConvertFloatToIntBSV()
func BenchmarkConvertFloatToIntBSV(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = ConvertFloatToIntBSV(10.00)
	}
}

// ExampleConvertFloatToIntBSV example using ConvertFloatToIntBSV()
func ExampleConvertFloatToIntBSV() {
	val := ConvertFloatToIntBSV(10.01)
	fmt.Printf("%d", val)
	// Output:1001000000
}

// TestTransformIntToCurrency will test the method TransformIntToCurrency()
func TestTransformIntToCurrency(t *testing.T) {
	t.Parallel()

	t.Run("valid cases", func(t *testing.T) {
		var tests = []struct {
			testCase       string
			integer        int
			currency       Currency
			expectedString string
		}{
			{"zero", 0, CurrencyDollars, "0.00"},
			{"negative", -1, CurrencyDollars, "-0.01"},
			{"dollar and change", 127, CurrencyDollars, "1.27"},
			{"less than 5", 1274, CurrencyDollars, "12.74"},
			{"more than 5", 1276, CurrencyDollars, "12.76"},
			{"large number", 1270000, CurrencyDollars, "12700.00"},
			{"127 sats", 127, CurrencyBitcoin, "0.00000127"},
			{"random number", 123456789123, CurrencyBitcoin, "1234.56789123"},
		}
		for _, test := range tests {
			t.Run(test.testCase, func(t *testing.T) {
				output, err := TransformIntToCurrency(test.integer, test.currency)
				assert.NoError(t, err)
				assert.Equal(t, test.expectedString, output)
			})
		}
	})

	t.Run("invalid cases", func(t *testing.T) {
		var tests = []struct {
			testCase       string
			integer        int
			currency       Currency
			expectedString string
		}{
			{"invalid currency", 111, 123, ""},
		}
		for _, test := range tests {
			t.Run(test.testCase, func(t *testing.T) {
				output, err := TransformIntToCurrency(test.integer, test.currency)
				assert.Error(t, err)
				assert.Equal(t, test.expectedString, output)
			})
		}
	})
}

// BenchmarkTransformIntToCurrency benchmarks the method TransformIntToCurrency()
func BenchmarkTransformIntToCurrency(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = TransformIntToCurrency(1000, CurrencyDollars)
	}
}

// ExampleTransformIntToCurrency example using TransformIntToCurrency()
func ExampleTransformIntToCurrency() {
	val, _ := TransformIntToCurrency(1000, CurrencyDollars)
	fmt.Printf("%s", val)
	// Output:10.00
}

// TestConvertIntToFloatUSD will test the method ConvertIntToFloatUSD()
func TestConvertIntToFloatUSD(t *testing.T) {
	t.Parallel()

	var tests = []struct {
		testCase      string
		integer       uint64
		expectedFloat float64
	}{
		{"zero", 0, 0.00},
		{"one", 1, 0.010000},
		{"ten", 10, 0.10000},
		{"hundred", 100, 1.0},
		{"thousand", 1000, 10.0},
		{"ten thousand", 10000, 100.0},
		{"penny", 10001, 100.01},
		{"ninety nine", 10099, 100.99},
	}
	for _, test := range tests {
		t.Run(test.testCase, func(t *testing.T) {
			output := ConvertIntToFloatUSD(test.integer)
			assert.Equal(t, test.expectedFloat, output)
		})
	}
}

// BenchmarkConvertIntPriceToFloat benchmarks the method ConvertIntToFloatUSD()
func BenchmarkConvertIntToFloatUSD(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = ConvertIntToFloatUSD(10000)
	}
}

// ExampleConvertIntToFloatUSD example using ConvertIntToFloatUSD()
func ExampleConvertIntToFloatUSD() {
	val := ConvertIntToFloatUSD(1000000)
	fmt.Printf("%f", val)
	// Output:10000.000000
}

// TestGetDollarsFromSatoshis will test the method GetDollarsFromSatoshis()
func TestGetDollarsFromSatoshis(t *testing.T) {
	t.Parallel()

	t.Run("valid cases", func(t *testing.T) {
		var tests = []struct {
			testCase        string
			currentRate     float64
			sats            int64
			expectedDollars float64
		}{
			{"zero rate", 0, 10000, 0},
			{"zero sats - 100", 100, 0, 0},
			{"rate 1 penny - 100", 100, 10000, 0.01},
			{"rate 99 pennies - 100", 100, 10000 * 99, 0.99},
			{"rate 1 dollar - 100", 100, 10000 * 100, 1.00},
			{"rate 10 dollars - 100", 100, 10000 * 100 * 10, 10.00},
			{"rate 100 dollars - 100", 100, 10000 * 100 * 100, 100.00},
			{"rate 1k dollars - 100", 100, 10000 * 100 * 100 * 10, 1000.00},
			{"rate 2 pennies - 100", 100, 20000, 0.02},
			{"rate 1/2 penny - 100", 100, 5000, 0.005},
			{"rate ~3 pennies - 150", 150, 20000, 0.03},
			{"rate ~3 pennies - 160", 160, 20000, 0.032},
			{"rate ~3 pennies - 170", 170, 20000, 0.034},
			{"rate ~4 pennies - 210", 210, 20000, 0.042},
			{"rate ~1 pennies - 99", 99, 20000, 0.0198},
			{"rate ~1 pennies - 99.99", 99.99, 20000, 0.019998},
			{"rate ~2 pennies - 217.67", 217.67, 10000, 0.021767},
			{"rate ~1 pennies - 313.33", 313.33, 4333, 0.0135765889},
		}
		for _, test := range tests {
			t.Run(test.testCase, func(t *testing.T) {
				output := GetDollarsFromSatoshis(test.currentRate, test.sats)
				assert.Equal(t, test.expectedDollars, output)
			})
		}
	})

	t.Run("invalid infinite case", func(t *testing.T) {
		assert.Panics(t, func() {
			_ = GetDollarsFromSatoshis(math.Inf(1), 0)
		})
	})
}

// ExampleGetDollarsFromSatoshis example using GetDollarsFromSatoshis()
func ExampleGetDollarsFromSatoshis() {
	dollars := GetDollarsFromSatoshis(100, 10000)
	fmt.Printf("%f", dollars)
	// Output:0.010000
}

// BenchmarkGetDollarsFromSatoshis benchmarks the method GetDollarsFromSatoshis()
func BenchmarkGetDollarsFromSatoshis(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = GetDollarsFromSatoshis(100, 10000)
	}
}

// TestGetCentsFromSatoshis will test the method GetCentsFromSatoshis()
func TestGetCentsFromSatoshis(t *testing.T) {
	t.Parallel()

	t.Run("valid cases", func(t *testing.T) {
		var tests = []struct {
			testCase      string
			currentRate   float64
			sats          int64
			expectedCents int64
		}{
			{"zero rate", 0, 10000, 0},
			{"zero sats", 100, 0, 0},
			{"rate 1 penny - 100", 100, 10000, 1},
			{"rate 99 pennies - 100", 100, 10000 * 99, 99},
			{"rate 1 dollar - 100", 100, 10000 * 100, 100},
			{"rate 10 dollars - 100", 100, 10000 * 100 * 10, 1000},
			{"rate 100 dollars - 100", 100, 10000 * 100 * 100, 10000},
			{"rate 1k dollars - 100", 100, 10000 * 100 * 100 * 10, 100000},
			{"rate 2 pennies - 100", 100, 20000, 2},
			{"rate 1/2 penny - 100", 100, 5000, 0},
			{"rate 4/10 penny - 100", 100, 4000, 0},
			{"rate 1 penny - 99", 99, 10000, 1},
			{"rate 1 penny - 99.99", 99.99, 10000, 1},
			{"rate 1 penny - 98.10", 98.10, 10000, 1},
			{"rate 1 penny - 150", 150, 10000, 1},
			{"rate 1 penny - 160", 160, 10000, 1},
			{"rate 1 penny - 170", 170, 10000, 1},
			{"rate 2 pennies - 200", 200, 10000, 2},
			{"rate 5 pennies - 500", 500, 10000, 5},
			{"rate 3 pennies - 333.33", 333.33, 10000, 3},
		}
		for _, test := range tests {
			t.Run(test.testCase, func(t *testing.T) {
				output := GetCentsFromSatoshis(test.currentRate, test.sats)
				assert.Equal(t, test.expectedCents, output)
			})
		}
	})

	t.Run("invalid infinite case", func(t *testing.T) {
		assert.Panics(t, func() {
			_ = GetCentsFromSatoshis(math.Inf(1), 0)
		})
	})
}

// ExampleGetCentsFromSatoshis example using GetCentsFromSatoshis()
func ExampleGetCentsFromSatoshis() {
	cents := GetCentsFromSatoshis(100, 20000)
	fmt.Printf("%d", cents)
	// Output:2
}

// BenchmarkGetCentsFromSatoshis benchmarks the method GetCentsFromSatoshis()
func BenchmarkGetCentsFromSatoshis(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = GetCentsFromSatoshis(100, 10000)
	}
}

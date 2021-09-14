package bob

import (
	"fmt"
	"testing"

	"github.com/bitcoinschema/go-bitcoin"
	"github.com/libsv/go-bt"
	"github.com/stretchr/testify/assert"
)

// TestNewFromBytes tests for nil case in NewFromBytes()
func TestNewFromBytes(t *testing.T) {

	t.Parallel()

	var (
		// Testing private methods
		tests = []struct {
			inputLine        []byte
			expectedTxString string
			expectedTxHash   string
			expectedNil      bool
			expectedError    bool
		}{
			{
				[]byte(""),
				"01000000000000000000",
				"",
				false,
				true,
			},
			{
				[]byte("invalid-json"),
				"01000000000000000000",
				"",
				false,
				true,
			},
			{
				[]byte(sampleBobTx),
				"0100000001f15a9d3c550c14e12ca066ad09edff31432f1e9f45894ecff5b70c8354c81f3d010000006b483045022100f012c3bd3781091aa8e53cab2ffcb90acced8c65500b41086fd225e48c98c1d702200b8ff117b8ecd2b2d7e95551bc5a1b3bbcca8049864479a28bed9dc842a86804412103ef5bb22964d529c0af748d9a6381432f05298e7a66ed2fe22e7975b1502528a7ffffffff0200000000000000001f006a15e4b880e781afe883bde999a4e58d83e5b9b4e69a970635386135393733b30100000000001976a9149c63715c6d1fa6c61b31d2911516e1c3db3bdfa888ac00000000",
				"207eaadc096849e037b8944df21a8bba6d91d8445848db047c0a3f963121e19d",
				false,
				false,
			},
		}
	)

	// Run tests
	var b *Tx
	var err error
	for _, test := range tests {
		if b, err = NewFromBytes(test.inputLine); err != nil && !test.expectedError {
			t.Errorf("%s Failed: [%v] inputted and error not expected but got: %s", t.Name(), test.inputLine, err.Error())
		} else if err == nil && test.expectedError {
			t.Errorf("%s Failed: [%v] inputted and error was expected", t.Name(), test.inputLine)
		} else if b == nil && !test.expectedNil {
			t.Errorf("%s Failed: [%v] inputted and nil was not expected", t.Name(), test.inputLine)
		} else if b != nil && test.expectedNil {
			t.Errorf("%s Failed: [%v] inputted and nil was expected", t.Name(), test.inputLine)
		} else if b != nil {

			var str string
			str, err = b.ToRawTxString()
			assert.NoError(t, err)
			assert.Equal(t, test.expectedTxString, str)
			assert.Equal(t, test.expectedTxHash, b.Tx.H)
		}
	}
}

// ExampleNewFromBytes example using NewFromBytes()
func ExampleNewFromBytes() {
	b, err := NewFromBytes([]byte(sampleBobTx))
	if err != nil {
		fmt.Printf("error occurred: %s", err.Error())
		return
	}
	fmt.Printf("found tx: %s", b.Tx.H)
	// Output:found tx: 207eaadc096849e037b8944df21a8bba6d91d8445848db047c0a3f963121e19d
}

// BenchmarkNewFromBytes benchmarks the method NewFromBytes()
func BenchmarkNewFromBytes(b *testing.B) {
	tx := []byte(sampleBobTx)
	for i := 0; i < b.N; i++ {
		_, _ = NewFromBytes(tx)
	}
}

// TestNewFromBytesPanic tests for nil case in NewFromBytes()
func TestNewFromBytesPanic(t *testing.T) {
	t.Parallel()

	assert.Panics(t, func() {
		b, err := NewFromBytes([]byte(sampleBobTxBadStrings))
		assert.NoError(t, err)
		assert.NotNil(t, b)
		_, _ = b.ToRawTxString()
	})
}

// TestNewFromString tests for nil case in NewFromString()
func TestNewFromString(t *testing.T) {
	t.Parallel()

	var (
		// Testing private methods
		tests = []struct {
			inputLine        string
			expectedTxString string
			expectedTxHash   string
			expectedNil      bool
			expectedError    bool
		}{
			{
				"",
				"01000000000000000000",
				"",
				false,
				true,
			},
			{
				"invalid-json",
				"01000000000000000000",
				"",
				false,
				true,
			},
			{
				sampleBobTx,
				"0100000001f15a9d3c550c14e12ca066ad09edff31432f1e9f45894ecff5b70c8354c81f3d010000006b483045022100f012c3bd3781091aa8e53cab2ffcb90acced8c65500b41086fd225e48c98c1d702200b8ff117b8ecd2b2d7e95551bc5a1b3bbcca8049864479a28bed9dc842a86804412103ef5bb22964d529c0af748d9a6381432f05298e7a66ed2fe22e7975b1502528a7ffffffff0200000000000000001f006a15e4b880e781afe883bde999a4e58d83e5b9b4e69a970635386135393733b30100000000001976a9149c63715c6d1fa6c61b31d2911516e1c3db3bdfa888ac00000000",
				"207eaadc096849e037b8944df21a8bba6d91d8445848db047c0a3f963121e19d",
				false,
				false,
			},
		}
	)

	// Run tests
	var b *Tx
	var err error
	for _, test := range tests {
		if b, err = NewFromString(test.inputLine); err != nil && !test.expectedError {
			t.Errorf("%s Failed: [%v] inputted and error not expected but got: %s", t.Name(), test.inputLine, err.Error())
		} else if err == nil && test.expectedError {
			t.Errorf("%s Failed: [%v] inputted and error was expected", t.Name(), test.inputLine)
		} else if b == nil && !test.expectedNil {
			t.Errorf("%s Failed: [%v] inputted and nil was not expected", t.Name(), test.inputLine)
		} else if b != nil && test.expectedNil {
			t.Errorf("%s Failed: [%v] inputted and nil was expected", t.Name(), test.inputLine)
		} else if b != nil {

			var str string
			str, err = b.ToRawTxString()
			assert.NoError(t, err)
			assert.Equal(t, test.expectedTxString, str)
			assert.Equal(t, test.expectedTxHash, b.Tx.H)
		}
	}
}

// ExampleNewFromString example using NewFromString()
func ExampleNewFromString() {
	b, err := NewFromString(sampleBobTx)
	if err != nil {
		fmt.Printf("error occurred: %s", err.Error())
		return
	}
	fmt.Printf("found tx: %s", b.Tx.H)
	// Output:found tx: 207eaadc096849e037b8944df21a8bba6d91d8445848db047c0a3f963121e19d
}

// BenchmarkNewFromString benchmarks the method NewFromString()
func BenchmarkNewFromString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = NewFromString(sampleBobTx)
	}
}

// TestNewFromRawTxString tests for nil case in NewFromRawTxString()
func TestNewFromRawTxString(t *testing.T) {
	t.Parallel()

	var (
		// Testing private methods
		tests = []struct {
			inputLine        string
			expectedTxString string
			expectedTxHash   string
			expectedNil      bool
			expectedError    bool
		}{
			{
				"",
				"01000000000000000000",
				"",
				false,
				true,
			},
			{
				"0",
				"01000000000000000000",
				"",
				false,
				true,
			},
			{
				"invalid-tx",
				"01000000000000000000",
				"",
				false,
				true,
			},
			{
				rawBobTx,
				"01000000018f81a0884a11452aa5860f3b0016db1ec58d0cd654b2fa11ebdfd7e87eabeb0e00000000964c948f81a0884a11452aa5860f3b0016db1ec58d0cd654b2fa11ebdfd7e87eabeb0e020000006b483045022100bfbaa9cb07155cd3690722a9d527c70f91a6fc79233b0d091729e457e7c59dd902203059e1f077593654d8f7d2e22a5a40013e8dbf6fcccc5595305144149e5ed9014121039c555f098562d5f6cff2764008d6491961ab51c49356fee349720781ff6dfff7ffffffff00000000030000000000000000fd9d04006a2231394878696756345179427633744870515663554551797131707a5a56646f417574200a746578742f706c61696e04746578740a7477657463682e7478747c223150755161374b36324d694b43747373534c4b79316b683536575755374d74555235035345540b7477646174615f6a736f6e4dbd027b22637265617465645f6174223a22576564204f63742032312031323a30363a3238202b303030302032303230222c227477745f6964223a2231333138383836333639363530303033393639222c2274657874223a2257534a20456469746f7269616c20426f6172643a204a6f6520426964656e204d75737420416e73776572205175657374696f6e732041626f75742048756e74657220426964656e20616e64204368696e612068747470733a2f2f7777772e6272656974626172742e636f6d2f6e6174696f6e616c2d73656375726974792f323032302f31302f32302f77736a2d656469746f7269616c2d626f6172642d6a6f652d626964656e2d6d7573742d616e737765722d7175657374696f6e732d61626f75742d68756e7465722d626964656e2d616e642d6368696e612f2076696120404272656974626172744e657773204a6f6520426964656e206973206120746f74616c6c7920636f727275707420706f6c6974696369616e2c20616e6420676f74206361756768742e204174206c65617374206e6f7720686520776f6ee28099742062652061626c6520746f20726169736520796f7572205461786573202d204269676765737420696e63726561736520696e20552e532e20686973746f727921222c2275736572223a7b226e616d65223a22446f6e616c64204a2e205472756d70222c2273637265656e5f6e616d65223a227265616c446f6e616c645472756d70222c22637265617465645f6174223a22576564204d61722031382031333a34363a3338202b303030302032303039222c227477745f6964223a223235303733383737222c2270726f66696c655f696d6167655f75726c223a22687474703a2f2f7062732e7477696d672e636f6d2f70726f66696c655f696d616765732f3837343237363139373335373539363637322f6b5575687430306d5f6e6f726d616c2e6a7067227d7d0375726c3e68747470733a2f2f747769747465722e636f6d2f7265616c446f6e616c645472756d702f7374617475732f3133313838383633363936353030303339363907636f6d6d656e74046e756c6c076d625f75736572046e756c6c057265706c79046e756c6c047479706504706f73740974696d657374616d70046e756c6c036170700674776574636807696e766f6963652434626130313735632d313738662d346636332d623737662d3536323737313562326563657c22313550636948473232534e4c514a584d6f53556157566937575371633768436676610d424954434f494e5f454344534122313438574448366e465776356748383177657043726b3566486b4a774550415134514c58494531786378574a6b4e364a6538683361426d644161574947487841773333556167515951586539704672794b4a55334f786875324c54646b784b364d4b5675624a4475592f516957743164776f7a782b796167696c553d00000000000000001976a91405186ff0710ed004229e644c0653b2985c648a2388ac00000000000000001976a9142f0fadb49432be5f3d13a7db410e7c2ddae5103188ac00000000",
				"9ec47d91ff11edb62f337dc828c52e39072d1a5a2f1b180bbfae9c3279d81a7c",
				false,
				false,
			},
		}
	)

	// Run tests
	var b *Tx
	var err error
	for _, test := range tests {
		if b, err = NewFromRawTxString(test.inputLine); err != nil && !test.expectedError {
			t.Errorf("%s Failed: [%v] inputted and error not expected but got: %s", t.Name(), test.inputLine, err.Error())
		} else if err == nil && test.expectedError {
			t.Errorf("%s Failed: [%v] inputted and error was expected", t.Name(), test.inputLine)
		} else if b == nil && !test.expectedNil {
			t.Errorf("%s Failed: [%v] inputted and nil was not expected", t.Name(), test.inputLine)
		} else if b != nil && test.expectedNil {
			t.Errorf("%s Failed: [%v] inputted and nil was expected", t.Name(), test.inputLine)
		} else if b != nil {

			var str string
			str, err = b.ToRawTxString()
			assert.NoError(t, err)
			assert.Equal(t, test.expectedTxString, str)
			assert.Equal(t, test.expectedTxHash, b.Tx.H)
		}
	}
}

// ExampleNewFromRawTxString example using NewFromRawTxString()
func ExampleNewFromRawTxString() {
	b, err := NewFromRawTxString(rawBobTx)
	if err != nil {
		fmt.Printf("error occurred: %s", err.Error())
		return
	}
	fmt.Printf("found tx: %s", b.Tx.H)
	// Output:found tx: 9ec47d91ff11edb62f337dc828c52e39072d1a5a2f1b180bbfae9c3279d81a7c
}

// BenchmarkNewFromRawTxString benchmarks the method NewFromRawTxString()
func BenchmarkNewFromRawTxString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = NewFromRawTxString(rawBobTx)
	}
}

func testExampleTx() (*bt.Tx, error) {
	pk := "80699541455b59a8a8a33b85892319de8b8e8944eb8b48e9467137825ae192e59f01"

	privateKey, err := bitcoin.PrivateKeyFromString(pk)
	if err != nil {
		return nil, err
	}

	opReturn1 := bitcoin.OpReturnData{[]byte("prefix1"), []byte("example data"), []byte{0x13, 0x37}, []byte{0x7c}, []byte("prefix2"), []byte("example data 2")}

	return bitcoin.CreateTx(nil, nil, []bitcoin.OpReturnData{opReturn1}, privateKey)
}

// TestNewFromTx tests for nil case in NewFromTx()
func TestNewFromTx(t *testing.T) {
	t.Parallel()

	validTx, exampleErr := testExampleTx()
	assert.NoError(t, exampleErr)
	assert.NotNil(t, validTx)

	var (
		// Testing private methods
		tests = []struct {
			inputTx          *bt.Tx
			expectedTxString string
			expectedTxHash   string
			expectedNil      bool
			expectedError    bool
		}{
			{
				&bt.Tx{},
				"01000000000000000000",
				"f702453dd03b0f055e5437d76128141803984fb10acb85fc3b2184fae2f3fa78",
				false,
				false,
			},
			{
				validTx,
				"010000000001000000000000000032006a07707265666978310c6578616d706c6520646174610213377c07707265666978320e6578616d706c652064617461203200000000",
				"f94e4adeac0cee5e9ff9985373622db9524e9f98d465dc024f85aec8acfeaf16",
				false,
				false,
			},
		}
	)

	// Run tests
	var b *Tx
	var err error
	for _, test := range tests {
		if b, err = NewFromTx(test.inputTx); err != nil && !test.expectedError {
			t.Errorf("%s Failed: [%v] inputted and error not expected but got: %s", t.Name(), test.inputTx, err.Error())
		} else if err == nil && test.expectedError {
			t.Errorf("%s Failed: [%v] inputted and error was expected", t.Name(), test.inputTx)
		} else if b == nil && !test.expectedNil {
			t.Errorf("%s Failed: [%v] inputted and nil was not expected", t.Name(), test.inputTx)
		} else if b != nil && test.expectedNil {
			t.Errorf("%s Failed: [%v] inputted and nil was expected", t.Name(), test.inputTx)
		} else if b != nil {

			var str string
			str, err = b.ToRawTxString()
			assert.NoError(t, err)
			assert.Equal(t, test.expectedTxString, str)
			assert.Equal(t, test.expectedTxHash, b.Tx.H)
		}
	}
}

// ExampleNewFromTx example using NewFromTx()
func ExampleNewFromTx() {
	// Use an example TX
	exampleTx, err := testExampleTx()
	if err != nil {
		fmt.Printf("error occurred: %s", err.Error())
		return
	}

	var b *Tx
	if b, err = NewFromTx(exampleTx); err != nil {
		fmt.Printf("error occurred: %s", err.Error())
		return
	}
	fmt.Printf("found tx: %s", b.Tx.H)
	// Output:found tx: f94e4adeac0cee5e9ff9985373622db9524e9f98d465dc024f85aec8acfeaf16
}

// BenchmarkNewFromTx benchmarks the method NewFromTx()
func BenchmarkNewFromTx(b *testing.B) {
	exampleTx, _ := testExampleTx()
	for i := 0; i < b.N; i++ {
		_, _ = NewFromTx(exampleTx)
	}
}

// TestNewFromTxPanic tests for nil case in NewFromTx()
func TestNewFromTxPanic(t *testing.T) {
	t.Parallel()

	assert.Panics(t, func() {
		b, err := NewFromTx(nil)
		assert.NoError(t, err)
		assert.NotNil(t, b)
		_, _ = b.ToRawTxString()
	})
}

// TestTx_ToTx tests for nil case in ToTx()
func TestTx_ToTx(t *testing.T) {

	bobTx, err := NewFromString(sampleBobTx)
	assert.NoError(t, err)
	assert.NotNil(t, bobTx)

	var tx *bt.Tx
	tx, err = bobTx.ToTx()
	assert.NoError(t, err)
	assert.NotNil(t, tx)
	assert.Equal(t, bobTx.Tx.H, tx.GetTxID())
}

// ExampleTx_ToTx example using ToTx()
func ExampleTx_ToTx() {
	// Use an example TX
	bobTx, err := NewFromString(sampleBobTx)
	if err != nil {
		fmt.Printf("error occurred: %s", err.Error())
		return
	}

	var tx *bt.Tx
	if tx, err = bobTx.ToTx(); err != nil {
		fmt.Printf("error occurred: %s", err.Error())
		return
	}

	fmt.Printf("found tx: %s", tx.GetTxID())
	// Output:found tx: 207eaadc096849e037b8944df21a8bba6d91d8445848db047c0a3f963121e19d
}

// BenchmarkTx_ToTx benchmarks the method ToTx()
func BenchmarkTx_ToTx(b *testing.B) {
	bobTx, _ := NewFromString(sampleBobTx)
	for i := 0; i < b.N; i++ {
		_, _ = bobTx.ToTx()
	}
}

// TestTx_ToRawTxString tests for nil case in ToRawTxString()
func TestTx_ToRawTxString(t *testing.T) {
	bobTx, err := NewFromString(sampleBobTx)
	assert.NoError(t, err)
	assert.NotNil(t, bobTx)

	testTx := "0100000001f15a9d3c550c14e12ca066ad09edff31432f1e9f45894ecff5b70c8354c81f3d010000006b483045022100f012c3bd3781091aa8e53cab2ffcb90acced8c65500b41086fd225e48c98c1d702200b8ff117b8ecd2b2d7e95551bc5a1b3bbcca8049864479a28bed9dc842a86804412103ef5bb22964d529c0af748d9a6381432f05298e7a66ed2fe22e7975b1502528a7ffffffff0200000000000000001f006a15e4b880e781afe883bde999a4e58d83e5b9b4e69a970635386135393733b30100000000001976a9149c63715c6d1fa6c61b31d2911516e1c3db3bdfa888ac00000000"

	var rawTx string
	rawTx, err = bobTx.ToRawTxString()
	assert.NoError(t, err)
	assert.Equal(t, testTx, rawTx)
}

// ExampleTx_ToRawTxString example using ToRawTxString()
func ExampleTx_ToRawTxString() {
	// Use an example TX
	bobTx, err := NewFromString(sampleBobTx)
	if err != nil {
		fmt.Printf("error occurred: %s", err.Error())
		return
	}

	var rawTx string
	if rawTx, err = bobTx.ToRawTxString(); err != nil {
		fmt.Printf("error occurred: %s", err.Error())
		return
	}

	fmt.Printf("found raw tx: %s", rawTx)
	// Output:found raw tx: 0100000001f15a9d3c550c14e12ca066ad09edff31432f1e9f45894ecff5b70c8354c81f3d010000006b483045022100f012c3bd3781091aa8e53cab2ffcb90acced8c65500b41086fd225e48c98c1d702200b8ff117b8ecd2b2d7e95551bc5a1b3bbcca8049864479a28bed9dc842a86804412103ef5bb22964d529c0af748d9a6381432f05298e7a66ed2fe22e7975b1502528a7ffffffff0200000000000000001f006a15e4b880e781afe883bde999a4e58d83e5b9b4e69a970635386135393733b30100000000001976a9149c63715c6d1fa6c61b31d2911516e1c3db3bdfa888ac00000000
}

// BenchmarkTx_ToRawTxString benchmarks the method ToRawTxString()
func BenchmarkTx_ToRawTxString(b *testing.B) {
	bobTx, _ := NewFromString(sampleBobTx)
	for i := 0; i < b.N; i++ {
		_, _ = bobTx.ToRawTxString()
	}
}

// TestTx_ToString tests for nil case in ToString()
func TestTx_ToString(t *testing.T) {

	bobTx := new(Tx)
	err := bobTx.FromRawTxString(rawBobTx)
	assert.NoError(t, err)

	// to string
	var txString string
	txString, err = bobTx.ToString()
	assert.NoError(t, err)

	// make another bob tx from string
	var otherBob *Tx
	otherBob, err = NewFromString(txString)
	assert.NoError(t, err)

	// check txid match
	assert.Equal(t, bobTx.Tx.H, otherBob.Tx.H)
}

// ExampleTx_ToString example using ToString()
func ExampleTx_ToString() {
	// Use an example TX
	bobTx, err := NewFromString(sampleBobTx)
	if err != nil {
		fmt.Printf("error occurred: %s", err.Error())
		return
	}

	var rawTx string
	if rawTx, err = bobTx.ToString(); err != nil {
		fmt.Printf("error occurred: %s", err.Error())
		return
	}

	fmt.Printf("found raw tx: %d (length)", len(rawTx)) // todo: show raw tx if possible
	// Output:found raw tx: 1782 (length)
}

// BenchmarkTx_ToString benchmarks the method ToString()
func BenchmarkTx_ToString(b *testing.B) {
	bobTx, _ := NewFromString(sampleBobTx)
	for i := 0; i < b.N; i++ {
		_, _ = bobTx.ToString()
	}
}

package aip

import (
	"fmt"
	"strings"
	"testing"

	"github.com/bitcoinschema/go-bob"
	"github.com/btcsuite/btcd/txscript"
)

const examplePrivateKey = "54035dd4c7dda99ac473905a3d82f7864322b49bab1ff441cc457183b9bd8abd"
const exampleMessage = "test message"

// TestSign will test the method Sign()
func TestSign(t *testing.T) {
	t.Parallel()

	var (
		// Testing private methods
		tests = []struct {
			inputPrivateKey   string
			inputAlgorithm    Algorithm
			inputMessage      string
			expectedSignature string
			expectedNil       bool
			expectedError     bool
		}{
			{
				"0499f8239bfe10eb0f5e53d543635a423c96529dd85fa4bad42049a0b435ebdd",
				BitcoinECDSA,
				exampleMessage,
				"HOpsJCCkmIOBs8HJIn3Od7aa/SLycQSsZ5QuLvaSlKobYvxpkE5Lcb4fAFLXp1h5pJTEHtm/SZICybovE8AcpiM=",
				false,
				false,
			},
			{
				"80699541455b59a8a8a33b85892319de8b8e8944eb8b48e9467137825ae192e59f01",
				Paymail,
				exampleMessage,
				"GwJif5FKr6rPpKYghZ7F7Rmaw00q5nWo6NiJOuje+6WBOM6qd8MlVnbWr2mrrA+aaDatHxavR8JTckpSfx1RJ1o=",
				false,
				false,
			},
			{
				"",
				BitcoinECDSA,
				exampleMessage,
				"",
				false,
				true,
			},
			{
				"",
				Paymail,
				exampleMessage,
				"",
				false,
				true,
			},
			{
				"80699541455b59a8a8a33b85892319de8b8e8944eb8b48e9467137825ae192e59f01",
				BitcoinECDSA,
				"",
				"G0ttfuC/XKY60ZRFmb12lARUJUekPJ1nD/f0WvJ94bJxT5U3SFHgHaJvAo1r/tVe1E0pMX+IuUxfOWckXdIS7wI=",
				false,
				false,
			},
			{
				"00000",
				BitcoinECDSA,
				"",
				"",
				false,
				true,
			},
			{
				"00000",
				Paymail,
				"",
				"",
				false,
				true,
			},
			{
				"e83385af76b2b1997326b567461fb73dd9c27eab9e1e86d26779f4650c5f2b75",
				BitcoinECDSA,
				exampleMessage,
				"G2m+m3KyHeoWcJF7Sj09bzF+td7QjVw+baFJJ3VCTG4qfaMGXtx8roWprfXi5qP7NihY5lkfWCKCngnodWmG104=",
				false,
				false,
			},
			{
				"e83385af76b2b1997326b567461fb73dd9c27eab9e1e86d26779f4650c5f2b75",
				BitcoinSignedMessage,
				exampleMessage,
				"G2m+m3KyHeoWcJF7Sj09bzF+td7QjVw+baFJJ3VCTG4qfaMGXtx8roWprfXi5qP7NihY5lkfWCKCngnodWmG104=",
				false,
				false,
			},
			{
				"e83385af76b2b1997326b567461fb73dd9c27eab9e1e86d26779f4650c5f2b75",
				Paymail,
				exampleMessage,
				"G2m+m3KyHeoWcJF7Sj09bzF+td7QjVw+baFJJ3VCTG4qfaMGXtx8roWprfXi5qP7NihY5lkfWCKCngnodWmG104=",
				false,
				false,
			},
			{
				"73646673676572676164666764666761646667616466",
				BitcoinECDSA,
				exampleMessage,
				"HIRS8UIWLYMwQUaiIDpe0ivhUqQVyHJg1kgOd/rviQJZWe2EFEI7PQblLaZofG+MjLCMbQLxzlV7DOAuFIdxNUc=",
				false,
				false,
			},
			{
				"73646673676572676164666764666761646667616466",
				BitcoinSignedMessage,
				exampleMessage,
				"HIRS8UIWLYMwQUaiIDpe0ivhUqQVyHJg1kgOd/rviQJZWe2EFEI7PQblLaZofG+MjLCMbQLxzlV7DOAuFIdxNUc=",
				false,
				false,
			},
			{
				"73646673676572676164666764666761646667616466",
				Paymail,
				exampleMessage,
				"HIRS8UIWLYMwQUaiIDpe0ivhUqQVyHJg1kgOd/rviQJZWe2EFEI7PQblLaZofG+MjLCMbQLxzlV7DOAuFIdxNUc=",
				false,
				false,
			},
		}
	)

	// Run tests
	for testNo, test := range tests {
		if a, err := Sign(test.inputPrivateKey, test.inputAlgorithm, test.inputMessage); err != nil && !test.expectedError {
			t.Errorf("%d %s Failed: [%s] [%s] [%s] inputted and error not expected but got: %s", testNo, t.Name(), test.inputPrivateKey, test.inputAlgorithm, test.inputMessage, err.Error())
		} else if err == nil && test.expectedError {
			t.Errorf("%d %s Failed: [%s] [%s] [%s] inputted and error was expected", testNo, t.Name(), test.inputPrivateKey, test.inputAlgorithm, test.inputMessage)
		} else if a == nil && !test.expectedNil {
			t.Errorf("%d %s Failed: [%s] [%s] [%s] inputted and nil was not expected (aip)", testNo, t.Name(), test.inputPrivateKey, test.inputAlgorithm, test.inputMessage)
		} else if a != nil && test.expectedNil {
			t.Errorf("%d %s Failed: [%s] [%s] [%s] inputted and nil was expected (aip)", testNo, t.Name(), test.inputPrivateKey, test.inputAlgorithm, test.inputMessage)
		} else if a != nil && a.Signature != test.expectedSignature {
			t.Errorf("%d %s Failed: [%s] [%s] [%s] inputted and expected [%s] but got [%s]", testNo, t.Name(), test.inputPrivateKey, test.inputAlgorithm, test.inputMessage, test.expectedSignature, a.Signature)
		} else if a != nil && err == nil {
			// Test validation - THIS WILL NOT WORK BECAUSE DATA IS NOT SET
			if _, err = a.Validate(); err != nil {
				t.Errorf("%d %s Failed: [%s] [%s] [%s] inputted and validation failed: %s", testNo, t.Name(), test.inputPrivateKey, test.inputAlgorithm, test.inputMessage, err.Error())
			}
		}
	}
}

// ExampleSign example using Sign()
func ExampleSign() {
	a, err := Sign(examplePrivateKey, BitcoinECDSA, exampleMessage)
	if err != nil {
		fmt.Printf("error occurred: %s", err.Error())
		return
	}
	fmt.Printf("address: %s signature: %s", a.AlgorithmSigningComponent, a.Signature)
	// Output:address: 1Dw6EeFNRZStXTUENRrV9tGUh1rT2hi6YP signature: HNQwm/7FV7S5wzDf4L+HayG8PVhenwgeZ0T5QuNnVGbtSe+7L+Um7lxcrjsj7eMi3N4K1dAOqrVbkESkQfV7odc=
}

// ExampleSign_paymail example using Sign()
func ExampleSign_paymail() {
	a, err := Sign(examplePrivateKey, Paymail, exampleMessage)
	if err != nil {
		fmt.Printf("error occurred: %s", err.Error())
		return
	}
	fmt.Printf("address: %s signature: %s", a.AlgorithmSigningComponent, a.Signature)
	// Output:address: 041b8c93100d35bd448f4646cc4678f278351b439b52b303ea31ec9edb5475e73f36e7ef720509250313fcf1b4c5af0dc7c5efa126efe2c3b7008e6f1487c61f31 signature: HNQwm/7FV7S5wzDf4L+HayG8PVhenwgeZ0T5QuNnVGbtSe+7L+Um7lxcrjsj7eMi3N4K1dAOqrVbkESkQfV7odc=
}

// BenchmarkSign benchmarks the method Sign()
func BenchmarkSign(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = Sign(examplePrivateKey, BitcoinECDSA, exampleMessage)
	}
}

// BenchmarkSign_paymail benchmarks the method Sign()
func BenchmarkSign_paymail(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = Sign(examplePrivateKey, Paymail, exampleMessage)
	}
}

// TestAip_Validate will test the method Validate()
func TestAip_Validate(t *testing.T) {

	t.Parallel()

	var (
		// Testing private methods
		tests = []struct {
			inputAip      *Aip
			expectedValid bool
		}{
			{&Aip{
				Algorithm:                 BitcoinECDSA,
				AlgorithmSigningComponent: "12SsqqYk43kggMBpSvWHwJwR31NsgMePKS",
				Data:                      []string{opReturn, exampleMessage},
				Signature:                 "HOpsJCCkmIOBs8HJIn3Od7aa/SLycQSsZ5QuLvaSlKobYvxpkE5Lcb4fAFLXp1h5pJTEHtm/SZICybovE8AcpiM=",
			}, true},
			{&Aip{
				Algorithm:                 BitcoinECDSA,
				AlgorithmSigningComponent: "12SsqqYk43kggMBpSvWHwJwR31NsgMePKS",
				Data:                      []string{"test message"},
				Signature:                 "HOpsJCCkmIOBs8HJIn3Od7aa/SLycQSsZ5QuLvaSlKobYvxpkE5Lcb4fAFLXp1h5pJTEHtm/SZICybovE8AcpiM=",
			}, false},
			{&Aip{}, false},
			{&Aip{
				Algorithm:                 BitcoinECDSA,
				AlgorithmSigningComponent: "1PASGrpoPtNXYVsWtRn3rR3JoesuZmK1Z5",
				Data:                      []string{exampleMessage},
				Signature:                 "invalid-sig",
			}, false},
			{&Aip{
				Algorithm:                 BitcoinECDSA,
				AlgorithmSigningComponent: "",
				Data:                      []string{exampleMessage},
				Signature:                 "IL1f9X5R//+1X+nBf4alcMe+Fom0Dtv5J4R+LBHiDyHYSt6OZqvuX3tTHwZefg/iXu/lsAScd2ekQci+wtbDyic=",
			}, false},
			{&Aip{
				Algorithm:                 BitcoinECDSA,
				AlgorithmSigningComponent: "invalid-address",
				Data:                      []string{exampleMessage},
				Signature:                 "IL1f9X5R//+1X+nBf4alcMe+Fom0Dtv5J4R+LBHiDyHYSt6OZqvuX3tTHwZefg/iXu/lsAScd2ekQci+wtbDyic=",
			}, false},
			{&Aip{
				Algorithm:                 BitcoinECDSA,
				AlgorithmSigningComponent: "invalid-address",
				Data:                      nil,
				Signature:                 "IL1f9X5R//+1X+nBf4alcMe+Fom0Dtv5J4R+LBHiDyHYSt6OZqvuX3tTHwZefg/iXu/lsAScd2ekQci+wtbDyic=",
			}, false},
			{&Aip{
				Algorithm:                 Paymail,
				AlgorithmSigningComponent: "0",
				Data:                      []string{string(rune(txscript.OP_RETURN)), exampleMessage},
				Signature:                 "HOpsJCCkmIOBs8HJIn3Od7aa/SLycQSsZ5QuLvaSlKobYvxpkE5Lcb4fAFLXp1h5pJTEHtm/SZICybovE8AcpiM=",
			}, false},
		}
	)

	// Run tests
	for _, test := range tests {
		if valid, err := test.inputAip.Validate(); valid && !test.expectedValid {
			t.Errorf("%s Failed: [%v] inputted and was valid but should NOT be valid", t.Name(), test.inputAip)
		} else if !valid && test.expectedValid && err != nil {
			t.Errorf("%s Failed: [%v] inputted and NOT valid but should be valid, error: %s", t.Name(), test.inputAip, err.Error())
		}
	}
}

// TestAip_ValidatePanic tests for nil case in Validate()
func TestAip_ValidatePanic(t *testing.T) {
	t.Parallel()

	defer func() {
		if r := recover(); r == nil {
			t.Errorf("the code did not panic")
		}
	}()

	a := &Aip{}
	if a.Algorithm != "" {
		t.Fatalf("algorithm should be empty")
	}
	a = nil
	if valid, _ := a.Validate(); valid {
		t.Fatalf("should be NOT valid")
	}

}

// ExampleAip_Validate example using Validate()
func ExampleAip_Validate() {
	a, err := Sign(examplePrivateKey, BitcoinECDSA, exampleMessage)
	if err != nil {
		fmt.Printf("error occurred: %s", err.Error())
		return
	}

	if valid, err := a.Validate(); valid {
		fmt.Printf("valid signature: %s", a.Signature)
	} else if err != nil {
		fmt.Printf("signature validation failed: %s", err.Error())
	}
	// Output:valid signature: HNQwm/7FV7S5wzDf4L+HayG8PVhenwgeZ0T5QuNnVGbtSe+7L+Um7lxcrjsj7eMi3N4K1dAOqrVbkESkQfV7odc=
}

// BenchmarkAip_Validate benchmarks the method Validate()
func BenchmarkAip_Validate(b *testing.B) {
	a, _ := Sign(examplePrivateKey, BitcoinECDSA, exampleMessage)
	for i := 0; i < b.N; i++ {
		_, _ = a.Validate()
	}
}

// TestSignOpReturnData tests for nil case in SignOpReturnData(), takes data including the OP_RETURN byte
func TestSignOpReturnData(t *testing.T) {
	t.Parallel()

	var (
		// Testing private methods
		tests = []struct {
			inputPrivateKey   string
			inputAlgorithm    Algorithm
			inputData         [][]byte
			expectedSignature string
			expectedOutput    string
			expectedAipNil    bool
			expectedOutNil    bool
			expectedError     bool
		}{
			{
				"80699541455b59a8a8a33b85892319de8b8e8944eb8b48e9467137825ae192e59f01",
				BitcoinECDSA,
				[][]byte{[]byte(exampleMessage)},
				"GwJif5FKr6rPpKYghZ7F7Rmaw00q5nWo6NiJOuje+6WBOM6qd8MlVnbWr2mrrA+aaDatHxavR8JTckpSfx1RJ1o=",
				"006a0c74657374206d65737361676522313550636948473232534e4c514a584d6f53556157566937575371633768436676610d424954434f494e5f45434453412131553151733836707847724e55796a37673752346d386b3879346b6d78766f756f4c5847774a696635464b72367250704b5967685a374637526d6177303071356e576f364e694a4f756a652b3657424f4d367164384d6c566e625772326d7272412b61614461744878617652384a54636b7053667831524a316f3d",
				false,
				false,
				false,
			},
			{
				"",
				BitcoinECDSA,
				[][]byte{[]byte(exampleMessage)},
				"",
				"",
				false,
				true,
				true,
			},
			{
				"80699541455b59a8a8a33b85892319de8b8e8944eb8b48e9467137825ae192e59f01",
				Paymail,
				[][]byte{[]byte(exampleMessage)},
				"GwJif5FKr6rPpKYghZ7F7Rmaw00q5nWo6NiJOuje+6WBOM6qd8MlVnbWr2mrrA+aaDatHxavR8JTckpSfx1RJ1o=",
				"006a0c74657374206d65737361676522313550636948473232534e4c514a584d6f5355615756693757537163376843667661077061796d61696c4c82303439393332396133303066333338653136343731373538313961666334356435353661366235666533633834313236663634633666353035616537616139333930343261346361633931326335396261663738323534346131626234356632333432613536303334343435656133313233643733393536663731306334333962654c5847774a696635464b72367250704b5967685a374637526d6177303071356e576f364e694a4f756a652b3657424f4d367164384d6c566e625772326d7272412b61614461744878617652384a54636b7053667831524a316f3d",
				false,
				false,
				false,
			},
		}
	)

	// Run tests
	for idx, test := range tests {
		if out, _, a, err := SignOpReturnData(test.inputPrivateKey, test.inputAlgorithm, test.inputData); err != nil && !test.expectedError {
			t.Errorf("%d %s Failed: [%s] [%s] [%v] inputted and error not expected but got: %s", idx, t.Name(), test.inputPrivateKey, test.inputAlgorithm, test.inputData, err.Error())
		} else if err == nil && test.expectedError {
			t.Errorf("%d %s Failed: [%s] [%s] [%v] inputted and error was expected", idx, t.Name(), test.inputPrivateKey, test.inputAlgorithm, test.inputData)
		} else if a == nil && !test.expectedAipNil {
			t.Errorf("%d %s Failed: [%s] [%s] [%v] inputted and nil was not expected (aip)", idx, t.Name(), test.inputPrivateKey, test.inputAlgorithm, test.inputData)
		} else if a != nil && test.expectedAipNil {
			t.Errorf("%d %s Failed: [%s] [%s] [%v] inputted and nil was expected (aip)", idx, t.Name(), test.inputPrivateKey, test.inputAlgorithm, test.inputData)
		} else if out == nil && !test.expectedOutNil {
			t.Errorf("%d %s Failed: [%s] [%s] [%v] inputted and nil was not expected (out)", idx, t.Name(), test.inputPrivateKey, test.inputAlgorithm, test.inputData)
		} else if out != nil && test.expectedOutNil {
			t.Errorf("%d %s Failed: [%s] [%s] [%v] inputted and nil was expected (out)", idx, t.Name(), test.inputPrivateKey, test.inputAlgorithm, test.inputData)
		} else if a != nil && a.Signature != test.expectedSignature {
			t.Errorf("%d %s Failed: [%s] [%s] [%v] inputted and expected signature [%s] but got [%s]", idx, t.Name(), test.inputPrivateKey, test.inputAlgorithm, test.inputData, test.expectedSignature, a.Signature)
		} else if out != nil && out.GetLockingScriptHexString() != test.expectedOutput {
			t.Logf("out %s", out.String())
			t.Errorf("%d %s Failed: [%s] [%s] [%v] inputted and expected output [%s] but got [%s]", idx, t.Name(), test.inputPrivateKey, test.inputAlgorithm, test.inputData, test.expectedOutput, out.GetLockingScriptHexString())
		}
	}
}

// ExampleSignOpReturnData example using SignOpReturnData()
func ExampleSignOpReturnData() {
	out, _, a, err := SignOpReturnData(examplePrivateKey, BitcoinECDSA, [][]byte{[]byte("some op_return data")})
	if err != nil {
		fmt.Printf("error occurred: %s", err.Error())
		return
	}
	fmt.Printf("signature: %s output: %s", a.Signature, out.GetLockingScriptHexString())
	// Output:signature: G7zptA7IbNaa7PQlblH1v5ElaOj3Zo49oiUrDMqfWM4QFNdIKDnXMkxLU1YgrxODd8uFNU279utUCC4MGPp5pjM= output: 006a13736f6d65206f705f72657475726e206461746122313550636948473232534e4c514a584d6f53556157566937575371633768436676610d424954434f494e5f454344534122314477364565464e525a5374585455454e52725639744755683172543268693659504c5847377a7074413749624e61613750516c626c48317635456c614f6a335a6f34396f695572444d7166574d3451464e64494b446e584d6b784c5531596772784f44643875464e553237397574554343344d47507035706a4d3d
}

// BenchmarkSignOpReturnData benchmarks the method SignOpReturnData()
func BenchmarkSignOpReturnData(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _, _, _ = SignOpReturnData(examplePrivateKey, BitcoinECDSA, [][]byte{[]byte("some op_return data")})
	}
}

func TestFromRawTx(t *testing.T) {

	tx, err := bob.NewFromRawTxString("01000000057d832704c486d74ae7a8960fd886464b1d9c0994412f9e807f2d03c1f6681f14020" +
		"000006a4730440220159bafed7ce0ec4bb39ad9c0fd3391d9ff19473f2419fe831892609c276ac9d702204b9b83d716280b48f4eb8c28a" +
		"ca3a5ce1461751bbaa82937b053276232e719a14121030a8844393b9910a573c4b5dd50376012388cda10ab2e03124b98b2fd954842d2" +
		"ffffffff7ec765003e78cb00caedbb4a11c79c2f3621931c0cc7e9941013b50b3c58ca4c050000006a47304402202813d834d159af027" +
		"8eb86b8506d8e001c6e5d6ce52c91e96f9b7c814381a62502200493bbb06560a24c8439e2a74ffcdd87be50338a638f2c2339d190710b" +
		"e71ff34121021470afef55ce278292015f384ac63a128f6bd0a6ca270b3a1bff7319731027efffffffff18630db3b58026616a208d5a2" +
		"450e8a4d1c5745b3d150e31bbfca954a2d64688030000006b483045022100ecc8c8bc635d7774765ceff31bf1d3f22cbdfa4c2477b3b2" +
		"82b38fd24c9bc66802201bfa43ddb0dfd864beeca4434c4aa5207b763be0af930ef0b20533db4e1b0ee3412103e72133552fff34fdd60" +
		"90409392065328f258299c6e5f2fc013d6089b015cfd2ffffffff18630db3b58026616a208d5a2450e8a4d1c5745b3d150e31bbfca954" +
		"a2d64688050000006a4730440220378bb4cb3448a707d42c40286e4c259011e42be8544b2decfd28ab95bfd2e2cb02205137d49d06468" +
		"d464ae59d5bf631912d4aa161a6cd55b4432324107b20b6c9394121024d853a0b84e261a0621b8d397a6525254999cee19d87295f59b0" +
		"4747f6769ef0ffffffff7abe2b4a292a027afa2c827d04dda67ab1558873366f5082a32f507c056a91c8040000006a47304402200b7ae" +
		"732deb10db268ee9a0a912ca1f2ad546d84497b4d43ff9874ff5e8bf52a022066789de229c02389b3ed263895174c6c45f586411e5af9" +
		"a9e79923eed6be1eff41210334acf1d2bfc5a5faddd400619f31a5b747e800cee265a1c7c7d62700682d59aaffffffff04000000000000" +
		"0000fd4802006a2231394878696756345179427633744870515663554551797131707a5a56646f4175744c9d477265617420636f6e746" +
		"56e742c2067726561742067726170686963732e2e2e2e6a757374206d697373696e67206f6e65207468696e67210a68747470733a2f2f" +
		"7777772e626974636f696e66696c65732e6f72672f742f633630303465363165316538353833633831363366336437363032616530313" +
		"0646262626131313734663430643564323235323037343962393737653466663320246f73670a746578742f706c61696e04746578740a" +
		"7477657463682e747874017c223150755161374b36324d694b43747373534c4b79316b683536575755374d74555235035345540b74776" +
		"46174615f6a736f6e046e756c6c0375726c046e756c6c07636f6d6d656e74046e756c6c076d625f75736572046e756c6c057265706c79" +
		"046e756c6c047479706504706f73740974696d657374616d70046e756c6c036170700674776574636807696e766f69636524333465633" +
		"96131642d653531312d343862652d386132322d356531346133353037623666017c22313550636948473232534e4c514a584d6f535561" +
		"57566937575371633768436676610d424954434f494e5f45434453412231376e326d4b64346b6267554255656a75396972436d4e464c4" +
		"e484c6862694669684c58494561672f73554a3249443445417378524c4c784c545549714367364f4c50393150792f6468437737797570" +
		"4274535a777451354346596f6b583474355a4732356e664a2f38786a7755356e597962654f2b36387964593d2f080000000000001976a" +
		"914a7878004efceb5abecac1a4132735fc52e9b388888ac16160000000000001976a9141d228ee969f1ba36e3fa4c31298c2678f06cd7" +
		"7488ac7f260000000000001976a91405186ff0710ed004229e644c0653b2985c648a2388ac00000000")
	if err != nil {
		t.Fatalf("error occurred: %s", err.Error())
	}

	a := NewFromTapes(tx.Out[0].Tape)

	// t.Log(a)
	t.Log(a.Algorithm)
	t.Log(a.AlgorithmSigningComponent)
	t.Log(a.Signature)
	t.Log(len(strings.Join(a.Data, "")))
	// t.Log(len(a.Data))

	t.Log(a.Validate())
}

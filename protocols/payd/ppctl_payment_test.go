package gopayd

import (
	"testing"

	"github.com/matryer/is"
	"gopkg.in/guregu/null.v3"
)

func TestCreatePayment_Validate(t *testing.T) {
	is := is.New(t)
	tests := map[string]struct {
		req CreatePayment
		exp string
	}{
		"valid request should return no errors": {
			req: CreatePayment{
				Transaction: "0200000004c4b8372f640f9fab1dc2c14eda6a9669d13ca0f4fff42c318f388cf917399fa9000000004847304402203f2c94003474010010a11cdc4bfac3065e117b22ff1e218fb31230be12a80d5202205b69e27a1815a7d6668a5b73e57b15a6117c94b15b3d915ff3304803e233af5341feffffff417e443a9da68f5bea767bb90f09737df50ff7592d662407dc16ed17af0b821d000000006a47304402200fe1bb41b168aa1e071b39c1bd00d7f960d98406b36c76cbeff98acbe20c117902205628cf5755676f85b2cd360406fc771ed3244395d2cd2bf2292e06e0a8f7e4dc412103b811b71802653c97388faa8a7275a49a2742896285515fb01e2801948ee9cc4cfeffffff94b976366984846918b8ef346da50db6231dcf870c6d48754a98976b3a989c23000000004847304402201baa75b71f066eaa5297efaa878f215fd08e3132e3de2d5c7038e8433ef49cf8022044655ef242869210ed8a9a290c5ccc7cfa70a0d6b8cc7d6dc832d1d728ef106341feffffff4383ff843f365a8c9a6ce44ba1c584840125227e7ad06409f7194423ca614aff000000006a4730440220328b446736fa1a47e8675e7ea31a86f6025ece36aa2e158e21e85758a1cf1db8022073cf6f9f3353337a537bfbfef818497941b6f00f6918d40e87d06751610e739e412102065bd35d20f59e1c8c1254690254f14e40710409481320df3854bbfc867b4698feffffff027a898400000000001976a914fc54fbfac51db40cd845ebe6d243d6c950f4bf4088ac0065cd1d000000001976a914ba903fcaa03a280a9577da32db79e52373b8d0e388ac1b040000",
				MerchantData: MerchantData{
					PaymentReference: "abc123",
				},
				RefundTo: null.String{},
				Memo:     "test this please",
			},
			exp: "no validation errors",
		}, "transaction with invalid prefix should error": {
			req: CreatePayment{
				Transaction: "0x74657374696e672074657374696e67",
				MerchantData: MerchantData{
					PaymentReference: "abc123",
				},
				Memo: "test this please",
			},
			exp: "[transaction: not a valid bitcoin transaction]",
		}, "transaction with invalid hex should error": {
			req: CreatePayment{
				Transaction: "74657374696e672074657374696ezz67",
				MerchantData: MerchantData{
					PaymentReference: "abc123",
				},
				Memo: "test this please",
			},
			exp: "[transaction: not a valid bitcoin transaction]",
		}, "merchant data missing payment reference should error": {
			req: CreatePayment{
				Transaction:  "0200000004c4b8372f640f9fab1dc2c14eda6a9669d13ca0f4fff42c318f388cf917399fa9000000004847304402203f2c94003474010010a11cdc4bfac3065e117b22ff1e218fb31230be12a80d5202205b69e27a1815a7d6668a5b73e57b15a6117c94b15b3d915ff3304803e233af5341feffffff417e443a9da68f5bea767bb90f09737df50ff7592d662407dc16ed17af0b821d000000006a47304402200fe1bb41b168aa1e071b39c1bd00d7f960d98406b36c76cbeff98acbe20c117902205628cf5755676f85b2cd360406fc771ed3244395d2cd2bf2292e06e0a8f7e4dc412103b811b71802653c97388faa8a7275a49a2742896285515fb01e2801948ee9cc4cfeffffff94b976366984846918b8ef346da50db6231dcf870c6d48754a98976b3a989c23000000004847304402201baa75b71f066eaa5297efaa878f215fd08e3132e3de2d5c7038e8433ef49cf8022044655ef242869210ed8a9a290c5ccc7cfa70a0d6b8cc7d6dc832d1d728ef106341feffffff4383ff843f365a8c9a6ce44ba1c584840125227e7ad06409f7194423ca614aff000000006a4730440220328b446736fa1a47e8675e7ea31a86f6025ece36aa2e158e21e85758a1cf1db8022073cf6f9f3353337a537bfbfef818497941b6f00f6918d40e87d06751610e739e412102065bd35d20f59e1c8c1254690254f14e40710409481320df3854bbfc867b4698feffffff027a898400000000001976a914fc54fbfac51db40cd845ebe6d243d6c950f4bf4088ac0065cd1d000000001976a914ba903fcaa03a280a9577da32db79e52373b8d0e388ac1b040000",
				MerchantData: MerchantData{},
				Memo:         "test this please",
			},
			exp: "[merchantData.PaymentReference: value cannot be empty]",
		}, "refundTo too long should error": {
			req: CreatePayment{
				Transaction: "0200000004c4b8372f640f9fab1dc2c14eda6a9669d13ca0f4fff42c318f388cf917399fa9000000004847304402203f2c94003474010010a11cdc4bfac3065e117b22ff1e218fb31230be12a80d5202205b69e27a1815a7d6668a5b73e57b15a6117c94b15b3d915ff3304803e233af5341feffffff417e443a9da68f5bea767bb90f09737df50ff7592d662407dc16ed17af0b821d000000006a47304402200fe1bb41b168aa1e071b39c1bd00d7f960d98406b36c76cbeff98acbe20c117902205628cf5755676f85b2cd360406fc771ed3244395d2cd2bf2292e06e0a8f7e4dc412103b811b71802653c97388faa8a7275a49a2742896285515fb01e2801948ee9cc4cfeffffff94b976366984846918b8ef346da50db6231dcf870c6d48754a98976b3a989c23000000004847304402201baa75b71f066eaa5297efaa878f215fd08e3132e3de2d5c7038e8433ef49cf8022044655ef242869210ed8a9a290c5ccc7cfa70a0d6b8cc7d6dc832d1d728ef106341feffffff4383ff843f365a8c9a6ce44ba1c584840125227e7ad06409f7194423ca614aff000000006a4730440220328b446736fa1a47e8675e7ea31a86f6025ece36aa2e158e21e85758a1cf1db8022073cf6f9f3353337a537bfbfef818497941b6f00f6918d40e87d06751610e739e412102065bd35d20f59e1c8c1254690254f14e40710409481320df3854bbfc867b4698feffffff027a898400000000001976a914fc54fbfac51db40cd845ebe6d243d6c950f4bf4088ac0065cd1d000000001976a914ba903fcaa03a280a9577da32db79e52373b8d0e388ac1b040000",
				RefundTo: func() null.String {
					bb := make([]byte, 0)
					// generate string 1 more byte than 10000
					for i := 0; i <= 100; i++ {
						bb = append(bb, 42)
					}
					return null.StringFrom(string(bb))
				}(),
				MerchantData: MerchantData{
					PaymentReference: "abc123",
				},
				Memo: "test this please",
			},
			exp: "[refundTo: value must be between 0 and 100 characters]",
		},
	}
	for name, test := range tests {
		t.Run(name, func(t *testing.T) {
			is.NewRelaxed(t)
			is.Equal(test.req.Validate() != nil, true)
			is.Equal(test.exp, test.req.Validate().Error())
		})
	}
}

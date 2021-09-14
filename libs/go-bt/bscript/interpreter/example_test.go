package interpreter_test

import (
	"fmt"

	"context"

	"github.com/libsv/go-bt/v2"
	"github.com/libsv/go-bt/v2/bscript/interpreter"
	"golang.org/x/sync/errgroup"
)

func ExampleEngine_Execute() {
	tx, err := bt.NewTxFromString("0200000003a9bc457fdc6a54d99300fb137b23714d860c350a9d19ff0f571e694a419ff3a0010000006b48304502210086c83beb2b2663e4709a583d261d75be538aedcafa7766bd983e5c8db2f8b2fc02201a88b178624ab0ad1748b37c875f885930166237c88f5af78ee4e61d337f935f412103e8be830d98bb3b007a0343ee5c36daa48796ae8bb57946b1e87378ad6e8a090dfeffffff0092bb9a47e27bf64fc98f557c530c04d9ac25e2f2a8b600e92a0b1ae7c89c20010000006b483045022100f06b3db1c0a11af348401f9cebe10ae2659d6e766a9dcd9e3a04690ba10a160f02203f7fbd7dfcfc70863aface1a306fcc91bbadf6bc884c21a55ef0d32bd6b088c8412103e8be830d98bb3b007a0343ee5c36daa48796ae8bb57946b1e87378ad6e8a090dfeffffff9d0d4554fa692420a0830ca614b6c60f1bf8eaaa21afca4aa8c99fb052d9f398000000006b483045022100d920f2290548e92a6235f8b2513b7f693a64a0d3fa699f81a034f4b4608ff82f0220767d7d98025aff3c7bd5f2a66aab6a824f5990392e6489aae1e1ae3472d8dffb412103e8be830d98bb3b007a0343ee5c36daa48796ae8bb57946b1e87378ad6e8a090dfeffffff02807c814a000000001976a9143a6bf34ebfcf30e8541bbb33a7882845e5a29cb488ac76b0e60e000000001976a914bd492b67f90cb85918494767ebb23102c4f06b7088ac67000000")
	if err != nil {
		fmt.Println(err)
		return
	}

	prevTx, err := bt.NewTxFromString("0200000001424408c9d997772e56112c731b6dc6f050cb3847c5570cea12f30bfbc7df0a010000000049483045022100fe759b2cd7f25bce4fcda4c8366891b0d9289dc5bac1cf216909c89dc324437a02204aa590b6e82764971df4fe741adf41ece4cde607cb6443edceba831060213d3641feffffff02408c380c010000001976a914f761fc0927a43f4fab5740ef39f05b1fb7786f5288ac0065cd1d000000001976a914805096c5167877a5799977d46fb9dee5891dc3cb88ac66000000")
	if err != nil {
		fmt.Println(err)
		return
	}

	inputIdx := 0
	input := tx.InputIdx(inputIdx)
	prevOutput := prevTx.OutputIdx(int(input.PreviousTxOutIndex))

	inputASM, err := input.UnlockingScript.ToASM()
	if err != nil {
		fmt.Println(err)
		return
	}

	outputASM, err := prevOutput.LockingScript.ToASM()
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println(inputASM)
	fmt.Println(outputASM)

	if err := interpreter.NewEngine().Execute(interpreter.ExecutionParams{
		Tx:            tx,
		InputIdx:      inputIdx,
		PreviousTxOut: prevOutput,
		Flags:         interpreter.ScriptEnableSighashForkID | interpreter.ScriptUTXOAfterGenesis,
	}); err != nil {
		fmt.Println(err)
		return
	}
	// Output:
	// 304502210086c83beb2b2663e4709a583d261d75be538aedcafa7766bd983e5c8db2f8b2fc02201a88b178624ab0ad1748b37c875f885930166237c88f5af78ee4e61d337f935f41 03e8be830d98bb3b007a0343ee5c36daa48796ae8bb57946b1e87378ad6e8a090d
	// OP_DUP OP_HASH160 805096c5167877a5799977d46fb9dee5891dc3cb OP_EQUALVERIFY OP_CHECKSIG
}

func ExampleEngine_Execute_error() {
	tx, err := bt.NewTxFromString("0200000003a9bc457fdc6a54d99300fb137b23714d860c350a9d19ff0f571e694a419ff3a0010000006b48304502210086c83beb2b2663e4709a583d261d75be538aedcafa7766bd983e5c8db2f8b2fc02201a88b178624ab0ad1748b37c875f885930166237c88f5af78ee4e61d337f935f412103e8be830d98bb3b007a0343ee5c36daa48796ae8bb57946b1e87378ad6e8a090dfeffffff0092bb9a47e27bf64fc98f557c530c04d9ac25e2f2a8b600e92a0b1ae7c89c20010000006b483045022100f06b3db1c0a11af348401f9cebe10ae2659d6e766a9dcd9e3a04690ba10a160f02203f7fbd7dfcfc70863aface1a306fcc91bbadf6bc884c21a55ef0d32bd6b088c8412103e8be830d98bb3b007a0343ee5c36daa48796ae8bb57946b1e87378ad6e8a090dfeffffff9d0d4554fa692420a0830ca614b6c60f1bf8eaaa21afca4aa8c99fb052d9f398000000006b483045022100d920f2290548e92a6235f8b2513b7f693a64a0d3fa699f81a034f4b4608ff82f0220767d7d98025aff3c7bd5f2a66aab6a824f5990392e6489aae1e1ae3472d8dffb412103e8be830d98bb3b007a0343ee5c36daa48796ae8bb57946b1e87378ad6e8a090dfeffffff02807c814a000000001976a9143a6bf34ebfcf30e8541bbb33a7882845e5a29cb488ac76b0e60e000000001976a914bd492b67f90cb85918494767ebb23102c4f06b7088ac67000000")
	if err != nil {
		fmt.Println(err)
		return
	}

	prevTx, err := bt.NewTxFromString("0200000001424408c9d997772e56112c731b6dc6f050cb3847c5570cea12f30bfbc7df0a010000000049483045022100fe759b2cd7f25bce4fcda4c8366891b0d9289dc5bac1cf216909c89dc324437a02204aa590b6e82764971df4fe741adf41ece4cde607cb6443edceba831060213d3641feffffff02408c380c010000001976a914f761fc0927a43f4fab5740ef39f05b1fb7786f5288ac0065cd1d000000001976a914805096c5167877a5799977d46fb9dee5891dc3cb88ac66000000")
	if err != nil {
		fmt.Println(err)
		return
	}

	// Use incorrect output for input
	inputIdx := 1
	prevOutput := prevTx.OutputIdx(0)

	if err := interpreter.NewEngine().Execute(interpreter.ExecutionParams{
		Tx:            tx,
		InputIdx:      inputIdx,
		PreviousTxOut: prevOutput,
		Flags:         interpreter.ScriptEnableSighashForkID | interpreter.ScriptUTXOAfterGenesis,
	}); err != nil {
		fmt.Println(err)
		return
	}
	// Output:
	// OP_EQUALVERIFY failed
}

func ExampleEngine_Execute_concurrent() {
	var params []interpreter.ExecutionParams

	tx, err := bt.NewTxFromString("0200000003a9bc457fdc6a54d99300fb137b23714d860c350a9d19ff0f571e694a419ff3a0010000006b48304502210086c83beb2b2663e4709a583d261d75be538aedcafa7766bd983e5c8db2f8b2fc02201a88b178624ab0ad1748b37c875f885930166237c88f5af78ee4e61d337f935f412103e8be830d98bb3b007a0343ee5c36daa48796ae8bb57946b1e87378ad6e8a090dfeffffff0092bb9a47e27bf64fc98f557c530c04d9ac25e2f2a8b600e92a0b1ae7c89c20010000006b483045022100f06b3db1c0a11af348401f9cebe10ae2659d6e766a9dcd9e3a04690ba10a160f02203f7fbd7dfcfc70863aface1a306fcc91bbadf6bc884c21a55ef0d32bd6b088c8412103e8be830d98bb3b007a0343ee5c36daa48796ae8bb57946b1e87378ad6e8a090dfeffffff9d0d4554fa692420a0830ca614b6c60f1bf8eaaa21afca4aa8c99fb052d9f398000000006b483045022100d920f2290548e92a6235f8b2513b7f693a64a0d3fa699f81a034f4b4608ff82f0220767d7d98025aff3c7bd5f2a66aab6a824f5990392e6489aae1e1ae3472d8dffb412103e8be830d98bb3b007a0343ee5c36daa48796ae8bb57946b1e87378ad6e8a090dfeffffff02807c814a000000001976a9143a6bf34ebfcf30e8541bbb33a7882845e5a29cb488ac76b0e60e000000001976a914bd492b67f90cb85918494767ebb23102c4f06b7088ac67000000")
	if err != nil {
		fmt.Println(err)
		return
	}

	prevTx, err := bt.NewTxFromString("0200000001424408c9d997772e56112c731b6dc6f050cb3847c5570cea12f30bfbc7df0a010000000049483045022100fe759b2cd7f25bce4fcda4c8366891b0d9289dc5bac1cf216909c89dc324437a02204aa590b6e82764971df4fe741adf41ece4cde607cb6443edceba831060213d3641feffffff02408c380c010000001976a914f761fc0927a43f4fab5740ef39f05b1fb7786f5288ac0065cd1d000000001976a914805096c5167877a5799977d46fb9dee5891dc3cb88ac66000000")
	if err != nil {
		fmt.Println(err)
		return
	}

	params = append(params, interpreter.ExecutionParams{
		Tx:            tx,
		InputIdx:      0,
		PreviousTxOut: prevTx.OutputIdx(int(tx.InputIdx(0).PreviousTxOutIndex)),
	})

	tx2, err := bt.NewTxFromString("020000000532bc3895b35a4d7b2da0103589a320e4eabeed08ef9777481b6f2475c0cf0084010000006a47304402206579610b3a845e7ffa58203c686ca86ed3f2f946454bcb5f78e960c8ec34617702206cf0f168267acbca0acdc7fe38311fd94fd821868891aa1da150fe0de6e0ff6c412103bb0164c11476e32287120301be5aca1310b0f72579f83e88cf6e10e42f6f78f1feffffff46987a5d7920f32aa950c9cd258fa918fcd03bea856233921f88b9eef32896e2000000006a47304402201cd57a7064c100bb7e565a9aeff12bfe4397d59bd3d44a89115f97e2bd04669e022020cba46c8ab99a763c983f7fb10d61875495af0d6f42e3dfe010b843cb9c0ceb4121033288af9d515600042c64a8a058e80ad0a70f885ab4fc2424da847b18b74335e8feffffff46987a5d7920f32aa950c9cd258fa918fcd03bea856233921f88b9eef32896e2010000006a473044022022ce6618dca7e4d38455f327987f43f1ea127081e51375efe311e310b309aaed0220397f92dcebca00027adcfc11231b490125299ce71c38ff18c096d2272354b85f4121034a4a9529513993c0c4f44a011b0e53180e6ebace7791abfd0e291f6c4aeccef8feffffff910084749909b991b60ff63962ba5f01f2fd30e155f5f5776e694a14eb58e76a000000006a47304402201a4a9c14879acdbde902d6ec27c680f6bbf7c399296b0da31eaaad896dd0451b02201defdcc8514d8fea8425bc18406adf23f4957c218c0f321b9db3850f0b16884e412102a4b2aabf9cbfb9031de4f00d1997f10fe232e7e344b7ceb39e382be9b2e5002dfeffffff910084749909b991b60ff63962ba5f01f2fd30e155f5f5776e694a14eb58e76a010000006a47304402200fe83fbb8c1055190395bf46f8e1521670b1da12680950ea7b40ef5ad02ab7ac02205794d2fba2353cf6e8c9372b9e8900fa40fb5574880be5b455d6927b28fcbfc24121034a4a9529513993c0c4f44a011b0e53180e6ebace7791abfd0e291f6c4aeccef8feffffff0294daf505000000001976a914a12a69314c08a5155d779a2ec247ea735ade23bd88ac006d7c4d000000001976a9146dbb06e4c0395ffdec982856beab28994a548dce88ac69000000")
	if err != nil {
		fmt.Println(err)
		return
	}

	prevTx2, err := bt.NewTxFromString("02000000016fc96646b49acbe283ca81813da5ce0cf6b34a79dda74d515eaf68236ac7e2ba000000006a47304402205c1a6ba8018fa5d8c8952d37e4e21b731ac09edb491a2f475133021e348a1e5c02205acba3d90d31738a192593b66940ca119fd7a2e018c198b28d432db68e182034412103e72d6d9988b7fffcdef654e3c40c1227539b90a89dc5f42cd3d850e74ad94503feffffff025e266bee000000001976a9143355c640863b680e977d3608075ee5749f98106188ac0065cd1d000000001976a914a0416fb58b878bfaede66f83bb0e8c9fe0b0619c88ac66000000")
	if err != nil {
		fmt.Println(err)
		return
	}

	params = append(params, interpreter.ExecutionParams{
		Tx:            tx2,
		InputIdx:      0,
		PreviousTxOut: prevTx2.OutputIdx(int(tx.InputIdx(0).PreviousTxOutIndex)),
	})

	vm := interpreter.NewEngine()
	errs, _ := errgroup.WithContext(context.TODO())
	for _, p := range params {
		param := p
		errs.Go(func() error {
			input := param.Tx.InputIdx(param.InputIdx)
			inputASM, err := input.UnlockingScript.ToASM()
			if err != nil {
				return err
			}

			outputASM, err := param.PreviousTxOut.LockingScript.ToASM()
			if err != nil {
				return err
			}

			fmt.Println(inputASM)
			fmt.Println(outputASM)
			return vm.Execute(param)
		})
	}

	if err := errs.Wait(); err != nil {
		fmt.Println(err)
		return
	}

	// Unordered output:
	// 304402206579610b3a845e7ffa58203c686ca86ed3f2f946454bcb5f78e960c8ec34617702206cf0f168267acbca0acdc7fe38311fd94fd821868891aa1da150fe0de6e0ff6c41 03bb0164c11476e32287120301be5aca1310b0f72579f83e88cf6e10e42f6f78f1
	// OP_DUP OP_HASH160 a0416fb58b878bfaede66f83bb0e8c9fe0b0619c OP_EQUALVERIFY OP_CHECKSIG
	// 304502210086c83beb2b2663e4709a583d261d75be538aedcafa7766bd983e5c8db2f8b2fc02201a88b178624ab0ad1748b37c875f885930166237c88f5af78ee4e61d337f935f41 03e8be830d98bb3b007a0343ee5c36daa48796ae8bb57946b1e87378ad6e8a090d
	// OP_DUP OP_HASH160 805096c5167877a5799977d46fb9dee5891dc3cb OP_EQUALVERIFY OP_CHECKSIG
}

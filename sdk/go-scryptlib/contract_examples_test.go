package scryptlib


import (
    "testing"
    "context"
    "math/big"
    "encoding/hex"

    "github.com/stretchr/testify/assert"

    "github.com/libsv/go-bk/bip32"
	"github.com/libsv/go-bk/crypto"
    "github.com/libsv/go-bt/v2"
    "github.com/libsv/go-bt/v2/sighash"
    "github.com/libsv/go-bt/v2/bscript/interpreter/scriptflag"
)


func TestContractDemo(t *testing.T) {
    compilerResult, err := compilerWrapper.CompileContractFile("./test/res/demo.scrypt")
    assert.NoError(t, err)

    desc, err := compilerResult.ToDescWSourceMap()
    assert.NoError(t, err)

    contractDemo, err := NewContractFromDesc(desc)
    assert.NoError(t, err)

    x := Int{big.NewInt(7)}
    y := Int{big.NewInt(4)}
    constructorParams := map[string]ScryptType {
        "x": x,
        "y": y,
    }

    err = contractDemo.SetConstructorParams(constructorParams)
    assert.NoError(t, err)

    sumCorrect := Int{big.NewInt(11)}
    addParams := map[string]ScryptType {
        "z": sumCorrect,
    }
    err = contractDemo.SetPublicFunctionParams("add", addParams)
    assert.NoError(t, err)

    success, err := contractDemo.EvaluatePublicFunction("add")
    assert.NoError(t, err)
    assert.Equal(t, true, success)

    subCorrect := Int{big.NewInt(3)}
    subParams := map[string]ScryptType {
        "z": subCorrect,
    }
    err = contractDemo.SetPublicFunctionParams("sub", subParams)
    assert.NoError(t, err)

    success, err = contractDemo.EvaluatePublicFunction("sub")
    assert.NoError(t, err)
    assert.Equal(t, true, success)
}

func TestContractP2PKH(t *testing.T) {
    compilerResult, err := compilerWrapper.CompileContractFile("./test/res/p2pkh.scrypt")
    assert.NoError(t, err)

    desc, err := compilerResult.ToDescWSourceMap()
    assert.NoError(t, err)

    contractP2PKH, err := NewContractFromDesc(desc)
    assert.NoError(t, err)

    privKeyBase58 := "xprv9s21ZrQH143K3QTDL4LXw2F7HEK3wJUD2nW2nRk4stbPy6cq3jPPqjiChkVvvNKmPGJxWUtg6LnF5kejMRNNU3TGtRBeJgk33yuGBxrMPHi"
    priv, err := bip32.NewKeyFromString(privKeyBase58)
    assert.NoError(t, err)
	pub, err := priv.ECPubKey()
    assert.NoError(t, err)
    addr := crypto.Hash160(pub.SerialiseCompressed())

    pubKeyHash := Ripemd160{addr}
    constructorParams := map[string]ScryptType {
        "pubKeyHash": pubKeyHash,
    }
    err = contractP2PKH.SetConstructorParams(constructorParams)
    assert.NoError(t, err)

    tx := bt.NewTx()
    assert.NotNil(t, tx)

    lockingScript, err := contractP2PKH.GetLockingScript()
    assert.NoError(t, err)
    lockingScriptHex := hex.EncodeToString(*lockingScript)
    err = tx.From(
        "07912972e42095fe58daaf09161c5a5da57be47c2054dc2aaa52b30fefa1940b", // Random TXID
        0,
        lockingScriptHex,
        5000)
    assert.NoError(t, err)

    ecPrivKey, err := priv.ECPrivKey()
    assert.NoError(t, err)
    localSigner := bt.LocalSigner{ecPrivKey}

    var shf sighash.Flag = sighash.AllForkID
    _, sigBytes, err := localSigner.Sign(context.Background(), tx, 0, shf)
    assert.NoError(t, err)
    sig, err := NewSigFromDECBytes(sigBytes, shf)
    assert.NoError(t, err)

    unlockParams := map[string]ScryptType {
        "sig": sig,
        "pubKey": PubKey{pub},
    }
    err = contractP2PKH.SetPublicFunctionParams("unlock", unlockParams)
    assert.NoError(t, err)

    executionContext := ExecutionContext{
        Tx:             tx,
        InputIdx:       0,
        Flags:          scriptflag.EnableSighashForkID | scriptflag.UTXOAfterGenesis,
    }

    contractP2PKH.SetExecutionContext(executionContext)

    success, err := contractP2PKH.EvaluatePublicFunction("unlock")
    assert.NoError(t, err)
    assert.Equal(t, true, success)
}

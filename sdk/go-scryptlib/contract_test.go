package scryptlib


import (
    "os"
    "testing"
    "math/big"

    "github.com/stretchr/testify/assert"
)


var compilerWrapper CompilerWrapper
var contractDemo Contract

func TestMain(m *testing.M) {
    compilerBin, _ := FindCompiler()

    compilerWrapper = CompilerWrapper {
            CompilerBin: compilerBin,
            OutDir: "./out",
            HexOut: true,
            Debug: true,
            Desc: true,
            Stack: true,
            Optimize: false,
            CmdArgs: "",
            Cwd: "./",
        }

    compilerResult, _ := compilerWrapper.CompileContractFile("./test/res/demo.scrypt")
    desc, _ := compilerResult.ToDescWSourceMap()
    contractDemo, _ = NewContractFromDesc(desc)

    os.Exit(m.Run())
}

func TestContractEval(t *testing.T) {
    x := Int{big.NewInt(7)}
    y := Int{big.NewInt(4)}
    constructorParams := map[string]ScryptType {
        "x": x,
        "y": y,
    }

    err := contractDemo.SetConstructorParams(constructorParams)
    assert.NoError(t, err)

    // Correct sum:
    sumCorrect := Int{big.NewInt(11)}
    addParams := map[string]ScryptType {
        "z": sumCorrect,
    }
    err = contractDemo.SetPublicFunctionParams("add", addParams)
    assert.NoError(t, err)

    success, err := contractDemo.EvaluatePublicFunction("add")
    assert.NoError(t, err)
    assert.Equal(t, true, success)

    // Incorect sum:
    sumIncorrect := Int{big.NewInt(-23)}
    addParams = map[string]ScryptType {
        "z": sumIncorrect,
    }
    err = contractDemo.SetPublicFunctionParams("add", addParams)
    assert.NoError(t, err)
    success, err = contractDemo.EvaluatePublicFunction("add")
    assert.Error(t, err)
    assert.Equal(t, false, success)
}

func TestContractParamCheck(t *testing.T) {
    // TODO
}

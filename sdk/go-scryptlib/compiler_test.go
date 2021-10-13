package scryptlib


import (
    "testing"

    "github.com/stretchr/testify/assert"
)


func TestCompiler0(t *testing.T) {
    compilerBin, err := FindCompiler()
    assert.NoError(t, err)

    compilerWrapper := CompilerWrapper {
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

    compilerResult, err := compilerWrapper.CompileContractFile("./test/res/p2pkh.scrypt")
    assert.NoError(t, err)

    _, err = compilerResult.ToDescWSourceMap()
    assert.NoError(t, err)


}

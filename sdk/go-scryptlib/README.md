# go-scryptlib
An [sCrypt](https://scrypt.io/) SDK for the Go language.

You can learn all about writing sCrypt smart contracts in the official [docs](https://scryptdoc.readthedocs.io/en/latest/intro.html).

## Installation

To use the SDK, you need to get a copy of the [sCrypt compiler](https://scrypt.io/#download).

For installing the SDK, run the following command
```sh
go get github.com/sCrypt-Inc/go-scryptlib
```

## Usage

### Compiling an sCrypt contract

To compile an sCrypt contract, we must first initialize a CompilerWrapper:
```go
compilerBin, _ := FindCompiler()

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
```

Once that is initialized, we can compile the contract:
```go
compilerResult, _ := compilerWrapper.CompileContractFile("./test/res/demo.scrypt")
```

This leaves us with a struct of type CompilerResult. This step also outputs results of the compiler, and a contract description file in the "./out" directory, which we passed as a parameter to the CompilerWrapper.

From the compiler results we can derive an in-memory representation of the contract description tree:
```go
desc, _ := compilerResult.ToDescWSourceMap()
```

This is the basis, that will be used to create a Contract struct, which represents our compiled contract.
```go
contractDemo, _ := NewContractFromDesc(desc)
```

Then we can set the values for the contracts constructor. This is needed to create a valid locking script for out contract.
```go
x := Int{big.NewInt(7)}
y := Int{big.NewInt(4)}
constructorParams := map[string]ScryptType {
    "x": x,
    "y": y,
}

contractDemo.SetConstructorParams(constructorParams)

fmt.Println(contractDemo.GetLockingScript())
```

The same is true for our contracts public functions. Because our contract can contain many public functions, we use the functions name for referencing.
```go
sumCorrect := Int{big.NewInt(11)}
addParams := map[string]ScryptType {
    "z": sumCorrect,
}

contractDemo.SetPublicFunctionParams("add", addParams)

fmt.Println(contractDemo.GetUnlockingScript("add"))
```

We can then localy check, if a public function calls successfully evaluates.
```go
success, err := contractDemo.EvaluatePublicFunction("add")
```

The above method call will use the parameter values, that we set in the previous steps.


## Testing

Run `go test -v` in the root of this project.


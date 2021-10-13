package scryptlib

import (
    "fmt"
    "errors"
    "strconv"
    "reflect"
    "strings"
    "math/big"

    "github.com/libsv/go-bt/v2"
    "github.com/libsv/go-bt/v2/bscript"
    "github.com/libsv/go-bt/v2/bscript/interpreter"
    "github.com/libsv/go-bt/v2/bscript/interpreter/scriptflag"
)


type functionParam struct {
    Name        string
    TypeString  string
    Value       ScryptType
}

func (param *functionParam) setParamValue(value ScryptType) error {
    // TODO: TypeString should already be resolved so make sure the parameter values that get to here are too!
    if param.TypeString != value.GetTypeString() {
        errMsg := fmt.Sprintf("Passed object of type \"%s\" for parameter with name \"%s\". Expected \"%s\".",
                        value.GetTypeString(), param.Name, param.TypeString)
        return errors.New(errMsg)
    }

    param.Value = value
    return nil
}


type publicFunction struct {
    FunctionName string
    Index        int
    Params       []functionParam    // TODO: Maybe make this a map, because order is set by the hex template.
}

type ExecutionContext struct {
	Tx              *bt.Tx
	InputIdx        int
	Flags           scriptflag.Flag
}

type Contract struct {
    lockingScriptHexTemplate  string
    aliases                   map[string]string
    constructorParams         []functionParam
    publicFunctions           map[string]publicFunction
    structTypes               map[string]Struct             // Templates of contracts struct types. Maps struct names to related templates.
    executionContext          ExecutionContext
    contextSet                bool
}

// Set values for the contracts constructors parameters. 
// The value of "params" must be a map, that maps a public function name (string) to an ScryptType object.
func (contract *Contract) SetConstructorParams(params map[string]ScryptType) error {
    if len(params) != len(contract.constructorParams) {
            errMsg := fmt.Sprintf("Passed %d parameter values to constructor, but %d expected.",
                        len(params), len(contract.constructorParams))
            return errors.New(errMsg)
    }

    for idx := range contract.constructorParams {
        paramPlaceholder := &contract.constructorParams[idx]
        value := params[paramPlaceholder.Name]

        typePlaceholder := reflect.TypeOf(paramPlaceholder.Value).Name()
        typeActualParam := reflect.TypeOf(value).Name()

        if typePlaceholder != typeActualParam {
            errMsg := fmt.Sprintf("Passed value for param with name \"%s\" is not of the right type. Got \"%s\" but expected \"%s\"",
                            paramPlaceholder.Name, typeActualParam, typeActualParam)
            return errors.New(errMsg)
        }

        if typePlaceholder == "Struct" {
            same := IsStructsSameStructure(paramPlaceholder.Value.(Struct), value.(Struct))
            if ! same {
                errMsg := fmt.Sprintf("Passed Struct value for param with name \"%s\" is not of the right structure.",
                                    paramPlaceholder.Name)
                return errors.New(errMsg)
            }
        } else if typePlaceholder == "Array" {
            same := IsArraySameStructure(paramPlaceholder.Value.(Array), value.(Array))
            if ! same {
                errMsg := fmt.Sprintf("Passed Array value for param with name \"%s\" is not of the right structure.",
                                    paramPlaceholder.Name)
                return errors.New(errMsg)
            }
        }

        paramPlaceholder.setParamValue(value)
    }

    return nil
}

// Set values for a specific public function parameters. 
// The value of "params" must be a map, that maps a public function name (string) to an ScryptType object.
func (contract *Contract) SetPublicFunctionParams(functionName string, params map[string]ScryptType) error {
    function := contract.publicFunctions[functionName]

    if len(params) != len(function.Params) {
            errMsg := fmt.Sprintf("Passed %d parameter values to function \"%s\", but %d expected.",
                        len(params), function.FunctionName, len(function.Params))
            return errors.New(errMsg)
    }

    for idx := range function.Params {
        paramPlaceholder := &function.Params[idx]
        value := params[paramPlaceholder.Name]

        typePlaceholder := reflect.TypeOf(paramPlaceholder.Value).Name()
        typeActualParam := reflect.TypeOf(value).Name()

        if typePlaceholder != typeActualParam {
            errMsg := fmt.Sprintf("Passed value for param with name \"%s\" is not of the right type. Got \"%s\" but expected \"%s\"",
                            paramPlaceholder.Name, typeActualParam, typeActualParam)
            return errors.New(errMsg)
        }

        if typePlaceholder == "Struct" {
            same := IsStructsSameStructure(paramPlaceholder.Value.(Struct), value.(Struct))
            if ! same {
                errMsg := fmt.Sprintf("Passed Struct value for param with name \"%s\" is not of the right structure.",
                                    paramPlaceholder.Name)
                return errors.New(errMsg)
            }
        } else if typePlaceholder == "Array" {
            same := IsArraySameStructure(paramPlaceholder.Value.(Array), value.(Array))
            if ! same {
                errMsg := fmt.Sprintf("Passed Array value for param with name \"%s\" is not of the right structure.",
                                    paramPlaceholder.Name)
                return errors.New(errMsg)
            }
        }

        paramPlaceholder.setParamValue(value)
    }

    return nil
}

// Returns if the contracts execution context was already set at least once.
func (contract *Contract) IsExecutionContextSet() bool {
    return contract.contextSet
}

// Set the execution context, that will be used while evaluating a contracts public function.
// The locking and unlocking scripts, that you wan't to evaluate can be just templates, as they will be substitued localy, 
// while calling the EvaluatePublicFunction method.
func (contract *Contract) SetExecutionContext(ec ExecutionContext) {
    contract.executionContext = ec
    contract.contextSet = true
}

// Evaluate a public function call locally and return whether the evaluation was successfull,
// meaning the public function call (unlocking script) successfully evaluated against the contract (lockingScript).
// Constructor parameter values and also the public function parameter values MUST be set.
func (contract *Contract) EvaluatePublicFunction(functionName string) (bool, error) {
    // TODO: Check if parameter vals haven't been set yet. Use flags.

    lockingScript, err := contract.GetLockingScript()
    if err != nil {
        return false, err
    }
    unlockingScript, err := contract.GetUnlockingScript(functionName)
    if err != nil {
        return false, err
    }

    if ! contract.contextSet {
        err = interpreter.NewEngine().Execute(interpreter.WithScripts(lockingScript, unlockingScript))
        if err != nil {
            return false, err
        }
    } else {
        //input := contract.executionContext.Tx.InputIdx(contract.executionContext.InputIdx)
        //if input == nil {
        //    return false, errors.New(fmt.Sprintf("Context transaction has no input with index %d.", contract.executionContext.InputIdx))
        //}
        contract.executionContext.Tx.Inputs[contract.executionContext.InputIdx].UnlockingScript = unlockingScript
        prevoutSats := contract.executionContext.Tx.InputIdx(contract.executionContext.InputIdx).PreviousTxSatoshis

        engine := interpreter.NewEngine()
        err = engine.Execute(
            //interpreter.WithScripts(
            //    lockingScript,
            //    unlockingScript,
            //),
            interpreter.WithTx(
                contract.executionContext.Tx,
                contract.executionContext.InputIdx,
                //contract.executionContext.PreviousTxOut,
                &bt.Output{LockingScript: lockingScript, Satoshis: prevoutSats},
            ),
            interpreter.WithFlags(
                contract.executionContext.Flags,
            ),
        )
        if err != nil {
            return false, err
        }
    }

    return true, nil
}

func (contract *Contract) GetUnlockingScript(functionName string) (*bscript.Script, error) {
    var res *bscript.Script
    var sb strings.Builder

    publicFunction := contract.publicFunctions[functionName]

    for _, param := range publicFunction.Params {
        paramHex, err := param.Value.Hex()
        if err != nil {
            return res, err
        }
        sb.WriteString(paramHex)
    }

    // Append public function index.
    index := Int{big.NewInt(int64(publicFunction.Index))}
    indexHex, err := index.Hex()
    if err != nil {
        return res, err
    }
    sb.WriteString(indexHex)

    unlockingScript, err := bscript.NewFromHexString(sb.String())
    if err != nil {
        return res, err
    }

    return unlockingScript, nil
}

func (contract *Contract) GetLockingScript() (*bscript.Script, error) {
    var res *bscript.Script
    lockingScriptHex := contract.lockingScriptHexTemplate

    // TODO: move to code part
    for _, param := range contract.constructorParams {
        paramHex, err := param.Value.Hex()
        if err != nil {
            return res, err
        }

        toReplace := fmt.Sprintf("<%s>", param.Name)
        lockingScriptHex = strings.Replace(lockingScriptHex, toReplace, paramHex, 1)
    }

    // TODO: Data part.

    lockingScript, err := bscript.NewFromHexString(lockingScriptHex)
    if err != nil {
        return res, err
    }

    return lockingScript, nil
}

// Returns a map, that maps struct names as defined in the contract to their respective instances of ScryptType.
func (contract *Contract) GetStructTypes() map[string]Struct {
    return contract.structTypes
}

func constructAbiPlaceholders(desc map[string]interface{}, structTypes map[string]Struct,
                                    aliases map[string]string) ([]functionParam, map[string]publicFunction, error) {
    var constructorParams []functionParam
    publicFunctions := make(map[string]publicFunction)

    // TODO: Pass this as a pram instead of recreating it here.
    structItemsByTypeString := getStructItemsByTypeString(desc)

    for _, abiItem := range desc["abi"].([]map[string]interface{}) {

        abiItemType := abiItem["type"].(string)
        params := abiItem["params"].([]map[string]string)

        var publicFunctionPlaceholder publicFunction
        var publicFunctionName string
        if abiItemType == "function" {
            publicFunctionName = abiItem["name"].(string)
            publicFunctionPlaceholder = publicFunction{
                                            FunctionName: publicFunctionName,
                                            Index: abiItem["index"].(int),
                                        }
        }

        for _, param := range params {
            var value ScryptType
            pName := param["name"]
            pType := param["type"]

            if IsStructType(pType) {
                // Create copy of struct template.
                structName := GetStructNameByType(pType)
                value = structTypes[structName]
            } else if IsArrayType(pType) {
                arrVal, err := constructArrayType(pType, structItemsByTypeString, aliases)
                if err != nil {
                    return nil, publicFunctions, err
                }
                value = arrVal
            } else {
                // Concrete values.
                val, err := createPrimitiveTypeWDefaultVal(pType)
                if err != nil {
                    return nil, publicFunctions, err
                }
                value = val
            }


            placeholder := functionParam{
                Name:       pName,
                TypeString: pType,
                Value:      value,
            }

            if abiItemType == "constructor" {
                constructorParams = append(constructorParams, placeholder)
            } else {
                publicFunctionPlaceholder.Params = append(publicFunctionPlaceholder.Params, placeholder)
            }
        }

        if abiItemType == "function" {
            publicFunctions[publicFunctionName] = publicFunctionPlaceholder
        }
    }

    return constructorParams, publicFunctions, nil
}

func getStructItemsByTypeString(desc map[string]interface{}) map[string]interface{} {
    structItemsByTypeString := make(map[string]interface{})
    for _, structItem := range desc["structs"].([]map[string]interface{}) {
        structType := structItem["name"].(string)
        structItemsByTypeString[structType] = structItem
    }
    return structItemsByTypeString
}


func constructStructTypes(structItemsByTypeString map[string]interface{}, aliases map[string]string) (map[string]Struct, error) {
    res := make(map[string]Struct)

    for structName, structItem := range structItemsByTypeString {
        structItem := structItem.(map[string]interface{})
        structType, err := constructStructType(structItem, structItemsByTypeString, aliases)
        if err != nil {
            return res, err
        }
        res[structName] = structType
    }

    return res, nil
}

func constructStructType(structItem map[string]interface{}, structItemsByTypeString map[string]interface{},
                                                    aliases map[string]string) (Struct, error) {

    var res Struct

    var keysInOrder []string
    values := make(map[string]ScryptType)

    params := structItem["params"].([]map[string]string)
    for _, param := range params {
        pName := param["name"]
        pType := param["type"]
        pTypeResolved := ResolveType(pType, aliases)

        keysInOrder = append(keysInOrder, pName)

        var val ScryptType
        var err error

        structItem, isStructType := structItemsByTypeString[pTypeResolved]
        if isStructType {
            val, err = constructStructType(structItem.(map[string]interface{}), structItemsByTypeString, aliases)
        } else if IsArrayType(pTypeResolved) {
            val, err = constructArrayType(pTypeResolved, structItemsByTypeString, aliases)
        } else {
            val, err = createPrimitiveTypeWDefaultVal(pTypeResolved)
        }

        if err != nil {
            return res, err
        }

        values[pName] = val
    }

    return Struct{
        keysInOrder: keysInOrder,
        values: values,
    }, nil
}

func constructArrayType(typeString string, structItemsByTypeString map[string]interface{},
                                       aliases map[string]string) (Array, error) {
    var res Array
    typeName, arraySizes := FactorizeArrayTypeString(typeString)

    var items []ScryptType
    for dimension := len(arraySizes) - 1; dimension >= 0; dimension-- {
        arraySize := arraySizes[dimension]
        nItems, _ := strconv.Atoi(arraySize)

        if dimension == len(arraySizes) - 1 {
            // Last dimension. Create concrete types here.
            for i := 0; i < nItems; i++ {
                var item ScryptType
                var err error

                structItem, isStructType := structItemsByTypeString[typeName]
                if isStructType {
                    item, err = constructStructType(structItem.(map[string]interface{}), structItemsByTypeString, aliases)
                } else {
                    item, err = createPrimitiveTypeWDefaultVal(typeName)
                }

                if err != nil {
                    return res, err
                }

                items = append(items, item)
            }
        } else {
            var itemsNewDimension []ScryptType
            for i := 0; i < nItems; i++ {
                // Copy items from level below.
                values := make([]ScryptType, len(items))
                copy(values, items)
                itemsNewDimension = append(itemsNewDimension, Array{values})
            }
            items = itemsNewDimension
        }
    }

    return Array{items}, nil
}

func createPrimitiveTypeWDefaultVal(typeString string) (ScryptType, error) {
    var res ScryptType
    switch typeString {
    case "bool":
        res = Bool{true}
    case "int":
        res = Int{big.NewInt(0)}
    case "bytes":
        res = Bytes{make([]byte, 0)}
    case "PrivKey":
        res = PrivKey{nil}
    case "PubKey":
        res = PubKey{nil}
    case "Sig":
        res = Sig{nil, 0}
    case "Ripemd160":
        res = Ripemd160{make([]byte, 0)}
    case "Sha1":
        res = Sha1{make([]byte, 0)}
    case "Sha256":
        res = Sha256{make([]byte, 0)}
    case "SigHashType":
        res = SigHashType{make([]byte, 0)}
    case "SigHashPreimage":
        res = SigHashPreimage{make([]byte, 0)}
    case "OpCodeType":
        res = OpCodeType{make([]byte, 0)}
    default:
        return res, errors.New(fmt.Sprintf("Unknown type string \"%s\".", typeString))
    }
    return res, nil
}

// Creates a new instance of Contract type from the contracts description tree.
func NewContractFromDesc(desc map[string]interface{}) (Contract, error) {
    var res Contract

    lockingScriptHexTemplate := desc["hex"].(string)
    aliases := ConstructAliasMap(desc["alias"].([]map[string]string))

    structItemsByTypeString := getStructItemsByTypeString(desc)

    // Construct instances of struct types.
    structTypes, err := constructStructTypes(structItemsByTypeString, aliases)
    if err != nil {
        return res, err
    }

    // structTypes should also contain keys for aliases.
    // TODO: Fix aliases for concrete vals.
    for key, val := range aliases {
        structTypes[key] = structTypes[val]
    }

    // Initialize constructor parameter placeholders and public functions along its parameter placeholders.
    constructorParams, publicFunctions, err := constructAbiPlaceholders(desc, structTypes, aliases)
    if err != nil {
        return res, err
    }

    return Contract{
        lockingScriptHexTemplate: lockingScriptHexTemplate,
        aliases: aliases,
        constructorParams: constructorParams,
        publicFunctions: publicFunctions,
        structTypes: structTypes,
        contextSet: false,
    }, nil
}

// Creates a new instance of Contract type.
func NewContract(compilerResult CompilerResult) (Contract, error) {
    var res Contract

    desc, err := compilerResult.ToDescWSourceMap()
    if err != nil {
        return res, err
    }

    res, err = NewContractFromDesc(desc)
    if err != nil {
        return res, err
    }

    return res, nil
}


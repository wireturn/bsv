package scryptlib

import (
    "fmt"
    "strings"
    "regexp"
    "strconv"
    "reflect"
    "math/big"
)


// Factor array declaration string to array type and sizes.
// e.g. 'int[N][N][4]' -> ('int', ['N', 'N', '4'])
func FactorizeArrayTypeString(typeStr string) (string, []string) {
    var arraySizes []string

    r := regexp.MustCompile(`\[([\w.]+)\]+`)
    matches := r.FindAllStringSubmatch(typeStr, -1)
    for _, match := range matches {
        arraySizes = append(arraySizes, match[1])
    }

    typeName := strings.Split(typeStr, "[")[0]

    return typeName, arraySizes
}

// Retruns array declaration string for given type name and sizes.
// Array sizes are passed as a slice of type []string.
func ToLiteralArrayTypeStr(typeName string, arraySizes []string)  string {
    var resBuff strings.Builder
    resBuff.WriteString(typeName)
    for _, size := range arraySizes {
        resBuff.WriteString(size)
    }
    return resBuff.String()
}

// Retruns array declaration string for given type name and sizes.
// Array sizes are passed as a slice of type []int.
// TODO: Change int types to *big.Int
func ToLiteralArrayTypeInt(typeName string, arraySizes []int)  string {
    var resBuff strings.Builder
    resBuff.WriteString(typeName)
    for _, size := range arraySizes {
        resBuff.WriteString(strconv.Itoa(size))
    }
    return resBuff.String()
}

// Check if string is of an array type.
// e.g. "int[2]" or "int[N][3]" 
func IsArrayType(typeStr string) bool {
    match, _ := regexp.MatchString(`^\w[\w.\s{}]+(\[[\w.]+\])+$`, typeStr)
    return match
}

// Check if string is of a struct type.
// e.g. "struct Point {}"
func IsStructType(typeStr string) bool {
    match, _ := regexp.MatchString(`^struct\s(\w+)\s\{\}$`, typeStr)
    return match
}

// Returns struct name from type string.
// e.g.: 'struct ST1 {}[2][2][2]' -> 'ST1'.
func GetStructNameByType(typeName string) string {
    r := regexp.MustCompile(`^struct\s(\w+)\s\{\}.*$`)
    match := r.FindStringSubmatch(typeName)
    if match != nil {
        return match[1]
    }
    return ""
}

func ResolveType(typeStr string, aliases map[string]string) string {
    if IsArrayType(typeStr) {
        typeName, arraySizes := FactorizeArrayTypeString(typeStr)
        return ToLiteralArrayTypeStr(ResolveType(typeName, aliases), arraySizes)
    }

    if IsStructType(typeStr) {
        return ResolveType(GetStructNameByType(typeStr), aliases)
    }

    resolvedType, ok := aliases[typeStr]
    if ok {
        return ResolveType(resolvedType, aliases)
    }

    if BASIC_SCRYPT_TYPES[typeStr] {
        return typeStr
    }
    return fmt.Sprintf("struct %s {}", typeStr)
}

func EvenHexStr(hexStr string) string {
    if len(hexStr) % 2 == 1 {
        return "0" + hexStr
    }
    return hexStr
}

func BigIntToHex_LE(value *big.Int) string {
    b := BigIntToBytes_LE(value)
    return EvenHexStr(fmt.Sprintf("%x", b))
}

func BigIntToBytes_LE(value *big.Int) []byte {
    b := value.Bytes()
    for i := 0; i < len(b)/2; i++ {
        b[i], b[len(b)-i-1] = b[len(b)-i-1], b[i]
    }
    return b
}

// Returns true if the passed Struct sCrypt types are of the same structure.
// Concrete values are not checked! It only recursively goes through Array and Struct types.
func IsStructsSameStructure(struct0 Struct, struct1 Struct) bool {
    if len(struct0.keysInOrder) != len(struct1.keysInOrder) {
        return false
    }
    if len(struct0.values) != len(struct1.values) {
        return false
    }

    for i, key := range struct0.keysInOrder {
        // Check key order.
        if struct1.keysInOrder[i] != key {
            return false
        }

        // Check values.
        type0 := reflect.TypeOf(struct0.values[key]).Name()
        type1 := reflect.TypeOf(struct1.values[key]).Name()
        if type0 != type1 {
            return false
        }

        // Go deeper if struct or array type.
        if type0 == "Struct" {
            return IsStructsSameStructure(struct0.values[key].(Struct), struct1.values[key].(Struct))
        }
        if type0 == "Array" {
            return IsArraySameStructure(struct0.values[key].(Array), struct1.values[key].(Array))
        }

    }

    return true
}

// Returns true if the passed Array sCrypt types are of the same structure.
// Concrete values are not checked! It only recursively goes through Array and Struct types.
func IsArraySameStructure(array0 Array, array1 Array) bool {
    if len(array0.values) != len(array1.values) {
        return false
    }

    for i, elem0 := range array0.values {
        elem1 := array1.values[i]

        // Check values.
        type0 := reflect.TypeOf(elem0).Name()
        type1 := reflect.TypeOf(elem1).Name()
        if type0 != type1 {
            return false
        }

        // Go deeper if struct or array type.
        if type0 == "Struct" {
            return IsStructsSameStructure(elem0.(Struct), elem1.(Struct))
        }
        if type0 == "Array" {
            return IsArraySameStructure(elem0.(Array), elem1.(Array))
        }
    }

    return true
}

// Construct a map for resolving alias types from the alias section of the contract description file.
func ConstructAliasMap(aliasesDesc []map[string]string) map[string]string {
    aliases := make(map[string]string)
    for _, item := range aliasesDesc {
        nameString := item["name"]
        typeString := item["type"]
        aliases[nameString] = typeString
    }
    return aliases
}

func reSubMatchMap(r *regexp.Regexp, str string) map[string]string {
    match := r.FindStringSubmatch(str)
    subMatchMap := make(map[string]string)
    for i, name := range r.SubexpNames() {
        if i != 0 {
            subMatchMap[name] = match[i]
        }
    }

    return subMatchMap
}

func reSubMatchMapAll(r *regexp.Regexp, str string) []map[string]string {
    var res []map[string]string

    matches := r.FindAllStringSubmatch(str, -1)
    for _, match := range matches {
        subMatchMap := make(map[string]string)
        for i, name := range r.SubexpNames() {
            if i != 0 {
                subMatchMap[name] = match[i]
            }
        }
        res = append(res, subMatchMap)
    }

    return res
}


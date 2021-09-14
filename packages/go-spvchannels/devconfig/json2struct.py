#!/usr/bin/env python3

import json
import argparse
import sys

def key2struct(jsonkey:str):
    """
    string key of type my_key --> MyKey
    """
    parts = jsonkey.split('_')
    ret = ""
    for w in parts:
        ret += w.capitalize()
    return ret

def json2struct(jobj:object, name:str, indent:int=0):
    """
    Parse a json string from a rest reply and convert it to a golang struct
    """

    ret = ""
    objbody = jobj
    if indent==0:
        ret+= 'type {} struct {{\n'.format(key2struct(name))
    else:
        indent_head = indent*4*' '
        if isinstance(jobj, list) and len(jobj)>0:
            ret += '{}{} []struct {{\n'.format(indent_head, key2struct(name))
            objbody = jobj[0]
        else:
            ret += '{}{} struct {{\n'.format(indent_head, key2struct(name))


    indent_body = 4*(indent+1)*' '
    for k, v in objbody.items():
        goname = key2struct(k)
        vtype = type(v)
        if vtype is bool:
            ret += '{}{} bool `json:"{}"`\n'.format(indent_body, goname, k)
        elif vtype is int:
            ret += '{}{} int `json:"{}"`\n'.format(indent_body, goname, k)
        elif vtype is str:
            ret += '{}{} string `json:"{}"`\n'.format(indent_body, goname, k)
        elif vtype is list or vtype is dict:
            ret += json2struct(v, k, indent+1)
        else:
            raise RuntimeError('Unhandled jsont key {}, value type {}'.format(k, vtype))

    ret += indent*4*' ' + '}'
    if indent < 1:
        ret +='\n'
    else:
        ret += ' `json:"{}"`\n'.format(name)
    return ret
#################################################################################################################
#################################################################################################################


def test_get_bitbucket_status():
    jstr = '{"id":"2vkapEui-Cfb3tY7l9FFviRjpsNGa0Iv4kFEHYoMWJdl4f9PSlvurjOCnTBzH1r_C8VUuvQsn-0NsO0Q2bKGUA","href":"https://localhost:5010/api/v1/channel/2vkapEui-Cfb3tY7l9FFviRjpsNGa0Iv4kFEHYoMWJdl4f9PSlvurjOCnTBzH1r_C8VUuvQsn-0NsO0Q2bKGUA","public_read":true,"public_write":true,"sequenced":true,"locked":false,"head":0,"retention":{"min_age_days":0,"max_age_days":99999,"auto_prune":true},"access_tokens":[{"id":"2","token":"cadpwRhA9H6N-_3vZf9uymj08iYVNPmdDMVUNcjfhD5UM5YpYpXXFB84WzOz6lwESuCYgMh8vD7XbrE4Z4a0PA","description":"Owner","can_read":true,"can_write":true}]}'
    jobj = json.loads(jstr )
    go_truct = json2struct(jobj, "MyCustomStruct", 0)
    print('\n\n{}'.format(go_truct))

def run_main():
    print("Run main here with input json string")

if __name__ == '__main__':
    run_main()

import json
import os
import sys


def convert_type_registry(name):
    module_path = os.path.dirname(__file__)
    path = os.path.join(module_path, '{}.json'.format(name))

    with open(os.path.abspath(path), 'r') as fp:
        data = fp.read()
    type_registry = json.loads(data)

    convert_dict = dict()

    for type_string, type_struct in type_registry.items():

        if type(type_struct) == dict:
            if '_enum' in type_struct:
                convert_dict[type_string] = {
                    "type": "enum",
                    "type_mapping": try_struct_convert(type_struct["_enum"])
                }
            elif '_struct' in type_struct:
                convert_dict[type_string] = type_struct["_struct"]
            else:
                convert_dict[type_string] = {
                    "type": "struct",
                    "type_mapping": try_struct_convert(type_struct)
                }
        else:
            convert_dict[type_string] = type_struct
    return convert_dict


def try_struct_convert(struct):
    n = []
    if type(struct) == list:
        return struct
    for type_string, type_struct in struct.items():
        if type_struct is None:
            type_struct = "null"
        n.append([type_string, type_struct])
    return n


if __name__ == '__main__':
    print(json.dumps(convert_type_registry(sys.argv[1])))

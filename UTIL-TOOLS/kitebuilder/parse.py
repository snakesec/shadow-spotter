import json
from json import JSONDecodeError
from os import listdir
from os.path import isfile, join
import logging
import hashlib
import ksuid
import yaml

from parserequestbody import requestBodyParse
from resolverefs import resolve_ref, resolve_schema_object

NRM="\001\033[0m\002"       

BLD="\001\033[1m\002"

BLK="\001\033[1;90m\002"       
RED="\001\033[1;91m\002"       
GRN="\001\033[1;92m\002"       
YEL="\001\033[1;93m\002"       
BLU="\001\033[1;94m\002"       
PUR="\001\033[1;95m\002"       
CYN="\001\033[1;96m\002"       
WHT="\001\033[1;97m\002"



seen_hashes = {}


class ParserSpecPathDuplicateException(Exception):
    pass


def parse_security_definitions(spec):
    """
    Parses and pulls the security definitions from the spec, excluding oauth2 definitions.
    :param spec: the swagger spec being parsed
    :return: dictionary of the specs security definitions
    """
    defs = spec.get("securityDefinitions", {})

    if type(defs) != dict:
        return

    return {
        k: v for (k, v) in defs.items() if
        type(v) == dict and v.get("type") != "oauth2"
    }


def parse_security_schemes(spec):
    """
    Parses and pulls the security schemes from the v3 spec, excluding oauth2 definitions.
    :param spec: the openapi spec being parsed
    :return: dictionary of the specs security schemes
    """
    defs = spec.get("components", {}).get("securitySchemes", {})

    if type(defs) != dict:
        return

    return {
        k: v for (k, v) in defs.items() if
        type(v) == dict and v.get("type") != "oauth2"
    }


def resolve_parameter(spec, content_types, parameter):
    if type(parameter) != dict:
        print("Return parameter is not a dict")
        return

    param_ref = parameter.get("$ref", None)
    if param_ref is not None:
        parameter = resolve_ref(spec, param_ref)
        if parameter is None:
            return

    # skip including any non dictionaries
    schema = parameter.get("schema", {})
    if schema and type(schema) == dict:
        resolved_schema = resolve_schema_object(spec, schema)

        if resolved_schema:
            parameter.update({"schema": resolved_schema})

    elif schema and type(schema) != dict:
        parameter.pop("schema")

    # handle Parameter array items
    items = parameter.get("items", {})
    if items and type(items) == dict:
        resolved_item = resolve_schema_object(spec, items)

        if resolved_item:
            parameter.update({"items": resolved_item})

    # basic check on name and description to determine if the string is a UUID
    param_desc = parameter.get("description", "")
    if param_desc and type(param_desc) == str and "uuid" in param_desc.lower():
        parameter["type"] = "uuid"

    param_name = parameter.get("name", "")
    if param_name and type(param_desc) == str and "uuid" in param_name.lower():
        parameter["type"] = "uuid"

    param_example = parameter.get("example", "")
    if param_example and type(param_example) == str and "uuid" in param_example.lower():
        parameter["type"] = "uuid"

    # check if a regex pattern exists, supplying it as the format if it does
    pattern = parameter.get("pattern", None)
    if pattern:
        parameter["format"] = pattern

    # x-examples
    x_examples = parameter.get("x-examples", None)
    _unknown_ct_idx = 0
    if x_examples is not None:
        if type(x_examples) == dict:
            new_examples = {}
            for content_type, example_data in x_examples.items():
                if content_type == "description":
                    continue

                if content_type not in content_types:
                    new_examples[f"unknown/{_unknown_ct_idx}"] = str(example_data)
                    _unknown_ct_idx += 1
                else:
                    new_examples[content_type] = example_data

            x_examples = new_examples

        if type(x_examples) == list:
            for idx, example in enumerate(x_examples):
                x_examples[f"unknown/{idx}"] = str(example)

        parameter["x-examples"] = x_examples

    return parameter

# check for extra parameters in path, such as :token and replace with {token} and a parameter in each route
def resolve_extra_path_params(path):
    extra_path_params = []
    for _idx, path_chunk in enumerate(path.split('/')):
        if len(path_chunk) > 0 and path_chunk[0] == ":":
            path_chunk = f"{{{path_chunk[1:]}}}"

            extra_path_params.append({
                "in": "path",
                "name": path_chunk,
                "required": True,
                "type": "number" if path_chunk.lower()[-2:] == "id" else "string"
            })

        if len(path_chunk) > 0 and path_chunk[0] == "[" and path_chunk[-1] == "]":
            path_chunk = f"{{{path_chunk[1:-1]}}}"

            extra_path_params.append({
                "in": "path",
                "name": path_chunk,
                "required": True,
                "type": "number" if path_chunk.lower()[-2:] == "id" else "string"
            })

    return extra_path_params


def parse_swagger_spec(file_name, blacklist):
    with open(file_name) as f:
        content = f.read()

        try:
            spec = None
            if "json" in file_name.lower():
                spec = json.loads(content)

            if "yaml" in file_name.lower():
                spec = yaml.safe_load(content)
            
        except JSONDecodeError as e:
            raise Exception(f"Failed to parse invalid spec (invalid JSON): {file_name}")

        except yaml.YAMLError:
                raise Exception(f"Failed to parse invalid spec (invalid YAML): {file_name}")

    if type(spec) != dict:
        raise Exception(f"Failed to parse invalid spec (invalid spec, invalid format): {file_name}")

    host = None
    security_definitions = None

    if spec.get("openapi") is not None:
        host = spec.get("servers", [])[0]["url"]
        security_definitions = parse_security_definitions(spec)

    if spec.get("swagger") is not None:
        host = spec.get("host")
        if type(host) == list and len(host) > 0:
            host = host[0]
        security_definitions = parse_security_schemes(spec)

    if host is not None and any(_host in host for _host in blacklist):
        raise Exception(f"Failed to parse invalid spec (blacklisted host): {file_name}")

    paths = spec.get("paths", {})
    if type(paths) != dict:
        raise Exception(f"Failed to parse invalid spec (invalid spec, missing paths): {file_name}")

    if paths is None:
        raise Exception(f"Failed to parse invalid spec (invalid spec, missing paths): {file_name}")

    base_path = spec.get('basePath', '/')
    if base_path is None:
        raise Exception(f"Failed to parse invalid spec (invalid spec, missing basePath)")

    # don't allow go templated base path, default to /
    if "{{." in base_path:
        base_path = '/'

    custom_schema = {
        "ksuid": str(ksuid.ksuid()),
        "url": host,
        "securityDefinitions": security_definitions,
        "paths": {}
    }

    for path in paths:
        if type(paths[path]) == str:
            print("Error path, just continue")
            continue

        full_path = f"{base_path}{path}".replace('//', '/')

        extra_path_params = resolve_extra_path_params(full_path)

        if full_path not in custom_schema["paths"]:
            custom_schema["paths"][full_path] = {}

        method_size = 0

        try:
            for method, endpoint_data in paths[path].items():

                method_size + 1
                
                if type(endpoint_data) != dict:
                    print("endpoint_data is not a dict")
                    continue

                path_params = endpoint_data.get("parameters", [])
                if type(path_params) != list:
                    print("path_params is not list")
                    continue

                body_params = None
                try:      
                    body_params = endpoint_data.get("requestBody", [])
                    
                except:
                    print("Body_params is wrong...")
                    None

                if isinstance(body_params, list) and all(isinstance(item, tuple) for item in body_params):
                    body_params = dict(body_params)

                if '$ref' in body_params:
                    ref = body_params['$ref']
                    body_params = resolve_schema_object(spec, resolve_ref(spec, ref))

                if type(body_params) != list:
                    body_params = list(body_params.items())

                description = endpoint_data.get("description", None)
                if description is None:
                    description = endpoint_data.get("summary", "")

                operation_id = spec["paths"][path][method].get("operationId", None)
                consumes = endpoint_data.get("consumes", None)

                if consumes is None:
                    consumes = []

                parameters = [
                    parameter for parameter in
                    list(map(lambda param_fields: resolve_parameter(spec, consumes, param_fields), path_params))
                    if parameter is not None
                ]

                parameters.extend(extra_path_params)

                print("Using Method: {}{}{}".format(CYN, method, NRM))

                seen_route_hashes = seen_hashes.get(path, {}).get(method, [])

                # skip this spec if circular references are detected
                try:
                    json_params = json.dumps(parameters, sort_keys=True).encode('utf-8')
                    param_hash = hashlib.md5(json_params).hexdigest()
                except ValueError:
                    logging.error(f"Failed to parse invalid spec (circular references): {file_name}")
                    print("Getting out due to circular reference")
                    return

                if param_hash in seen_route_hashes:
                    raise ParserSpecPathDuplicateException

                seen_route_hashes.append(param_hash)
                if seen_hashes.get(path, None) is None:
                    seen_hashes[path] = {}

                seen_hashes[path][method] = seen_route_hashes

                tmp_para = requestBodyParse(spec, body_params, 0)

                real_parameter = []

                if method == "get":
                    real_parameter = parameters
                else:
                    real_parameter.extend(parameters)

                    if tmp_para == None:
                        tmp_para = {}
                    else:
                        real_parameter.extend([tmp_para])

                    consumes = requestBodyParse(spec, body_params, 1)

                if real_parameter == None:
                    real_parameter = []

                print("Parameters: {}{}{}".format(WHT, json.dumps(real_parameter, indent=4), NRM))

                custom_schema["paths"][full_path][method] = {
                    
                    "description": description,
                    "operationId": operation_id,
                    "parameters": real_parameter,
                    "consumes": consumes,
                    "produces": endpoint_data.get("produces", None),
                    
                }

                print("Custom Schema: {}{}{}".format(PUR, custom_schema["paths"][full_path][method], NRM))

        except ParserSpecPathDuplicateException:
            print("ParserSpecPathDuplicateException")
            continue

        except RecursionError:
            print("RecursionError")
            continue

    if len(custom_schema.get("paths", [])) == 0:
        print("Len of paths is 0")
        return

    return custom_schema, method_size


def parse_specs(scrape_dir, output_file, blacklist):
    try:
        files = [f"{scrape_dir}/{f}" for f in listdir(scrape_dir) if isfile(join(scrape_dir, f))]
    except FileNotFoundError:
        logging.error("Invalid spec directory supplied.")
        return

    output_contents = []

    for idx, file in enumerate(files):
        try:
            custom_schema, method_size = parse_swagger_spec(file, blacklist)
        except KeyboardInterrupt:
            exit()

        except Exception as e:
            logging.error(f"[{idx}, {file}] {e} AAAAA")
            continue

        if custom_schema is None:
            print("Custom_schema is None")
            continue

        
        output_contents.append(custom_schema)

        logging.info(f"[{idx}] Successfully parsed spec: {file}")

    with open(output_file, 'w+') as f:
        print("\n{}{}{}\n".format(WHT, output_contents, NRM))
        f.write(json.dumps(output_contents, indent=4))
        logging.info(f"Wrote output to file: {output_file}")

    return len(output_contents)

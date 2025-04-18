"""
Shadow-Spotter Next Gen Content Discovery
Copyright (C) 2024  Weidsom Nascimento - SNAKE Security

Based on kiterunner from AssetNote

This program is free software: you can redistribute it and/or modify
it under the terms of the GNU Affero General Public License as
published by the Free Software Foundation, either version 3 of the
License, or (at your option) any later version.

This program is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Affero General Public License for more details.

You should have received a copy of the GNU Affero General Public License
along with this program.  If not, see <https://www.gnu.org/licenses/>.
"""

from resolverefs import resolve_ref, resolve_schema_object
import json

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

def find_properties(data):
    all_properties = {}

    if isinstance(data, dict):
        if 'properties' in data:
            all_properties.update(data['properties'])
        
        for value in data.values():
            all_properties.update(find_properties(value))
    
    elif isinstance(data, list):
        for item in data:
            all_properties.update(find_properties(item))

    return all_properties

def requestBodyParse(spec, body_params, required_consumes):
    content_body = []
    reqbodyserialized = []
    
    try:
        for idx, elmt in enumerate(body_params):
            if 'content' in elmt:
                content_idx = idx
                break

        for ab in body_params[content_idx]:
            if isinstance(ab, dict):
                first_key = next(iter(ab))
                consumes_type = first_key

                if ab[first_key]['schema'].get('properties') is not None:
                    # Lista para armazenar todos os objetos processados
                    all_reqbodyserialized = []

                    for prop_key, prop_value in ab[first_key]['schema'].get('properties', {}).items():
                        
                        prop_ref = prop_value.get("$ref", None)
                        
                        if prop_ref:
                            resolved_prop_ref = resolve_ref(spec, prop_ref)
                            reqbodyserialized = {}
                            
                            for rpr_keys, rpr_value in resolved_prop_ref['properties'].items():    
                                reqbodyserialized[rpr_keys] = {    
                                    "name": rpr_keys,
                                    "in": "body",
                                    "description": "null",
                                    "required": True,
                                    "schema": {
                                        "type": rpr_value['type'] or "",
                                    }
                                }

                                if "format" in rpr_value:
                                    reqbodyserialized[rpr_keys]['schema']['format'] = rpr_value['format']

                            # Adiciona os objetos processados à lista
                            all_reqbodyserialized.extend(list(reqbodyserialized.values()))

                        else:
                            reqbodyserialized = {}
                            prop_resolved = ab[first_key]['schema'].get('properties', {})

                            for rpr_keys, rpr_value in prop_resolved.items(): 
                                reqbodyserialized[rpr_keys] = {    
                                    "name": rpr_keys,
                                    "in": "body",
                                    "description": "null",
                                    "required": True,
                                    "schema": {
                                        "type": rpr_value['type'] or "",
                                    }
                                }

                                if "format" in rpr_value:
                                    reqbodyserialized[rpr_keys]['schema']['format'] = rpr_value['format']

                    # Adiciona os objetos processados à lista
                    all_reqbodyserialized.extend(list(reqbodyserialized.values()))

                else:
                    all_reqbodyserialized = []

                    schema_ref = ab[first_key]['schema'].get("$ref", None)
                    schema_ref_tmp = next(iter(ab[first_key]['schema'].values()), None)
                    schema_ref_bkp = next((item['$ref'] for item in schema_ref_tmp if '$ref' in item), None)
                    
                    if schema_ref:
                        resolved_schema = resolve_ref(spec, schema_ref)
                        resolved_schema = resolve_schema_object(spec, resolved_schema)
                        resolved_schema = find_properties(resolved_schema)

                        if resolved_schema:
                            schema = resolved_schema
                        else:
                            schema = dict(type="object")

                    elif schema_ref_bkp:
                          resolved_schema = resolve_ref(spec, schema_ref_bkp)

                          if resolved_schema:
                            schema = resolved_schema['properties']
                          else:
                            schema = dict(type="object")

                    for prop_items in schema.items():

                        prop_ref = prop_items[1].get("$ref", None)
                        
                        if prop_ref:
                            resolved_prop_ref = resolve_ref(spec, prop_ref)
                            reqbodyserialized = {}
                            
                            for rpr_keys, rpr_value in resolved_prop_ref['properties'].items(): 
                                   
                                reqbodyserialized[rpr_keys] = {    
                                    "name": rpr_keys,
                                    "in": "body",
                                    "description": "null",
                                    "required": True,
                                    "schema": {
                                        "type": rpr_value['type'] or "",
                                    }
                                }

                                if "format" in rpr_value:
                                    reqbodyserialized[rpr_keys]['schema']['format'] = rpr_value['format']

                            # Adiciona os objetos processados à lista
                            all_reqbodyserialized.extend(list(reqbodyserialized.values()))

                        else:
                            reqbodyserialized = {}
                            prop_resolved = ab[first_key]['schema'].get('properties', {})

                            if len(prop_items) == 2:

                                key, value = prop_items

                                if isinstance(value, dict):

                                    reqbodyserialized[key] = {
                                        "name": key,
                                        "in": "body",
                                        "description": value.get('description', 'null'),
                                        "required": True,
                                        "schema": {
                                            "type": value.get('type', ''),
                                        }
                                    }

                                    if "format" in value:
                                        reqbodyserialized[key]['schema']['format'] = value['format']

                            # Adiciona os objetos processados à lista
                            all_reqbodyserialized.extend(list(reqbodyserialized.values()))


        if consumes_type is None:
            consumes_type = "PLACEHOLDER"

        if required_consumes:
            return [consumes_type]

        else:
            if not any(all_reqbodyserialized):
                pass

            return all_reqbodyserialized

        #
        # The real one is Dead... no one can survive a head shot...
        #

    except Exception as e:
        None
        print("RequestBody Error: {}{}{} Just going around it... We are smart as fuck, don't you think?".format(RED, e, NRM))
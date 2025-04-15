# resolve JSON Reference
def resolve_ref(spec, ref_name):
    """
    Resolve and replace a specs JSON schema $ref value
    :param spec: the swagger/openapi spec being parsed
    :param ref_name: the $ref value to be resolved
    :return: None if unable to parse the $ref value, otherwise will resolve to that ref definition
    """

    # don't parse any external definitions
    if ref_name[0:1] == "./" or ".json" in ref_name:
        return

    resolved_ref = None

    for part in ref_name.split('/'):
        if part == "#":
            resolved_ref = spec
            continue

        if resolved_ref is None:
            return

        resolved_ref = resolved_ref.get(part)

    return resolved_ref

# resolve Schema Object
def resolve_schema_object(spec, schema):
    if type(schema) != dict:
        return

    # handle Schema $ref (but don't return, we may need to resolve nested $refs)
    schema_ref = schema.get("$ref", None)
    if schema_ref:
        resolved_schema = resolve_ref(spec, schema_ref)

        if resolved_schema:
            schema = resolved_schema
        else:
            schema = dict(type="object")

    # handle Schema array items
    items = schema.get("items", {})
    if items and type(items) == dict:
        resolved_item = resolve_schema_object(spec, items)

        if resolved_item:
            schema.update({"items": resolved_item})

    # handle Schema object additionalProperties
    additional_properties = schema.get("additionalProperties", {})
    if additional_properties and type(additional_properties) == dict:
        resolved = resolve_schema_object(spec, additional_properties)

        if resolved:
            schema.update({"additionalProperties": resolved})

    # handle Schema object properties
    properties = schema.get("properties", {})
    if properties and type(properties) == dict:
        parsed_properties = {}

        for prop_name, prop_value in properties.items():
            # why is this even necessary ugh
            if prop_name == "$ref":
                print("{}Error{} {}$ref{}".format(RED, NRM, YEL, NRM))
                continue

            resolved = resolve_schema_object(spec, prop_value)

            if resolved:
                parsed_properties[prop_name] = resolved

        if parsed_properties:
            schema.update({"properties": parsed_properties})

    # resolve and replace any $refs in allOf
    all_of = schema.get("allOf", [])
    if all_of and type(all_of) == list:
        for idx, item in enumerate(all_of):
            resolved_item = resolve_schema_object(spec, item)

            if resolved_item:
                all_of[idx] = resolved_item

        schema.update({"allOf": all_of})

    # resolve and replace any $refs in oneOf
    one_of = schema.get("oneOf", [])
    if one_of and type(one_of) == list:
        for idx, item in enumerate(one_of):
            resolved_item = resolve_schema_object(spec, item)

            if resolved_item:
                one_of[idx] = resolved_item

        schema.update({"oneOf": one_of})

    # resolve and replace any $refs in anyOf
    any_of = schema.get("anyOf", [])
    if any_of and type(any_of) == list:
        for idx, item in enumerate(any_of):
            resolved_item = resolve_schema_object(spec, item)

            if resolved_item:
                any_of[idx] = resolved_item

        schema.update({"oneOf": any_of})

    # handle examples (limit to only one, return str)
    example = schema.get("example", None)
    if example:
        if type(example) == dict:
            example = json.dumps(example)
        elif type(example) == list and len(example) >= 1:
            example = example[0]
        else:
            example = str(example)

        schema.update({"example": example})

    return schema
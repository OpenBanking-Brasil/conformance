import argparse
import os
import json
import re


def validate_json(json_obj, api_name):
    uri_regex = r'^https://web\.conformance\.directory\.openbankingbrasil\.org\.br/plan-detail\.html\?.*$'

    def is_valid_api(api):
        # TODO: check what the standard is for OL files
        return api.removesuffix('-OL') == api_name.removesuffix('-OL')

    def is_valid_sd(sd):
        return re.match(r'^\d+$', sd)

    def is_valid_docusign_id(docusign_id):
        return re.match(r'^[a-zA-Z\d]{8}-[a-zA-Z\d]{4}-[a-zA-Z\d]{4}-[a-zA-Z\d]{4}-[a-zA-Z\d]{12}$', docusign_id)

    def is_valid_test_plan_uri(test_plan_uri):
        if isinstance(test_plan_uri, str):
            return re.match(uri_regex, test_plan_uri)
        elif isinstance(test_plan_uri, list):
            return all(re.match(uri_regex, uri) for uri in test_plan_uri)
        return False

    if len(json_obj) != 4:
        return False, 'json contains more keys than expected'
    if 'api' not in json_obj or not is_valid_api(json_obj['api']):
        return False, 'api field has an invalid value'
    if 'sd' not in json_obj or not is_valid_sd(json_obj['sd']):
        return False, 'sd field has an invalid value'
    if 'docusign_id' not in json_obj or not is_valid_docusign_id(json_obj['docusign_id']):
        return False, 'docusign_id field has an invalid value'
    if 'test_plan_uri' not in json_obj or not is_valid_test_plan_uri(json_obj['test_plan_uri']):
        return False, 'test_plan_uri field has an invalid value'
    return True, 'approved'


def main(argv=None):
    parser = argparse.ArgumentParser(
        description="Checks if the json objects are valid"
    )
    parser.add_argument(
        "apis",
        nargs="+",
        help="Every api with version to be checked by the script (<api_name>/<full_version>/, for example payments/2.0.0/)."
    )
    args = parser.parse_args(argv)
        
    base_path = './submissions/functional/'
    post_fixes = args.apis
    directories = [os.path.join(base_path, post_fix) for post_fix in post_fixes]
    wrong_files = []

    for directory in directories:
        for filename in os.listdir(directory):
            if filename.endswith(".json"):
                file_path = os.path.join(directory, filename)
                name_parts = filename.split('_')
                api_name = f"{name_parts[2]}_{name_parts[3]}"
                with open(file_path) as f:
                    json_obj = json.load(f)
                    approved, message = validate_json(json_obj, api_name)
                    if not approved:
                        wrong_files.append((filename, message))

    if wrong_files:
        print('The following files contain invalid json objects:')
        for filename, message in wrong_files:
            print(f'{filename} ({message})')
        return 1

    return 0


if __name__ == '__main__':
    raise SystemExit(main())

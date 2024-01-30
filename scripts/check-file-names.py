import argparse
from os import listdir
import re


def check_file_names(api_list):
    wrong_files = []

    for api_with_version in api_list:
        api, version = api_with_version.split('_v')
        if api == 'payments_webhook':
            api = 'webhook'
        
        directory = f"./submissions/functional/{api}/{version}"
        
        pattern = re.compile(r"^\d{8}_.+_(?P<api>[A-Za-z-]+)_v(?P<version>[\d\.]+)(-OL)?_(0[1-9]|[12]\d|3[01])-(0[1-9]|1[012])-(20\d\d).(zip|ZIP|json)$")

        for file in listdir(directory):
            match = pattern.match(file)
            if (
                file != ".DS_Store" and file != "readme.md" and
                (len(file.split('_')) != 5 or not match or
                match.group('api').removesuffix('-timezone') != api or
                not version.startswith(match.group('version')))
            ):
                wrong_files.append(file)

    return wrong_files


def main(argv=None):
    parser = argparse.ArgumentParser(
        description="Checks if the file names are correct"
    )
    parser.add_argument(
        "apis",
        nargs="+",
        help="Every api with version to be checked by the script (<api_name>_v<full_version>, for example resources_v1.0.0)."
    )
    args = parser.parse_args(argv)
    
    api_list = args.apis

    wrong_files = check_file_names(api_list)

    if wrong_files:
        print("The following file names are wrong: ", wrong_files)
        return 1
        
    return 0


if __name__ == '__main__':
    raise SystemExit(main())

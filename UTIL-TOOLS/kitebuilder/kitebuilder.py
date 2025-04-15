# encoding=utf-8


import logging
import argparse
import base64

from convert import convert_csv
from parse import parse_specs

import sys

def main():
    root_parser = argparse.ArgumentParser(description="OpenAPI/Swagger API schema parser originally made by Assetnote")
    action_subparser = root_parser.add_subparsers(title="action", dest="action")

    parse_parser = action_subparser.add_parser(
        "parse",
        help="Parse a directory of OpenAPI/Swagger files into a single JSON file for Shadow-Spotter"
    )
    parse_parser.add_argument(
        "--blacklist",
        metavar="HOSTS",
        type=list,
        default=[
            "googleapis",
            "azure",
            "petstore",
            "amazon"
        ],
        help="Exclude specs with host field matching any of these strings (default googleapis, azure, petstore, amazon)"
    )
    parse_parser.add_argument(
        "--scrape-dir",
        metavar="DIR",
        type=str,
        help="Directory to read list of specs from (default ./scrape)",
        default="./scrape"
    )
    parse_parser.add_argument(
        "--output-file",
        metavar="FILE",
        type=str,
        help="File to output resulting schema to (default output.json)",
        default="output.json"
    )

    convert_parser = action_subparser.add_parser(
        "convert",
        help="Convert a file to a number of swagger JSON files in the provided output directory"
    )
    convert_parser.add_argument(
        "--file",
        metavar="FILE",
        type=str,
        help="File to convert to a number of swagger spec files.",
        required=True
    )
    convert_parser.add_argument(
        "--format",
        metavar="FORMAT",
        default="CSV",
        choices=["CSV"],
        type=str,
        help="File format to convert. Only CSV files supported. Must be in the format 'id,name,content'",
    )
    convert_parser.add_argument(
        "--scrape-dir",
        metavar="DIR",
        type=str,
        help="File to output resulting schema files to (defaults to ./scrape)",
        default="./scrape"
    )

    if len(sys.argv)==1:
       root_parser.print_help(sys.stderr)
       sys.exit(1)

    args = root_parser.parse_args()

    logging.basicConfig(
        level=logging.INFO,
        format="[%(asctime)s] %(message)s",
        datefmt="%m/%d/%Y %I:%M:%S %p"
    )

    if args.action == "parse":
        spec_count = parse_specs(args.scrape_dir, args.output_file, args.blacklist)
        logging.info(f"Finished parsing {spec_count}")

    elif args.action == "convert":
        if args.format == "CSV":
            convert_csv(args.file, args.scrape_dir)


if __name__ == "__main__":

    banner_kitebuilder = "4qCA4qCA4qCA4qCA4qCA4qCA4qCA4qCA4qCA4qCA4qCA4qCA4qCA4qCA4qCA4qCA4qCA4qCA4qCA4qCA4qCA4qCA4qCA4qCA4qCA4qCA4qCA4qOg4qOE4qCA4qCA4qCA4qCA4qCA4qCA4qCA4qCA4qCA4qCA4qCA4qCA4qCA4qCA4qCA4qCA4qCA4qCA4qCA4qCA4qCA4qCA4qCA4qCA4qCA4qCA4qCACuKggOKggOKggOKggOKggOKggOKggOKggOKggOKggOKggOKggOKggOKggOKggOKggOKggOKggOKggOKggOKggOKggOKggOKggOKjoOKjtOKjv+Kjv+Kjv+Kjv+KjpuKjhOKggOKggOKggOKggOKggOKggOKggOKggOKggOKggOKggOKggOKggOKggOKggOKggOKggOKggOKggOKggOKggOKggOKggOKggArioIDioIDioIDioIDioIDioIDioIDioIDioIDioIDioIDioIDioIDioIDioIDioIDioIDioIDioIDioIDioIDio6Dio7Tio7/io7/ioJ/ioIvioIDioIDioJnioLvio7/io7/io6bio4TioIDioIDioIDioIDioIDioIDioIDioIDioIDioIDioIDioIDioIDioIDioIDioIDioIDioIDioIDioIDioIAK4qCA4qCA4qCA4qCA4qCA4qCA4qCA4qCA4qCA4qCA4qCA4qCA4qCA4qCA4qCA4qCA4qCA4qCA4qOA4qO04qO+4qO/4qC/4qCL4qCA4qCA4qCA4qCA4qCA4qCA4qCA4qCA4qCZ4qC74qO/4qO34qOm4qOA4qCA4qCA4qCA4qCA4qCA4qCA4qCA4qCA4qCA4qCA4qCA4qCA4qCA4qCA4qCA4qCA4qCA4qCACuKggOKggOKggOKggOKggOKggOKggOKggOKggOKggOKggOKggOKggOKggOKggOKigOKjtOKjvuKjv+Kgv+Kgi+KggeKggOKggOKggOKggOKggOKggOKggOKggOKggOKggOKggOKggOKgiOKgmeKgv+Kjv+Kjt+KjpuKhgOKggOKggOKggOKggOKggOKggOKggOKggOKggOKggOKggOKggOKggOKggOKggArioIDioIDioIDioIDioIDioIDioIDioIDioIDioIDioIDioIDiooDio6Tio77io7/iob/ioIvioIHioIDioIDioIDioIDioIDioIDioIDioIDioIDioIDioIDioIDioIDioIDioIDioIDioIDioIDioIjioJnior/io7/io7fio6TioYDioIDioIDioIDioIDioIDioIDioIDioIDioIDioIDioIDioIAK4qCA4qCA4qCA4qCA4qCA4qCA4qCA4qCA4qCA4qKA4qOk4qO+4qO/4qG/4qCb4qCB4qCA4qCA4qOA4qOk4qOk4qOk4qOk4qOA4qCA4qCA4qCA4qCA4qCA4qCA4qCA4qCA4qOA4qOk4qOk4qOk4qOk4qOE4qCA4qCA4qCI4qCb4qK/4qO/4qO34qOk4qGA4qCA4qCA4qCA4qCA4qCA4qCA4qCA4qCA4qCACuKggOKggOKggOKggOKggOKggOKigOKjoOKjtuKjv+Khv+Kgm+KggeKggOKggOKggOKggOKjgOKgu+Kjv+Kjv+Kjv+Kjv+Kjv+Kjv+Kjt+KjtuKjtuKjtuKjtuKjvuKjv+Kjv+Kjv+Kjv+Kjv+Kjv+Kgn+KjgOKggOKggOKggOKggOKgiOKgu+Kiv+Kjv+KjtuKjhOKhgOKggOKggOKggOKggOKggOKggArioIDioIDioIDioIDioIDioLDio7/io7/io5/ioInioIDioIDioIDioIDioIDio7Tio6bio7/io4DioIniorvio7/io7/io7/io7/io7/io7/io7/io7/io7/io7/io7/io7/io7/io7/io5viooHio6Dio7/io7bio6bioIDioIDioIDioIDioIDioInio7/io7/io7/ioIbioIDioIDioIDioIDioIAK4qCA4qCA4qCA4qCA4qCA4qCA4qCI4qK/4qO/4qO34qGA4qCA4qCA4qCA4qCA4qCY4qK/4qO/4qO/4qO/4qO/4qO/4qO/4qO/4qO/4qO/4qO/4qO/4qO/4qO/4qO/4qO/4qO/4qO/4qO/4qO/4qO/4qO/4qO/4qG/4qCD4qCA4qCA4qCA4qCA4qKA4qO+4qO/4qG/4qCB4qCA4qCA4qCA4qCA4qCA4qCACuKggOKggOKggOKggOKggOKggOKggOKggOKgu+Kjv+Kjv+KjhuKggOKggOKggOKgseKjpOKhiOKgmeKgm+Kgv+Kgv+Kgv+Kjv+KjheKjqeKjveKjv+Kjv+Kjr+KjheKjuOKjv+Kgv+Kgv+Kgn+Kgm+Kgi+KigeKjpOKgjuKggOKggOKggOKjsOKjv+Kjv+Kgn+KggOKggOKggOKggOKggOKggOKggOKggArioIDioIDioIDioIDioIDioIDioIDioIDioIDioJjior/io7/io7fioYDioIDioIDioJjior/io4TioIDioLnio7/io7fioYbioInioJnioInioInioIjioInioInioIniorDio77io7/ioIPioIDio6Diob/ioIPioIDioIDiooDio77io7/iob/ioIPioIDioIDioIDioIDioIDioIDioIDioIDioIAK4qCA4qCA4qCA4qCA4qCA4qCA4qCA4qCA4qCA4qCA4qCA4qC74qO/4qO/4qOG4qCA4qCA4qCA4qC74qOG4qCA4qCI4qC74qOH4qCA4qCA4qCA4qCA4qCA4qCA4qCA4qCA4qO44qCf4qCB4qCA4qOw4qCf4qCA4qCA4qCA4qOw4qO/4qO/4qCf4qCA4qCA4qCA4qCA4qCA4qCA4qCA4qCA4qCA4qCA4qCACuKggOKggOKggOKggOKggOKggOKggOKggOKggOKggOKggOKggOKgmOKiv+Kjv+Kjp+KhgOKggOKggOKgmeKjp+KggOKggOKgiOKggOKggOKggOKggOKggOKggOKggOKggOKggeKggOKggOKjvOKgi+KggOKggOKigOKjvuKjv+Khv+Kgg+KggOKggOKggOKggOKggOKggOKggOKggOKggOKggOKggOKggArioIDioIDioIDioIDioIDioIDioIDioIDioIDioIDioIDioIDioIDioIDioLvio7/io7/io4TioIDioIDioJjioIfioIDioIDioIDioIDioIDioIDioIDioIDioIDioIDioIDioIDiorjioIPioIDioIDio6Dio7/io7/ioJ/ioIDioIDioIDioIDioIDioIDioIDioIDioIDioIDioIDioIDioIDioIAK4qCA4qCA4qCA4qCA4qCA4qCA4qCA4qCA4qCA4qCA4qCA4qCA4qCA4qCA4qCA4qCY4qK/4qO/4qOn4qGA4qCA4qCA4qCA4qCA4qOg4qCA4qCA4qCA4qCA4qCA4qCA4qGE4qCA4qCA4qCA4qCA4qKA4qO84qO/4qG/4qCD4qCA4qCA4qCA4qCA4qCA4qCA4qCA4qCA4qCA4qCA4qCA4qCA4qCA4qCA4qCACuKggOKggOKggOKggOKggOKggOKggOKggOKggOKggOKggOKggOKggOKggOKggOKggOKggOKgu+Kjv+Kjv+KjhOKggOKggOKggOKiueKjt+KjpOKjgOKjgOKjpOKjvuKgh+KggOKggOKggOKjoOKjv+Kjv+Kgn+KggOKggOKggOKggOKggOKggOKggOKggOKggOKggOKggOKggOKggOKggOKggOKggOKggArioIDioIDioIDioIDioIDioIDioIDioIDioIDioIDioIDioIDioIDioIDioIDioIDioIDioIDioJnior/io7/io6fioYDioIDioIDioJnioJvioL/ioL/ioJvioInioIDioIDiooDio7zio7/iob/ioIvioIDioIDioIDioIDioIDioIDioIDioIDioIDioIDioIDioIDioIDioIDioIDioIDioIDioIAK4qCA4qCA4qCA4qCA4qCA4qCA4qCA4qCA4qCA4qCA4qCA4qCA4qCA4qCA4qCA4qCA4qCA4qCA4qCA4qCI4qC74qO/4qO/4qOE4qCA4qCA4qCA4qCA4qCA4qCA4qCA4qCA4qOg4qO/4qO/4qCf4qCB4qCA4qCA4qCA4qCA4qCA4qCA4qCA4qCA4qCA4qCA4qCA4qCA4qCA4qCA4qCA4qCA4qCA4qCA4qCACuKggOKggOKggOKggOKggOKggOKggOKggOKggOKggOKggOKggOKggOKggOKggOKggOKggOKggOKggOKggOKggOKgmeKiv+Kjv+Kjp+KhgOKggOKggOKggOKggOKigOKjvOKjv+Khv+Kgi+KggOKggOKggOKggOKggOKggOKggOKggOKggOKggOKggOKggOKggOKggOKggOKggOKggOKggOKggOKggOKggArioIDioIDioIDioIDioIDioIDioIDioIDioIDioIDioIDioIDioIDioIDioIDioIDioIDioIDioIDioIDioIDioIDioIjioLvio7/io7/io4TioIDioIDio6Dio7/io7/ioJ/ioIHioIDioIDioIDioIDioIDioIDioIDioIDioIDioIDioIDioIDioIDioIDioIDioIDioIDioIDioIDioIDioIDioIAK4qCA4qCA4qCA4qCA4qCA4qCA4qCA4qCA4qCA4qCA4qCA4qCA4qCA4qCA4qCA4qCA4qCA4qCA4qCA4qCA4qCA4qCA4qCA4qCA4qCZ4qK/4qO/4qOn4qO84qO/4qG/4qCL4qCA4qCA4qCA4qCA4qCA4qCA4qCA4qCA4qCA4qCA4qCA4qCA4qCA4qCA4qCA4qCA4qCA4qCA4qCA4qCA4qCA4qCA4qCA4qCACuKggOKggOKggOKggOKggOKggOKggOKggOKggOKggOKggOKggOKggOKggOKggOKggOKggOKggOKggOKggOKggOKggOKggOKggOKggOKgiOKgu+Kjv+Kjv+Khn+KggeKggOKggOKggOKggOKggOKggOKggOKggOKggOKggOKggOKggOKggOKggOKggOKggOKggOKggOKggOKggOKggOKggOKggOKggOKggArioIDioIDioIDioIDioIDioIDioIDioIDioIDioIDioIDioIDioIDioIDioIDioIDioIDioIDioIDioIDioIDioIDioIDioIDioIDioIDioIDioJnioIvioIDioIDioIDioIDioIDioIDioIDioIDioIDioIDioIDioIDioIDioIDioIDioIDioIDioIDioIDioIDioIDioIDioIDioIDioIDioIDioIAK"
    
    print("\001\033[1;94m\002")
    print(base64.b64decode(banner_kitebuilder).decode('utf-8'))
    print("    \033[1;97mCopyright\033[0m \033[1;92m©\033[0m \033[1;96m2024\033[0m \033[1;92mSNAKE Security\033[0m \033[1;97m-\033[0m \033[1;91mJust one bite\033[0m\n")
    
    main()
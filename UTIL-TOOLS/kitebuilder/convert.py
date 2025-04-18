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

import csv
import sys
from os import listdir
from os.path import isfile, join
import logging


def convert_csv(csv_file, scrape_dir):
    csv.field_size_limit(sys.maxsize)

    existing_files = [f.rsplit('.', 1)[1] for f in listdir(scrape_dir) if isfile(join(scrape_dir, f))]
    idx = max(int(x) for x in existing_files if x.isdigit()) + 1

    with open(csv_file) as file:
        reader = csv.reader(file)

        for row in reader:
            current_file = f"{scrape_dir}/openapi.json.{idx}"

            try:
                with open(current_file, "x") as f:
                    f.write(row[2])

                logging.info(f"Wrote {current_file} from csv: {csv_file}")
            except FileExistsError:
                logging.warning(f"Failed to write {current_file} from csv: {csv_file}")
                continue
            finally:
                idx += 1

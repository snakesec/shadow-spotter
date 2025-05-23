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

import requests
from requests.exceptions import ConnectionError, HTTPError
from json.decoder import JSONDecodeError
from math import ceil
import asyncio
import concurrent.futures
import aiohttp
from aiohttp import ClientSession
import time
import os

#
#
# This script is a piece of shit, a big one, don't use this crap...
#
#

def _make_request(url):
    return requests.get(url).json()


def get_spec_count():
    try:
        return requests.get("https://app.swaggerhub.com/apiproxy/specs?limit=1&page=1").json().get("totalCount")

    except ConnectionError:
        print(f"Failed to get spec count: Connection Failed")
        return

    except JSONDecodeError:
        print(f"Failed to get spec count: Invalid JSON")
        return

    except Exception as e:
        print(f"Failed to get spec count: {e}")
        return


async def get_spec_list(session, page, sort_by="BEST_MATCH", order="ASC", limit=100):
    try:
        response = await session.request(method='GET', url=(
            f"https://app.swaggerhub.com/apiproxy/specs"
            f"?sort={sort_by}&order={order}&limit={limit}&page={page}"
        ))
        response.raise_for_status()
        response_json = await response.json()
        return list(map(lambda api: api["properties"][0]["url"], response_json.get("apis")))

    except HTTPError as http_err:
        print(f"Failed to get spec list: HTTP error occurred ({http_err})")

    except Exception as err:
        print(f"An error occurred: {err}")


"""
async def save_spec(session, url, file_name):
    try:
        response = await session.request(method='GET', url=url)
        response.raise_for_status()
        response_text = await response.text()

        with open(f"scrape/swaggerhub/{file_name}", "w+") as f:
            f.write(response_text)

        print(f"Saved new swagger spec: {file_name}")
        time.sleep(3)

    except HTTPError as http_err:
        print(f"Failed to save file ({file_name}): HTTP Error ({http_err})")
    except ConnectionError:
        print(f"Failed to save file ({file_name}): Connection Failed")
    except JSONDecodeError:
        print(f"Failed to save file ({file_name}): Invalid JSON")
    except Exception as e:
        print(f"Failed to save file ({file_name}): {e}")
"""

import os

async def save_spec(session, url, file_name):
    try:
        if os.path.exists(f"scrape/swaggerhub/{file_name}"):
            print(f"File {file_name} already exists. Skipping...")
            return

        response = await session.request(method='GET', url=url)
        response.raise_for_status()
        response_text = await response.text()

        with open(f"scrape/swaggerhub/{file_name}", "w+") as f:
            f.write(response_text)

        print(f"Saved new swagger spec: {file_name}")
        time.sleep(3)

    except HTTPError as http_err:
        print(f"Failed to save file ({file_name}): HTTP Error ({http_err})")
    except ConnectionError:
        print(f"Failed to save file ({file_name}): Connection Failed")
    except JSONDecodeError:
        print(f"Failed to save file ({file_name}): Invalid JSON")
    except Exception as e:
        print(f"Failed to save file ({file_name}): {e}")


async def save_spec_page(session, page):
    print(f"Saving page: {page}")

    try:
        urls = await get_spec_list(session, page)

        for idx, url in enumerate(urls):
            time.sleep(1)
            await save_spec(session, url, f"swagger.json.{((page-1)*100)+idx}")

    except Exception as err:
        print(f"Exception: {err}")
        pass


async def main():
    async with ClientSession() as session:
        spec_count = get_spec_count()
        page_count = ceil(spec_count / 100)

        print(f"Collecting {spec_count} specs (pages: {page_count})")
        await asyncio.gather(*[save_spec_page(session, page) for page in range(1, 99)])


loop = asyncio.get_event_loop()
loop.run_until_complete(main())

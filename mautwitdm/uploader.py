# Copyright (c) 2020 Tulir Asokan
#
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at http://mozilla.org/MPL/2.0/.
from typing import Dict
import logging
import asyncio
import math

from aiohttp import ClientSession, MultipartWriter
from yarl import URL

from .types import MediaUploadResponse
from .errors import check_error


class TwitterUploader:
    upload_url: URL = URL("https://upload.twitter.com/i/media/upload.json")

    http: ClientSession
    log: logging.Logger
    headers: Dict[str, str]

    async def _wait_processing(self, media_id: str, wait_requests: int = 1) -> MediaUploadResponse:
        query_req = self.upload_url.with_query({
            "command": "STATUS",
            "media_id": media_id
        })
        while True:
            await asyncio.sleep(wait_requests)
            async with self.http.get(query_req, headers=self.headers) as resp:
                data = await check_error(resp)
            state = data["processing_info"]["state"]
            if state == "succeeded":
                self.log.debug(f"Server completed processing of {media_id}")
                return MediaUploadResponse.deserialize(data)
            elif state == "in_progress":
                progress = data["processing_info"]["progress_percent"]
                wait_requests = data["processing_info"].get("check_after_secs", wait_requests)
                self.log.debug(f"Upload of {media_id} at {progress} %, "
                               f"re-checking after {wait_requests} seconds")
            else:
                raise RuntimeError(f"Unknown state {state}")

    async def _upload_data(self, media_id: str, data: bytes) -> MediaUploadResponse:
        max_size = 2 ** 17
        base_upload_req = self.upload_url.with_query({
            "command": "APPEND",
            "media_id": media_id,
        })
        for i in range(math.ceil(len(data) / max_size)):
            multipart_data = MultipartWriter("form-data")
            part = multipart_data.append(data[i * max_size: (i + 1) * max_size])
            part.set_content_disposition("form-data", name="media", filename="blob")
            req = base_upload_req.update_query({"segment_index": i})
            async with self.http.post(req, data=multipart_data, headers=self.headers) as resp:
                await check_error(resp)
                self.log.debug(f"Uploaded segment {i} of {media_id}")
        finalize_req = self.upload_url.with_query({
            "command": "FINALIZE",
            "media_id": media_id,
        })
        async with self.http.post(finalize_req, headers=self.headers) as resp:
            resp_data = await check_error(resp)
        processing_info = resp_data.get("processing_info", {})
        if processing_info.get("state", None) == "pending":
            self.log.debug(f"Finished uploading {media_id}, but server is still processing it")
            check_after = processing_info.get("check_after_secs", 1)
            return await self._wait_processing(media_id, check_after)
        self.log.debug(f"Finished uploading {media_id}")
        return MediaUploadResponse.deserialize(resp_data)

    async def upload(self, data: bytes = None, url: str = None, mime_type: str = None
                     ) -> MediaUploadResponse:
        if mime_type == "image/gif":
            category = "dm_gif"
            size_limit = 15 * 1024 * 1024
        elif mime_type.startswith("image/"):
            category = "dm_image"
            size_limit = 5 * 1024 * 1024
        elif mime_type.startswith("video/"):
            category = "dm_video"
            size_limit = 15 * 1024 * 1024
        else:
            raise ValueError("Unsupported mime type")
        if len(data) > size_limit:
            raise ValueError("File too big")
        init_req = {
            "command": "INIT",
            "media_type": mime_type,
            "media_category": category,
        }
        if data is not None:
            init_req["total_bytes"] = len(data)
        elif url is not None:
            init_req["source_url"] = url
        else:
            raise ValueError("Either data bytes or url must be provided")
        init_url = self.upload_url.with_query(init_req)
        async with self.http.post(init_url, headers=self.headers) as resp:
            resp_data = await check_error(resp)
        media_id = resp_data["media_id_string"]
        self.log.debug(f"Started upload, got media ID {media_id}")

        if url is not None:
            processing_info = resp_data.get("processing_info", {})
            if processing_info.get("state") == "succeeded":
                return MediaUploadResponse.deserialize(data)
            else:
                check_after = processing_info.get("check_after_secs", 1)
                return await self._wait_processing(media_id, check_after)
        else:
            return await self._upload_data(media_id, data)

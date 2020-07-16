# Copyright (c) 2020 Tulir Asokan
#
# This Source Code Form is subject to the terms of the Mozilla Public
# License, v. 2.0. If a copy of the MPL was not distributed with this
# file, You can obtain one at http://mozilla.org/MPL/2.0/.
import logging
import asyncio
import math

from aiohttp import ClientSession, MultipartWriter
from yarl import URL

from .types import MediaUploadResponse


class TwitterUploader:
    upload_url: URL
    http: ClientSession
    log: logging.Logger

    async def _wait_processing(self, media_id: str, wait_requests: int = 1) -> MediaUploadResponse:
        query_req = self.upload_url.with_query({
            "command": "STATUS",
            "media_id": media_id
        })
        while True:
            await asyncio.sleep(wait_requests)
            async with self.http.get(query_req) as resp:
                resp.raise_for_status()
                data = await resp.json()
            state = data["processing_info"]["state"]
            if state == "succeeded":
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
            multipart_data.set_content_disposition("form-data", name="media", filename="blob")
            multipart_data.append(data[i * max_size: (i + 1) * max_size])
            upload_req = base_upload_req.update_query({"segment_index": i})
            async with self.http.post(upload_req, data=multipart_data) as resp:
                resp.raise_for_status()
        finalize_req = self.upload_url.with_query({
            "command": "FINALIZE",
            "media_id": media_id,
        })
        async with self.http.post(finalize_req) as resp:
            resp.raise_for_status()
            resp_data = await resp.json()
        processing_info = resp_data.get("processing_info", {})
        if processing_info.get("state", None) == "pending":
            check_after = processing_info.get("check_after_secs", 1)
            return await self._wait_processing(media_id, check_after)
        return MediaUploadResponse.deserialize(resp_data)

    async def upload(self, data: bytes = None, url: str = None, mime_type: str = None
                     ) -> MediaUploadResponse:
        if mime_type == "image/gif":
            category = "dm_gif"
        elif mime_type.startswith("image/"):
            category = "dm_image"
        elif mime_type.startswith("video/"):
            category = "dm_video"
        else:
            raise ValueError("Unsupported mime type")
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
        async with self.http.post(self.upload_url.with_query(init_req)) as resp:
            resp.raise_for_status()
            resp_data = await resp.json()
        media_id = resp_data["media_id_string"]

        if url is not None:
            processing_info = resp_data.get("processing_info", {})
            if processing_info.get("state") == "succeeded":
                return MediaUploadResponse.deserialize(data)
            else:
                check_after = processing_info.get("check_after_secs", 1)
                return await self._wait_processing(media_id, check_after)
        else:
            return await self._upload_data(media_id, data)

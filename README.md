# StreamSight Development Setup

## Overview
StreamSight uses MediaMTX (RTSP server), FFmpeg (streaming), and a Go Worker for real-time video detection.

---

## Prerequisites
- Docker installed
- FFmpeg (can be run via Docker)
- Video file `video.mp4` for streaming tests

---

## Run MediaMTX RTSP Server

Start MediaMTX server on your local machine with:

```docker run --rm -it -p 8554:8554 bluenviron/mediamtx:latest```

text

This exposes an RTSP server on port 8554.

---

## Stream Video File to MediaMTX with FFmpeg

Use the FFmpeg Docker image to stream your local video file continuously:

docker run --rm -v <path>:/videos jrottenberg/ffmpeg
-re -stream_loop -1 -i /videos/video.mp4 -c copy -f rtsp -rtsp_transport tcp rtsp://host.docker.internal:8554/mystream

text

- Adjust `<path>` to your own path.
- The stream URL will be: `rtsp://localhost:8554/mystream`

---

## Test RTSP Stream Playback

Verify your stream in VLC or with `ffplay`:

vlc rtsp://localhost:8554/mystream

ffplay rtsp://localhost:8554/mystream

---

## Configure Worker

Set your worker environment variable to consume the stream:

```RTSP_STREAMS=rtsp://localhost:8554/mystream```

text

---

## Common Docker Commands

- Run MediaMTX:

```docker run --rm -it -p 8554:8554 bluenviron/mediamtx:latest```

- Run FFmpeg streaming:

```docker run --rm -v PATH_TO_VIDEO:/videos jrottenberg/ffmpeg -re -stream_loop -1 -i /videos/video.mp4 -c copy -f rtsp -rtsp_transport tcp rtsp://host.docker.internal:8554/mystream```

---

- Build and run worker on docker

```docker build -t streamsight-worker .```
```docker run --rm --network="host" streamsight-worker```

## Notes

- Use `-rtsp_transport tcp` in FFmpeg commands to avoid broken pipe errors.
- This setup provides a reliable RTSP streaming pipeline for local development and testing.

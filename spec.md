# Spec

## Directories

- Video list JSONs are downloaded to disk (and optionally compressed with brotli)
- We only save the audio/video format specified in `format_id`
- We only save the thumbnail in WebP
- This gets turned into a torrent cuz why not lol

```
dir/
├─ lists/
│  ├─ 2021-08-15T02.30.45Z.json
│  ├─ UCdn5BQ06XqgXoAxIhbqw5Rg.json
│  ├─ local.json
├─ videos/
│  ├─ UCdn5BQ06XqgXoAxIhbqw5Rg/
│  │  ├─ gekKWg0yoOE/
│  │  │  ├─ f140.m4a
│  │  │  ├─ thumbnail.webp
│  │  │  ├─ f248.webm
│  │  │  ├─ info.json
├─ torrents/
|  ├─ UCdn5BQ06XqgXoAxIhbqw5Rg.torrent
```

## Lists database

```json
[
    {
        "c": "UCdn5BQ06XqgXoAxIhbqw5Rg",
        "v": "gekKWg0yoOE",
        "u": "2021-08-15T02.30.45Z",
    }
]
```

## Downloading
- FIFO queue via https://github.com/enriquebris/goconcurrentqueue
- goroutine stuff ig?
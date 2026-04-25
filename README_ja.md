# llama-run

[English README](./README.md)

`llama-run` は [llama.cpp](https://github.com/ggml-org/llama.cpp) の `llama-server` を、YAML による宣言的な設定で起動できるようにするラッパーです。

`llama-server` が持つ多数のコマンドラインフラグを構造化された YAML に落とし込むことで、コンテナや systemd 環境で扱いやすい構成管理を可能にします。設定の読み込み・モデルのダウンロード・マルチモデル用 preset 生成を行ったのち、`syscall.Exec(2)` により自身を `llama-server` プロセスへ置き換えます。ラッパープロセスは残らず、PID をそのまま `llama-server` が引き継ぐため、コンテナ内 PID 1 としての実行にも適しています。

## 主な機能

- **YAML による宣言的設定** — `llama-server` のフラグを構造化された YAML で記述できます。
- **`syscall.Exec` によるプロセス置換** — ラッパープロセスは残らず、`llama-server` が直接 PID を引き継ぎます。Docker / Kubernetes / systemd での運用に適しています。
- **モデルの自動ダウンロード** — Hugging Face および S3 (S3 互換ストレージを含む) からモデルを取得します。ダウンロードしたモデルファイルは SHA-256 による検証を行います。
- **マルチモデルルータ対応** — 複数のモデルにエイリアスを割り当て、`llama-server` の `--models-preset` 機能を使った切り替えが行えます。`preset.ini` は自動生成されます。
- **cgroup 対応の CPU スレッド自動検出** — `threads: Auto` 指定時、`/sys/fs/cgroup/cpu.max` を参照してコンテナに割り当てられた CPU 量を考慮します。

## インストール

### Docker イメージ

`llama-server` 本体も同梱された Docker イメージを配布予定です。設定ファイルをマウントして起動してください。

```sh
docker run --rm -p 8080:8080 \
  -v $(pwd)/config.yaml:/etc/llama-run/config.yaml \
  -v llama-run-data:/var/lib/llama-server \
  ghcr.io/sileader/llama-run:latest
```

### バイナリ

GitHub Releases にてバイナリを配布予定です。`llama-server` は別途インストールし、`config.yaml` の `llamaServer.executable` でそのパスを指定してください。

### ソースからビルド

```sh
go build -o llama-run .
```

Go 1.25 以上が必要です。

## 使い方

```sh
llama-run --config /etc/llama-run/config.yaml
```

`--config` を省略した場合、デフォルトで `/etc/llama-run/config.yaml` が読み込まれます。

## 設定

以下はすべての項目を含んだサンプルです。未指定の項目は既定値が使われます。

```yaml
llamaServer:
  executable: /usr/bin/llama-server
  arguments: []
  directory:
    model: /var/lib/llama-server/model
    config: /var/lib/llama-server/config

downloader:
  s3:
    region: us-east-1
    endpoint: null           # S3 互換ストレージを使う場合のみ指定
    accessKeyEnv: null       # アクセスキーを読み込む環境変数名 (省略時 AWS_ACCESS_KEY_ID)
    secretKeyEnv: null       # シークレットキーを読み込む環境変数名 (省略時 AWS_SECRET_ACCESS_KEY)

server:
  host: 0.0.0.0
  port: 8080
  reusePort: false
  apiPrefix: null
  staticPath: null
  unsafeApiKey: []
  apiKeyFile: null
  tls: null                  # { certFile: ..., keyFile: ... }

features:
  embedding:
    enabled: false
    pooling: null            # None | Mean | Cls | Last | Rank
  rerank:
    enabled: false
  webui:
    enabled: false
    config: null
    configFile: null
  metrics:
    enabled: false
  properties:
    enabled: false
  jinja:
    enabled: true

log:
  enabled: true
  file: null
  level: Info                # Debug | Info | Warn | Error | Generic
  timestamp: true
  color: Auto                # Auto | On | Off

chat:
  template: null
  templateFile: null
  templateArguments: {}

reasoning:
  mode: Auto                 # Auto | On | Off
  format: None               # None | Deepseek | DeepseekLegacy
  budget: Unrestricted       # Unrestricted | Immediate | <数値>
  budgetMessage: ""

device:
  cpu:
    threads: Auto            # Auto | <数値>
  memory:
    mmap: true
  gpu:
    layers: Auto             # Auto | All | <数値>
    mainIndex: 0

model:
  alias: null
  aliases: []
  docker: null               # llama-server の --docker-repo に渡す値
  huggingFace: null          # llama-server の --hf-repo に渡す値
  router:
    default:
      context: null
      gpuLayers: null
    models:
      - alias: gemma4-e4b
        huggingFace: ggml-org/gemma-4-E4B-it-GGUF:Q4_K_M

sampling:
  samplers: null
  seed: Random               # Random | <数値>
  temperature: null
  topK: null
  topP: null
  minP: null
  repeatLastN: 64            # <数値> | Disabled | Context
  repeatPenalty: Disabled    # Disabled | <数値>
  presencePenalty: Disabled  # Disabled | <数値>
  frequencyPenalty: Disabled # Disabled | <数値>
```

### `llamaServer`

起動対象となる `llama-server` バイナリのパスと、モデル・設定ファイルを配置するディレクトリを指定します。`directory.model` はダウンロードしたモデルの保存先、`directory.config` は自動生成される `preset.ini` の出力先です。

### `model.router`

マルチモデル運用のための設定です。`models` 配下の各エイリアスについて、以下の処理が行われます。

1. `huggingFace` または `s3` で指定されたソースからモデルを `directory.model/<alias>.gguf` へダウンロード。
2. エイリアスと対応するモデルパスを含む `preset.ini` を `directory.config/preset.ini` に生成。
3. `llama-server` に `--models-preset <path>` を付与して起動。

エイリアスごとに個別の `context` を設定できるほか、`default` ブロックでルータ全体の既定値 (`context`, `gpuLayers`) を指定できます。

## モデルのダウンロードソース

### Hugging Face

`huggingFace: <org>/<repo>:<quantize>` の形式で指定します（例: `ggml-org/gemma-4-E4B-it-GGUF:Q4_K_M`）。認証が必要な場合は環境変数 `HF_TOKEN` にトークンを設定してください。ダウンロード後は Hugging Face が返す SHA-256 と照合した上で確定されます。

### S3 / S3 互換ストレージ

`s3://<bucket>/<key>` の形式で指定します。

- **AWS S3** — `downloader.s3.endpoint` を `null` のままにしておくと、`aws-sdk-go-v2` の標準的な認証解決 (環境変数・共有 config・IAM ロール等) に従います。
- **S3 互換ストレージ (MinIO 等)** — `endpoint` を指定すると path-style のアクセスに切り替わり、`accessKeyEnv` / `secretKeyEnv` で指定した環境変数から認証情報を読み取ります（未指定の場合は `AWS_ACCESS_KEY_ID` / `AWS_SECRET_ACCESS_KEY`）。

## 対応する `llama-server` のバージョン

本ツールは `llama-server` の以下のような比較的新しい機能を前提としています。

- `--models-preset` (マルチモデル)
- `--reasoning` / `--reasoning-format` / `--reasoning-budget`
- `--hf-repo` / `--docker-repo`
- `--webui` / `--no-webui` / `--webui-config*`

利用する機能に応じて、対応する `llama-server` のバージョンを用意してください。

## ライセンス

MIT License. Copyright (c) SiLeader.

# clscli

Command-line tool for Tencent Cloud CLS (Cloud Log Service). Query and analyze logs, list topics, and fetch log context from the terminal.

## Features

- **List topics** – List log topics by region with optional filters (topic name, logset name/ID) so you can pick the right `--region` and topic for queries.
- **Query logs** – Search with CQL or SQL; time range (`--last`, `--from`/`--to`), multiple topics, pagination, and output as JSON/CSV or to a file.
- **Log context** – Retrieve context around a specific log (by PkgId and PkgLogId from search results).

## Installation

### Homebrew (macOS)

```bash
brew tap dbwang0130/clscli
brew install dbwang0130/clscli/clscli
```

### From source (Go 1.21+)

```bash
go install github.com/clscli/clscli@latest
```

Or clone and build:

```bash
git clone https://github.com/clscli/clscli.git
cd clscli
make build
```

Binary is written to `./clscli` (or `clscli.exe` on Windows).

## Prerequisites

- Tencent Cloud account with CLS enabled.
- API credentials and a CLS region. See [Tencent Cloud API common parameters](https://cloud.tencent.com/document/api/614/56474) for regions and auth.

Set environment variables (same as Tencent Cloud API):

```bash
export TENCENTCLOUD_SECRET_ID="your-secret-id"
export TENCENTCLOUD_SECRET_KEY="your-secret-key"
```

Region is passed per command with `--region` (e.g. `ap-guangzhou`).

## Quick start

1. List topics in a region to get topic IDs:
   ```bash
   clscli topics --region ap-guangzhou
   ```
2. Query logs (e.g. last 1 hour, CQL + optional SQL):
   ```bash
   clscli query -q "level:ERROR" --region ap-guangzhou -t <TopicId> --last 1h
   ```
3. Get context around a log (use PkgId and PkgLogId from query results):
   ```bash
   clscli context <PkgId> <PkgLogId> --region ap-guangzhou -t <TopicId>
   ```

## Usage

| Command   | Description                    |
|----------|---------------------------------|
| `topics` | List log topics (with filters) |
| `query`  | Search/analyze logs (CQL/SQL)  |
| `context`| Get log context by PkgId/PkgLogId |

Run `clscli <command> --help` for options. Main flags:

- **Global:** `--region` (required), `--output` / `-o` (json, csv, or file path).
- **query:** `-q` / `--query`, `-t` / `--topic` or `--topics`, `--last` or `--from` / `--to`, `--limit`, `--max` (cap total logs with auto-pagination).

Query syntax: CQL (recommended) or Lucene; optional SQL part after `|`. See [CLS SearchLog API](https://cloud.tencent.com/document/product/614/56447) and [SKILL.md](SKILL.md) for CQL/SQL details.

## Development

```bash
make build   # build binary
make run ARGS="--help"   # build and run with args
make tidy    # go mod tidy
make clean   # remove binary and dist/
```

## References

- [Tencent Cloud API common parameters (region, auth)](https://cloud.tencent.com/document/api/614/56474)
- [CLS SearchLog API](https://cloud.tencent.com/document/product/614/56447)
- [tencentcloud-sdk-go](https://github.com/TencentCloud/tencentcloud-sdk-go)

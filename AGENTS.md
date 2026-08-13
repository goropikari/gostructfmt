# AGENTS.md

## Testing

テスト方針の要点。
詳細なテストガイドラインは [TESTING.md](TESTING.md) を参照すること。

- `testify` を使う
- AAA パターンで書く
- `t.Run` を常に使う
- テーブル駆動テストは原則使わない
- テスト関数名は `Test{対象関数名}` を基本にする

## Development

- ユーザーから指摘を受けた場合、同様の問題の再発防止に役立つ内容は AGENTS.md またはその他の適切な規約ドキュメントにも追記する
- 全体テスト: `go test ./...`
- format: `make fmt`
- lint: `make lint`
- コードを編集したら `make fmt` と `make lint` でエラーが出ないことも確認する

### Change size guidance (800 lines)

Unless the change is mechanical the total number of changed lines should not exceed 800 lines. For complex logic changes the size should be under 500 lines.

If the change is larger, explore whether it can be split into reviewable stages and identify the smallest coherent stage to land first. Base the staging suggestion on the actual diff, dependencies, and affected call sites.

## Pull Requests

- PR を作るときは `.github/pull_request_template.md` に必ず従う
- PR 本文は GitHub Markdown として正しい記法で書く
- issue 由来の PR では、`Fixes #<issue number>` または `Closes #<issue number>` を PR 本文に含め、merge 時に GitHub が issue を自動 close するようにする

## graphify

This project has a knowledge graph at graphify-out/ with god nodes, community structure, and cross-file relationships.

When the user types `/graphify`, use the installed graphify skill or instructions before doing anything else.

Semantic graph extraction in this workspace is performed by Codex. Do not route the corpus to Ollama or another external backend.

Rules:

- For codebase questions, first run `graphify query "<question>"` when graphify-out/graph.json exists. Use `graphify path "<A>" "<B>"` for relationships and `graphify explain "<concept>"` for focused concepts. These return a scoped subgraph, usually much smaller than GRAPH_REPORT.md or raw grep output.
- Dirty graphify-out/ files are expected after hooks or incremental updates; dirty graph files are not a reason to skip graphify. Only skip graphify if the task is about stale or incorrect graph output, or the user explicitly says not to use it.
- If graphify-out/wiki/index.md exists, use it for broad navigation instead of raw source browsing.
- Read graphify-out/GRAPH_REPORT.md only for broad architecture review or when query/path/explain do not surface enough context.
- After modifying code, run `graphify update .` to keep the graph current (AST-only, no API cost).

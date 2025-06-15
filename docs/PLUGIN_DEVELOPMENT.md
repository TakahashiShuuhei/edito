# Edito プラグイン開発ガイド

## 概要

Editoでは、プラグインとユーザー設定ファイルの両方をGoで記述できます。
設定ファイルはGo言語のみをサポートし、タイプセーフで補完が効く環境で設定を記述できます。

## 1. アーキテクチャ

```
┌─────────────────┐    ┌─────────────────┐
│   edito本体     │    │  ユーザー設定    │
│   (バイナリ)    │    │  (config.go)    │
└─────────────────┘    └─────────────────┘
         │                       │
         │              ┌─────────────────┐
         │              │ edito-config    │
         │              │ (コンパイラ)    │
         │              └─────────────────┘
         │                       │
         ▼                       ▼
┌─────────────────┐    ┌─────────────────┐
│  API Layer      │◄───┤  config.so      │
│                 │    │  (コンパイル済み) │
└─────────────────┘    └─────────────────┘
         ▲
         │
┌─────────────────┐
│   プラグイン    │
│   (plugin.so)   │
└─────────────────┘
```

## 2. ユーザー設定ファイル

### 2.1 設定ファイルの作成

`~/.config/edito/config.go`:
```go
package config

import "github.com/TakahashiShuuhei/edito/pkg/edito"

func init() {
    // 基本設定
    edito.SetOption("tab-width", 4)
    edito.SetOption("theme", "dark")
    
    // キーバインド
    edito.BindKey("C-x C-s", "save-buffer")
    edito.BindKey("C-c C-c", "comment-region")
    
    // プラグイン読み込み
    edito.LoadPlugin("go-mode")
    edito.LoadPlugin("syntax-highlighting")
    
    // フック登録
    edito.RegisterHook("file-opened", func() {
        buf := edito.GetCurrentBuffer()
        if buf != nil && strings.HasSuffix(buf.GetFilename(), ".go") {
            edito.ShowMessage("Go file opened!")
        }
    })
}
```

### 2.2 設定ファイルのコンパイル

```bash
# edito-configツールを使用
edito-config ~/.config/edito/config.go

# または手動で
cd ~/.config/edito
go build -buildmode=plugin -o config.so config.go
```

### 2.3 設定ファイルの読み込み

Editoは起動時に以下の順序で設定を読み込みます：
1. `~/.config/edito/config.so` (コンパイル済みGo設定)
2. `~/.config/edito/config.go` (Go設定ソースファイル - fallback)

## 3. プラグイン開発

### 3.1 プラグインの構造

```go
package main

import (
    "github.com/TakahashiShuuhei/edito/internal/api"
    "github.com/TakahashiShuuhei/edito/pkg/edito"
)

type GoModePlugin struct {
    api *api.EditorAPI
}

func (p *GoModePlugin) Name() string {
    return "go-mode"
}

func (p *GoModePlugin) Version() string {
    return "1.0.0"
}

func (p *GoModePlugin) Init(api *api.EditorAPI) error {
    p.api = api
    
    // コマンド登録
    p.registerCommands()
    
    // キーバインド登録
    p.registerKeyBindings()
    
    return nil
}

func (p *GoModePlugin) Cleanup() error {
    return nil
}

func (p *GoModePlugin) registerCommands() {
    // Go固有のコマンドを登録
}

func (p *GoModePlugin) registerKeyBindings() {
    // Go固有のキーバインドを登録
}

// プラグインのエクスポート
var Plugin GoModePlugin
```

### 3.2 プラグインのビルド

```bash
# プラグインをビルド
go build -buildmode=plugin -o go-mode.so go-mode.go

# プラグインディレクトリに配置
mv go-mode.so ~/.local/share/edito/plugins/
```

## 4. APIリファレンス

### 4.1 基本操作

```go
// 設定の変更
edito.SetOption("key", value)

// キーバインドの設定
edito.BindKey("C-c C-g", "goto-definition")

// コマンドの実行
edito.ExecuteCommand("save-buffer", []string{})

// メッセージ表示
edito.ShowMessage("Hello from plugin!")
```

### 4.2 バッファ操作

```go
// 現在のバッファを取得
buf := edito.GetCurrentBuffer()
if buf != nil {
    // ファイル名を取得
    filename := buf.GetFilename()
    
    // カーソル位置を取得
    x, y := buf.GetCursorPosition()
    
    // テキストを取得
    lines := buf.GetLines()
    
    // テキストを設定
    buf.SetLines([]string{"new", "content"})
}
```

### 4.3 イベントフック

```go
// ファイルが開かれた時
edito.RegisterHook("file-opened", func() {
    // 処理
})

// ファイルが保存される前
edito.RegisterHook("before-save", func() {
    // 処理
})

// ファイルが保存された後
edito.RegisterHook("after-save", func() {
    // 処理
})
```

## 5. 開発フロー

### 5.1 設定ファイル開発

1. `~/.config/edito/config.go` を作成
2. `edito-config` でコンパイル
3. Editoを再起動して確認

### 5.2 プラグイン開発

1. プラグインソースを作成
2. `go build -buildmode=plugin` でビルド
3. プラグインディレクトリに配置
4. 設定ファイルで `edito.LoadPlugin()` を呼び出し

## 6. トラブルシューティング

### よくある問題

**Q: 設定ファイルが読み込まれない**
A: `config.so` が存在し、正しくコンパイルされているか確認してください。

**Q: プラグインが認識されない**
A: プラグインファイルが正しいディレクトリにあり、`var Plugin` が正しくエクスポートされているか確認してください。

**Q: APIが呼び出せない**
A: `github.com/TakahashiShuuhei/edito/pkg/edito` パッケージを正しくインポートしているか確認してください。
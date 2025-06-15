# Edito - Emacs風CLIエディタ

Emacsライクなキーバインディング、バッファシステム、プラグイン拡張機能を持つGoで実装されたCLIテキストエディタです。

## 特徴

- **Emacsライクなキーバインディング**: 馴染みのあるCtrl+A, Ctrl+E, Ctrl+P, Ctrl+Nなどのキーバインド
- **バッファシステム**: 複数ファイルの同時編集と切り替え
- **M-x コマンド実行**: Alt+Xでコマンドパレットを起動
- **Go設定ファイル**: Go言語でタイプセーフな設定ファイルを記述可能
- **XDG Base Directory準拠**: 標準的な設定ディレクトリ構造
- **プラグインシステム**: Goで書かれたプラグインによる機能拡張
- **パッケージマネージャ**: プラグインの簡単なインストール・管理
- **軽量**: 高速な起動とレスポンス

## インストール

```bash
go build -o edito
```

## 使用方法

```bash
./edito <filename>
```

## キーバインディング

| キー | 機能 |
|------|------|
| Ctrl+Q | 終了 |
| Ctrl+S | 保存 |
| Ctrl+A | 行頭に移動 |
| Ctrl+E | 行末に移動 |
| Ctrl+P | 上の行に移動 |
| Ctrl+N | 下の行に移動 |
| Ctrl+F | 右に移動 |
| Ctrl+B | 左に移動 |
| Alt+X | コマンドパレット起動 |
| Enter | 改行 |
| Backspace | 文字削除 |

## 利用可能なコマンド (M-x)

| コマンド | 機能 |
|----------|------|
| save-buffer | 現在のバッファを保存 |
| kill-buffer | 現在のバッファを閉じる |
| switch-to-buffer | 別のバッファに切り替え |
| find-file | ファイルを開く |
| list-buffers | 開いているバッファ一覧 |
| goto-line | 指定行に移動 |
| quit | エディタを終了 |

## プラグイン開発

プラグインはGoのpluginパッケージを使用してロードされる共有ライブラリ(.so)ファイルです。

### プラグインの実装例

```go
package main

import "edito/internal/plugin"

type MyPlugin struct{}

func (p *MyPlugin) Name() string {
    return "my-plugin"
}

func (p *MyPlugin) Version() string {
    return "1.0.0"
}

func (p *MyPlugin) Init(api *plugin.API) error {
    // プラグイン初期化処理
    return nil
}

func (p *MyPlugin) Execute(command string, args []string) error {
    // コマンド実行処理
    return nil
}

var Plugin MyPlugin
```

### プラグインのビルド

```bash
go build -buildmode=plugin -o myplugin.so myplugin.go
```

## 設定ファイル

設定ファイルは `~/.config/edito/config.go` にGo言語で記述します。

```go
package config

import (
    "strings"
    "edito/pkg/edito"
)

func init() {
    // 基本設定
    edito.SetOption("tab-width", 4)
    edito.SetOption("show-line-numbers", true)
    edito.SetOption("theme", "dark")
    
    // キーバインド設定
    edito.BindKey("C-x C-s", "save-buffer")
    edito.BindKey("C-x C-c", "quit")
    edito.BindKey("M-g g", "goto-line")
    
    // プラグイン読み込み
    edito.LoadPlugin("syntax-highlighting")
    edito.LoadPlugin("go-mode")
    
    // イベントフック - ファイル種別に応じた設定
    edito.RegisterHook("file-opened", func() {
        buf := edito.GetCurrentBuffer()
        if buf != nil {
            filename := buf.GetFilename()
            if strings.HasSuffix(filename, ".go") {
                edito.SetOption("use-tabs", true)
                edito.ShowMessage("Go mode activated")
            }
        }
    })
}
```

設定ファイルをコンパイルして使用します：

```bash
edito-config ~/.config/edito/config.go
```

## ディレクトリ構造

```
edito/
├── main.go                         # エントリーポイント
├── internal/
│   ├── buffer/                     # バッファ管理
│   │   └── buffer.go
│   ├── command/                    # コマンドシステム
│   │   └── command.go
│   ├── config/                     # 設定管理
│   │   ├── config.go
│   │   └── editorc.go
│   ├── editor/                     # エディタコア機能
│   │   ├── editor.go
│   │   ├── commands.go
│   │   └── editor_test.go
│   ├── keybinding/                 # キーバインディングシステム
│   │   └── keybinding.go
│   ├── minibuffer/                 # コマンドパレット
│   │   └── minibuffer.go
│   ├── plugin/                     # プラグインシステム
│   │   └── plugin.go
│   └── package_manager/            # パッケージマネージャ
│       └── manager.go
├── example-config/                 # 設定例
│   └── .config/edito/
│       └── config.go               # Go設定ファイル例
├── example-plugin/                 # プラグイン例
│   └── syntax-highlighting.go
├── go.mod
└── README.md
```

## XDG Base Directory 準拠

- 設定ファイル: `$XDG_CONFIG_HOME/edito/config.go` (デフォルト: `~/.config/edito/config.go`)
- コンパイル済み設定: `$XDG_CONFIG_HOME/edito/config.so`
- データファイル: `$XDG_DATA_HOME/edito/` (デフォルト: `~/.local/share/edito/`)
- キャッシュファイル: `$XDG_CACHE_HOME/edito/` (デフォルト: `~/.cache/edito/`)
- プラグイン: `$XDG_DATA_HOME/edito/plugins/`

## ライセンス

MIT License
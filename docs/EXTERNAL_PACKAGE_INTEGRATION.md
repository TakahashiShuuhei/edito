# 外部パッケージ統合ガイド

editoに外部パッケージを組み込む方法を説明します。

## 方法1: go.mod依存関係として追加（推奨）

### 1. go.modに追加
```go
require (
    github.com/nsf/termbox-go v1.1.1
    github.com/TakahashiShuuhei/edito-file-tree v0.1.0
)
```

### 2. 統合コードを作成
```go
// internal/editor/file_tree_integration.go
package editor

import (
    filetree "github.com/TakahashiShuuhei/edito-file-tree"
)

func (e *Editor) setupFileTree() {
    tree := filetree.New()
    
    // キーバインド追加
    e.keyMap.BindKey(termbox.KeyCtrlT, func() {
        e.toggleFileTree()
    })
    
    // コマンド追加
    e.commandRegistry.Register("toggle-file-tree", "Toggle file tree", func(args []string) error {
        return e.toggleFileTree()
    })
}
```

### 3. エディタ初期化時に呼び出し
```go
// internal/editor/editor.go
func New() *Editor {
    // ...
    e.setupFileTree()  // 追加
    return e
}
```

## 方法2: プラグインとして動的ロード

### 1. edito-file-treeをプラグイン化
```go
// edito-file-tree側でプラグインインターフェースを実装
package main

import "github.com/TakahashiShuuhei/edito/internal/plugin"

type FileTreePlugin struct {
    api *plugin.API
}

func (p *FileTreePlugin) Name() string { return "file-tree" }
func (p *FileTreePlugin) Version() string { return "0.1.0" }

func (p *FileTreePlugin) Init(api *plugin.API) error {
    p.api = api
    // ファイルツリー機能を登録
    return nil
}

var Plugin FileTreePlugin
```

### 2. .soファイルとしてビルド
```bash
cd edito-file-tree
go build -buildmode=plugin -o file-tree.so
cp file-tree.so ~/.local/share/edito/plugins/
```

### 3. 設定ファイルでロード
```go
// ~/.config/edito/config.go
func init() {
    edito.LoadPlugin("file-tree")
}
```

## 方法3: 設定ファイルでimport

### 1. 設定ファイルに直接import
```go
// ~/.config/edito/config.go
package config

import (
    "github.com/TakahashiShuuhei/edito/pkg/edito"
    filetree "github.com/TakahashiShuuhei/edito-file-tree"
)

func init() {
    // ファイルツリーを初期化
    tree := filetree.New()
    
    // カスタムキーバインド
    edito.RegisterHook("editor-ready", func() {
        // ファイルツリーを統合
    })
}
```

### 2. 設定をコンパイル
```bash
cd ~/.config/edito
go mod init config
go get github.com/TakahashiShuuhei/edito-file-tree
edito-config config.go
```

## 推奨アプローチ

**edito-file-tree**のような機能なら、**方法1**（go.mod依存関係）が最適です：

1. **型安全**: コンパイル時に依存関係をチェック
2. **パフォーマンス**: 実行時ロードのオーバーヘッドなし
3. **簡単**: 通常のGoライブラリと同じ扱い

## 実装例

```go
// edito-file-tree パッケージがこのインターフェースを提供する想定
type FileTree interface {
    Show() error
    Hide() error
    Toggle() error
    GetSelectedPath() string
    SetRootPath(path string) error
}

// edito本体での統合
func (e *Editor) integrateFileTree(tree FileTree) {
    e.keyMap.BindKey(termbox.KeyCtrlT, func() {
        tree.Toggle()
    })
    
    e.commandRegistry.Register("file-tree-toggle", "Toggle file tree", func(args []string) error {
        return tree.Toggle()
    })
}
```
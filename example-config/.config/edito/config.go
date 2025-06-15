package config

import (
	"strings"
	"github.com/TakahashiShuuhei/edito/pkg/edito"
)

func init() {
	// 基本設定
	edito.SetOption("tab-width", 4)
	edito.SetOption("show-line-numbers", true)
	edito.SetOption("auto-save", true)
	edito.SetOption("theme", "dark")
	
	// キーバインド設定
	edito.BindKey("C-x C-s", "save-buffer")
	edito.BindKey("C-x C-c", "quit")
	edito.BindKey("C-x b", "switch-to-buffer")
	edito.BindKey("C-x C-f", "find-file")
	edito.BindKey("C-x k", "kill-buffer")
	edito.BindKey("M-g g", "goto-line")
	edito.BindKey("C-c C-c", "comment-region")
	
	// プラグイン読み込み
	edito.LoadPlugin("syntax-highlighting")
	edito.LoadPlugin("auto-complete")
	edito.LoadPlugin("git-integration")
	edito.LoadPlugin("file-tree")
	
	// プラグイン固有設定
	edito.SetOption("syntax-highlighting-theme", "monokai")
	edito.SetOption("auto-complete-delay", 500)
	edito.SetOption("git-show-diff", true)
	
	// イベントフック登録
	edito.RegisterHook("file-opened", func() {
		buf := edito.GetCurrentBuffer()
		if buf != nil {
			filename := buf.GetFilename()
			if strings.HasSuffix(filename, ".go") {
				edito.ShowMessage("Go file opened: " + filename)
				// Go特有の設定を適用
				edito.SetOption("tab-width", 4)
				edito.SetOption("use-tabs", true)
			} else if strings.HasSuffix(filename, ".py") {
				edito.ShowMessage("Python file opened: " + filename)
				// Python特有の設定を適用
				edito.SetOption("tab-width", 4)
				edito.SetOption("use-tabs", false)
			}
		}
	})
	
	edito.RegisterHook("before-save", func() {
		buf := edito.GetCurrentBuffer()
		if buf != nil && strings.HasSuffix(buf.GetFilename(), ".go") {
			// Goファイル保存前にフォーマット実行
			edito.ExecuteCommand("format-buffer", []string{})
		}
	})
	
	edito.RegisterHook("after-save", func() {
		buf := edito.GetCurrentBuffer()
		if buf != nil {
			edito.ShowMessage("File saved: " + buf.GetFilename())
		}
	})
}
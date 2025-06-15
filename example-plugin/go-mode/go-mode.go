package main

import (
	"fmt"
	"os/exec"
	"strings"
	
	"github.com/TakahashiShuuhei/edito/internal/api"
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

func (p *GoModePlugin) Init(editorAPI *api.EditorAPI) error {
	p.api = editorAPI
	
	// Go固有のコマンドを登録
	p.registerCommands()
	
	return nil
}

func (p *GoModePlugin) Cleanup() error {
	return nil
}

func (p *GoModePlugin) registerCommands() {
	// Go固有のコマンドをエディタに登録する処理
	// 実際の実装では editorAPI を通じてコマンドを登録
}

func (p *GoModePlugin) formatBuffer() error {
	buf := p.api.GetCurrentBuffer()
	if buf == nil {
		return fmt.Errorf("no current buffer")
	}
	
	// gofmtを実行
	lines := buf.GetLines()
	code := strings.Join(lines, "\n")
	
	cmd := exec.Command("gofmt")
	cmd.Stdin = strings.NewReader(code)
	
	output, err := cmd.Output()
	if err != nil {
		p.api.ShowMessage("gofmt failed: " + err.Error())
		return err
	}
	
	// フォーマット済みコードをバッファに設定
	formattedLines := strings.Split(string(output), "\n")
	buf.SetLines(formattedLines)
	
	p.api.ShowMessage("Buffer formatted with gofmt")
	return nil
}

func (p *GoModePlugin) runGoTest() error {
	buf := p.api.GetCurrentBuffer()
	if buf == nil {
		return fmt.Errorf("no current buffer")
	}
	
	filename := buf.GetFilename()
	if !strings.HasSuffix(filename, ".go") {
		p.api.ShowMessage("Not a Go file")
		return fmt.Errorf("not a go file")
	}
	
	// go testを実行
	cmd := exec.Command("go", "test", "-v")
	output, err := cmd.CombinedOutput()
	
	if err != nil {
		p.api.ShowMessage("Tests failed: " + string(output))
	} else {
		p.api.ShowMessage("Tests passed!")
	}
	
	return nil
}

func (p *GoModePlugin) addImport(packageName string) error {
	buf := p.api.GetCurrentBuffer()
	if buf == nil {
		return fmt.Errorf("no current buffer")
	}
	
	lines := buf.GetLines()
	
	// import文を見つけて追加
	for i, line := range lines {
		if strings.Contains(line, "import (") {
			// 複数import文のブロック内に追加
			newLines := make([]string, len(lines)+1)
			copy(newLines[:i+1], lines[:i+1])
			newLines[i+1] = "\t\"" + packageName + "\""
			copy(newLines[i+2:], lines[i+1:])
			buf.SetLines(newLines)
			return nil
		}
	}
	
	// 単一import文の後に追加
	for i, line := range lines {
		if strings.HasPrefix(strings.TrimSpace(line), "import ") {
			newLines := make([]string, len(lines)+1)
			copy(newLines[:i+1], lines[:i+1])
			newLines[i+1] = "import \"" + packageName + "\""
			copy(newLines[i+2:], lines[i+1:])
			buf.SetLines(newLines)
			return nil
		}
	}
	
	p.api.ShowMessage("Could not find import section")
	return fmt.Errorf("import section not found")
}

// プラグインのエクスポート - これが重要！
var Plugin GoModePlugin
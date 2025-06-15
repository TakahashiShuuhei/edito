// Example config.go showing automatic plugin installation
package config

import "github.com/TakahashiShuuhei/edito/pkg/edito"

func init() {
	// Install the file-tree plugin automatically from the user's repository
	edito.InstallPlugin("file-tree", "github.com/TakahashiShuuhei/edito-file-tree", "v0.1.0")
	
	// Load the plugin after installation
	edito.LoadPlugin("file-tree")
	
	// Bind a key to toggle the file tree
	edito.BindKey("C-x C-t", "file-tree-toggle")
	
	// Set some editor options
	edito.SetOption("tab-width", 4)
	edito.SetOption("auto-indent", true)
}
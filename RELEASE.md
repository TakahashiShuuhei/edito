# Release Process

## Creating a New Release

1. **Update version in code**
   ```bash
   # Update version.go
   sed -i 's/Version   = "[0-9.]*"/Version   = "0.2.0"/' version.go
   
   # Update Makefile
   sed -i 's/edito version [0-9.]*/edito version 0.2.0/' Makefile
   ```

2. **Commit version bump**
   ```bash
   git add version.go Makefile
   git commit -m "Bump version to 0.2.0"
   git push
   ```

3. **Create and push tag**
   ```bash
   git tag v0.2.0
   git push origin v0.2.0
   ```

4. **Create GitHub Release**
   - Go to https://github.com/TakahashiShuuhei/edito/releases
   - Click "Create a new release"
   - Choose tag: v0.2.0
   - Release title: "v0.2.0"
   - Describe changes and new features
   - Click "Publish release"

5. **Test installation**
   ```bash
   go install github.com/TakahashiShuuhei/edito@v0.2.0
   go install github.com/TakahashiShuuhei/edito/cmd/edito-config@v0.2.0
   ```

## Current Release: v0.1.0

Initial release featuring:
- Emacs-like CLI text editor
- Go-based configuration system
- Plugin architecture
- Buffer management
- M-x command palette
- XDG Base Directory compliance
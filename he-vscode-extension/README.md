# HE VS Code Language Support

Minimal VS Code language extension for HE (.he) files.

## Install locally

1. Install Node.js and npm (if you don't already have them).
2. Install `vsce` packaging tool: `npm install -g vsce`.
3. From the extension root (folder containing package.json), run: `vsce package`.
4. Install the produced `he-language-0.0.1.vsix` file in VS Code: `code --install-extension he-language-0.0.1.vsix`

## Developer test

Open the folder in VSCode and press `F5` to launch an Extension Development Host with the extension loaded.
```

---

## Packaging & Installing (quick commands)

```bash
# (1) ensure Node.js/npm installed
npm --version

# (2) install vsce if missing
npm install -g vsce

# (3) package the extension (run in extension root)
vsce package
# produces he-language-0.0.1.vsix

# (4) install the .vsix into your local VS Code
code --install-extension he-language-0.0.1.vsix

# (5) test in dev mode
# open the extension folder in VS Code and press F5
const vscode = require("vscode");
const path = require("path");

function activate(context) {
    const runCmd = vscode.commands.registerCommand("he.runFile", () => {
        const editor = vscode.window.activeTextEditor;
        if (!editor) return;

        const file = editor.document.fileName;
        const folder = vscode.workspace.workspaceFolders?.[0]?.uri.fsPath;
        const exe = path.join(folder, "heinterp.exe");

        const terminal = vscode.window.createTerminal("HE Runner");
        terminal.sendText(`"${exe}" "${file}"`);
        terminal.show();
    });

    context.subscriptions.push(runCmd);
}

function deactivate() {}

module.exports = { activate, deactivate };

import * as vscode from 'vscode';
import { exec } from 'child_process';

function getActiveFilePath(): string | undefined {
  const editor = vscode.window.activeTextEditor;
  if (!editor) return undefined;
  const doc = editor.document;
  return doc.uri.fsPath;
}

function runHeFile(debug: boolean): void {
  const filePath = getActiveFilePath();
  if (!filePath) {
    vscode.window.showErrorMessage('No active editor / file to run.');
    return;
  }

  // HE CLI in this repo is implemented by Go at ./cmd/hunterlang
  // We run via `go run` to avoid needing a separate binary install.
  const cmd = `go run ./cmd/hunterlang -file "${filePath}"`;

  vscode.window.withProgress(
    {
      location: vscode.ProgressLocation.Notification,
      title: debug ? 'HE: Debug run' : 'HE: Run',
      cancellable: false,
    },
    async () => {
      const channel = vscode.window.createOutputChannel(debug ? 'HE Debug' : 'HE Run');
      channel.clear();
      channel.show(true);

      channel.appendLine(`$ ${cmd}`);
      exec(
        cmd,
        { cwd: vscode.workspace.rootPath || undefined },
        (error: Error | null, stdout: string, stderr: string) => {
          if (stdout) channel.appendLine(stdout);
          if (stderr) channel.appendLine(stderr);

          if (error) {
            vscode.window.showErrorMessage(`HE failed: ${error.message}`);
            return;
          }
          vscode.window.showInformationMessage(debug ? 'HE debug run completed' : 'HE run completed');
        }
      );
    }
  );
}

export function activate(context: vscode.ExtensionContext) {
  const runCmd = vscode.commands.registerCommand('hunterlang.runHe', () => runHeFile(false));
  const debugCmd = vscode.commands.registerCommand('hunterlang.debugHe', () => runHeFile(true));

  context.subscriptions.push(runCmd, debugCmd);
}

export function deactivate() {}

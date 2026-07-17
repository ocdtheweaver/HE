// HE Language VS Code Extension
// Handles: run, check, build, repl commands + check-on-save

const vscode = require('vscode');
const path = require('path');
const fs = require('fs');

let outputChannel;
let diagnosticCollection;

function activate(context) {
  outputChannel = vscode.window.createOutputChannel('HE');
  diagnosticCollection = vscode.languages.createDiagnosticCollection('he');

  context.subscriptions.push(
    vscode.commands.registerCommand('he.runFile', cmdRunFile),
    vscode.commands.registerCommand('he.checkFile', cmdCheckFile),
    vscode.commands.registerCommand('he.buildFile', cmdBuildFile),
    vscode.commands.registerCommand('he.openRepl', cmdOpenRepl),
    vscode.commands.registerCommand('he.runGenesis', cmdRunGenesis),
  );

  // Check on save
  context.subscriptions.push(
    vscode.workspace.onDidSaveTextDocument(doc => {
      if (doc.languageId === 'he') {
        const cfg = vscode.workspace.getConfiguration('he');
        if (cfg.get('checkOnSave')) {
          runCheck(doc.fileName, true);
        }
      }
    })
  );

  context.subscriptions.push(diagnosticCollection);
}

function deactivate() {
  if (diagnosticCollection) {
    diagnosticCollection.clear();
    diagnosticCollection.dispose();
  }
}

// ── Helpers ────────────────────────────────────────────────────────────────────

function heExe() {
  return vscode.workspace.getConfiguration('he').get('executablePath') || 'he';
}

function activeFile() {
  const ed = vscode.window.activeTextEditor;
  if (!ed || ed.document.languageId !== 'he') {
    vscode.window.showWarningMessage('Open a .he file first.');
    return null;
  }
  return ed.document.fileName;
}

function runInTerminal(cmd, cwd) {
  let terminal = vscode.window.terminals.find(t => t.name === 'HE');
  if (!terminal) {
    terminal = vscode.window.createTerminal({ name: 'HE', cwd });
  } else {
    terminal.sendText(`cd "${cwd || path.dirname(activeFile() || '.')}"`);
  }
  terminal.show(true);
  terminal.sendText(cmd);
}

// ── Commands ───────────────────────────────────────────────────────────────────

function cmdRunFile() {
  const file = activeFile();
  if (!file) return;
  const dir = path.dirname(file);
  const name = path.basename(file);
  runInTerminal(`${heExe()} run "${name}"`, dir);
}

function cmdCheckFile() {
  const file = activeFile();
  if (!file) return;
  runCheck(file, false);
}

async function runCheck(filePath, silent) {
  const { exec } = require('child_process');
  const dir = path.dirname(filePath);
  const name = path.basename(filePath);
  const cmd = `${heExe()} check "${name}"`;

  exec(cmd, { cwd: dir }, (err, stdout, stderr) => {
    const output = stdout + stderr;

    // Clear old diagnostics for this file
    const uri = vscode.Uri.file(filePath);
    diagnosticCollection.set(uri, []);

    if (!silent) {
      outputChannel.clear();
      outputChannel.appendLine(`he check: ${name}`);
      outputChannel.appendLine('─'.repeat(50));
      outputChannel.appendLine(output);
      outputChannel.show(true);
    }

    // Parse any "Type errors" lines into VS Code diagnostics
    const diagnostics = [];
    const typeErrPattern = /line (\d+):\s+(.+)/g;
    let m;
    while ((m = typeErrPattern.exec(output)) !== null) {
      const line = parseInt(m[1], 10) - 1;
      const msg = m[2].trim();
      const range = new vscode.Range(line, 0, line, 200);
      diagnostics.push(new vscode.Diagnostic(
        range,
        msg,
        vscode.DiagnosticSeverity.Warning
      ));
    }

    if (diagnostics.length > 0) {
      diagnosticCollection.set(uri, diagnostics);
      if (silent) {
        // Show a subtle status bar hint instead of popping the output panel
        vscode.window.setStatusBarMessage(
          `HE: ${diagnostics.length} issue(s) in ${name}`,
          5000
        );
      }
    } else if (!silent) {
      vscode.window.setStatusBarMessage('HE: no issues found ✓', 3000);
    }
  });
}

function cmdBuildFile() {
  const file = activeFile();
  if (!file) return;

  const dir = path.dirname(file);
  const base = path.basename(file, '.he');
  const cfg = vscode.workspace.getConfiguration('he');
  const outDir = cfg.get('buildOutputDir') || '.';
  const outPath = path.join(outDir, `${base}.hex`);

  runInTerminal(`${heExe()} build "${path.basename(file)}" --o "${outPath}"`, dir);
}

function cmdOpenRepl() {
  const file = vscode.window.activeTextEditor?.document?.fileName;
  const cwd = file ? path.dirname(file) : undefined;
  runInTerminal(`${heExe()} repl`, cwd);
}

function cmdRunGenesis() {
  // Find genesis.he in the workspace root or current file's folder
  const workspaceFolders = vscode.workspace.workspaceFolders;
  let cwd;

  if (workspaceFolders && workspaceFolders.length > 0) {
    cwd = workspaceFolders[0].uri.fsPath;
    // Check if genesis.he exists in workspace root
    if (!fs.existsSync(path.join(cwd, 'genesis.he'))) {
      // Fall back to active file's folder
      const file = vscode.window.activeTextEditor?.document?.fileName;
      if (file) cwd = path.dirname(file);
    }
  } else {
    const file = vscode.window.activeTextEditor?.document?.fileName;
    if (file) cwd = path.dirname(file);
  }

  if (!cwd) {
    vscode.window.showWarningMessage('Open a workspace or a .he file first.');
    return;
  }

  runInTerminal(`${heExe()} run`, cwd);
}

module.exports = { activate, deactivate };

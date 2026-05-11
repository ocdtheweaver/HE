"use strict";
var __createBinding = (this && this.__createBinding) || (Object.create ? (function(o, m, k, k2) {
    if (k2 === undefined) k2 = k;
    var desc = Object.getOwnPropertyDescriptor(m, k);
    if (!desc || ("get" in desc ? !m.__esModule : desc.writable || desc.configurable)) {
      desc = { enumerable: true, get: function() { return m[k]; } };
    }
    Object.defineProperty(o, k2, desc);
}) : (function(o, m, k, k2) {
    if (k2 === undefined) k2 = k;
    o[k2] = m[k];
}));
var __setModuleDefault = (this && this.__setModuleDefault) || (Object.create ? (function(o, v) {
    Object.defineProperty(o, "default", { enumerable: true, value: v });
}) : function(o, v) {
    o["default"] = v;
});
var __importStar = (this && this.__importStar) || (function () {
    var ownKeys = function(o) {
        ownKeys = Object.getOwnPropertyNames || function (o) {
            var ar = [];
            for (var k in o) if (Object.prototype.hasOwnProperty.call(o, k)) ar[ar.length] = k;
            return ar;
        };
        return ownKeys(o);
    };
    return function (mod) {
        if (mod && mod.__esModule) return mod;
        var result = {};
        if (mod != null) for (var k = ownKeys(mod), i = 0; i < k.length; i++) if (k[i] !== "default") __createBinding(result, mod, k[i]);
        __setModuleDefault(result, mod);
        return result;
    };
})();
Object.defineProperty(exports, "__esModule", { value: true });
exports.activate = activate;
exports.deactivate = deactivate;
const vscode = __importStar(require("vscode"));
const child_process_1 = require("child_process");
function getActiveFilePath() {
    const editor = vscode.window.activeTextEditor;
    if (!editor)
        return undefined;
    const doc = editor.document;
    return doc.uri.fsPath;
}
function runHeFile(debug) {
    const filePath = getActiveFilePath();
    if (!filePath) {
        vscode.window.showErrorMessage('No active editor / file to run.');
        return;
    }
    // HE CLI in this repo is implemented by Go at ./cmd/hunterlang
    // We run via `go run` to avoid needing a separate binary install.
    const cmd = `go run ./cmd/hunterlang -file "${filePath}"`;
    vscode.window.withProgress({
        location: vscode.ProgressLocation.Notification,
        title: debug ? 'HE: Debug run' : 'HE: Run',
        cancellable: false,
    }, async () => {
        const channel = vscode.window.createOutputChannel(debug ? 'HE Debug' : 'HE Run');
        channel.clear();
        channel.show(true);
        channel.appendLine(`$ ${cmd}`);
        (0, child_process_1.exec)(cmd, { cwd: vscode.workspace.rootPath || undefined }, (error, stdout, stderr) => {
            if (stdout)
                channel.appendLine(stdout);
            if (stderr)
                channel.appendLine(stderr);
            if (error) {
                vscode.window.showErrorMessage(`HE failed: ${error.message}`);
                return;
            }
            vscode.window.showInformationMessage(debug ? 'HE debug run completed' : 'HE run completed');
        });
    });
}
function activate(context) {
    const runCmd = vscode.commands.registerCommand('hunterlang.runHe', () => runHeFile(false));
    const debugCmd = vscode.commands.registerCommand('hunterlang.debugHe', () => runHeFile(true));
    context.subscriptions.push(runCmd, debugCmd);
}
function deactivate() { }

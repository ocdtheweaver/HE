# HE Language — Windows Uninstaller
# Run as Administrator

param(
    [string]$InstallDir = "C:\Program Files\HE"
)

$isAdmin = ([Security.Principal.WindowsPrincipal][Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)
if (-not $isAdmin) {
    Write-Host "Please run as Administrator." -ForegroundColor Red
    pause; exit 1
}

Write-Host ""
Write-Host "  HE Language — Uninstaller" -ForegroundColor DarkYellow
Write-Host ""

# Remove registry entries
Remove-Item -Path "HKCU:\Software\Classes\.he" -Recurse -Force -ErrorAction SilentlyContinue
Remove-Item -Path "HKCU:\Software\Classes\HESourceFile" -Recurse -Force -ErrorAction SilentlyContinue
Write-Host "  ✓ Registry entries removed" -ForegroundColor Green

# Remove from PATH
$currentPath = [Environment]::GetEnvironmentVariable("PATH", "Machine")
$newPath = ($currentPath.Split(';') | Where-Object { $_ -ne $InstallDir }) -join ';'
[Environment]::SetEnvironmentVariable("PATH", $newPath, "Machine")
Write-Host "  ✓ Removed from PATH" -ForegroundColor Green

# Remove files
if (Test-Path $InstallDir) {
    Remove-Item -Path $InstallDir -Recurse -Force -ErrorAction SilentlyContinue
    Write-Host "  ✓ Removed $InstallDir" -ForegroundColor Green
}

# Refresh shell
try {
    $code = @"
using System;
using System.Runtime.InteropServices;
public class Shell32 {
    [DllImport("shell32.dll")] public static extern void SHChangeNotify(int wEventId, int uFlags, IntPtr dwItem1, IntPtr dwItem2);
}
"@
    Add-Type -TypeDefinition $code -ErrorAction SilentlyContinue
    [Shell32]::SHChangeNotify(0x08000000, 0, [IntPtr]::Zero, [IntPtr]::Zero)
} catch {}

Write-Host ""
Write-Host "  HE uninstalled." -ForegroundColor DarkYellow
Write-Host ""
pause

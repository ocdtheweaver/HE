# HE Language — Windows Installer
# Run this as Administrator in PowerShell:
#   Right-click install.ps1 → Run as Administrator
# Or from an admin PowerShell prompt:
#   Set-ExecutionPolicy Bypass -Scope Process -Force; .\install.ps1

param(
    [string]$InstallDir = "C:\Program Files\HE"
)

$ErrorActionPreference = "Stop"

Write-Host ""
Write-Host "  HE Language v5.0.0 — Windows Installer" -ForegroundColor DarkYellow
Write-Host "  ─────────────────────────────────────────" -ForegroundColor DarkGray
Write-Host ""

# ── Check admin ───────────────────────────────────────────────────────────────

$isAdmin = ([Security.Principal.WindowsPrincipal][Security.Principal.WindowsIdentity]::GetCurrent()).IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)
if (-not $isAdmin) {
    Write-Host "  [ERROR] Please run as Administrator." -ForegroundColor Red
    Write-Host "  Right-click install.ps1 and choose 'Run as Administrator'."
    pause
    exit 1
}

# ── Locate files ──────────────────────────────────────────────────────────────

$scriptDir = Split-Path -Parent $MyInvocation.MyCommand.Path
$exeSource  = Join-Path $scriptDir "he.exe"
$icoSource  = Join-Path $scriptDir "hefile.ico"

# he.exe might be one folder up (in the HE root)
if (-not (Test-Path $exeSource)) {
    $exeSource = Join-Path (Split-Path -Parent $scriptDir) "he.exe"
}
if (-not (Test-Path $icoSource)) {
    $icoSource = Join-Path (Split-Path -Parent $scriptDir) "windows\hefile.ico"
}
if (-not (Test-Path $icoSource)) {
    $icoSource = Join-Path $scriptDir "..\vscode-he\images\hefile.ico"
}

if (-not (Test-Path $exeSource)) {
    Write-Host "  [ERROR] he.exe not found near $scriptDir" -ForegroundColor Red
    Write-Host "  Place he.exe next to install.ps1 and try again."
    pause
    exit 1
}

if (-not (Test-Path $icoSource)) {
    Write-Host "  [WARNING] hefile.ico not found — file icon will use Windows default" -ForegroundColor Yellow
    $icoSource = $null
}

# ── Install files ─────────────────────────────────────────────────────────────

Write-Host "  Installing to: $InstallDir" -ForegroundColor Cyan
New-Item -ItemType Directory -Force -Path $InstallDir | Out-Null

$exeDest = Join-Path $InstallDir "he.exe"
$icoDest = Join-Path $InstallDir "hefile.ico"

Copy-Item $exeSource $exeDest -Force
Write-Host "  ✓ Copied he.exe" -ForegroundColor Green

if ($icoSource) {
    Copy-Item $icoSource $icoDest -Force
    Write-Host "  ✓ Copied hefile.ico" -ForegroundColor Green
}

# ── Add to PATH ───────────────────────────────────────────────────────────────

$currentPath = [Environment]::GetEnvironmentVariable("PATH", "Machine")
if ($currentPath -notlike "*$InstallDir*") {
    [Environment]::SetEnvironmentVariable("PATH", "$currentPath;$InstallDir", "Machine")
    Write-Host "  ✓ Added to system PATH" -ForegroundColor Green
} else {
    Write-Host "  ✓ Already on system PATH" -ForegroundColor Green
}

# ── Register .he file type ────────────────────────────────────────────────────

Write-Host ""
Write-Host "  Registering .he file type..." -ForegroundColor Cyan

# File extension
New-Item -Path "HKCU:\Software\Classes\.he" -Force | Out-Null
Set-ItemProperty -Path "HKCU:\Software\Classes\.he" -Name "(default)" -Value "HESourceFile"
Set-ItemProperty -Path "HKCU:\Software\Classes\.he" -Name "Content Type" -Value "text/x-he"
Set-ItemProperty -Path "HKCU:\Software\Classes\.he" -Name "PerceivedType" -Value "text"

# File type
New-Item -Path "HKCU:\Software\Classes\HESourceFile" -Force | Out-Null
Set-ItemProperty -Path "HKCU:\Software\Classes\HESourceFile" -Name "(default)" -Value "HE Source File"

# Icon
if ($icoSource) {
    New-Item -Path "HKCU:\Software\Classes\HESourceFile\DefaultIcon" -Force | Out-Null
    Set-ItemProperty -Path "HKCU:\Software\Classes\HESourceFile\DefaultIcon" -Name "(default)" -Value "`"$icoDest`",0"
    Write-Host "  ✓ Registered wolf icon for .he files" -ForegroundColor Green
}

# Open command: he run <file>
New-Item -Path "HKCU:\Software\Classes\HESourceFile\shell\open\command" -Force | Out-Null
Set-ItemProperty -Path "HKCU:\Software\Classes\HESourceFile\shell" -Name "(default)" -Value "open"
Set-ItemProperty -Path "HKCU:\Software\Classes\HESourceFile\shell\open" -Name "(default)" -Value "Run with HE"
Set-ItemProperty -Path "HKCU:\Software\Classes\HESourceFile\shell\open\command" -Name "(default)" -Value "`"$exeDest`" run `"%1`""

# Check command: he check <file> (opens a cmd window so you can read output)
New-Item -Path "HKCU:\Software\Classes\HESourceFile\shell\check\command" -Force | Out-Null
Set-ItemProperty -Path "HKCU:\Software\Classes\HESourceFile\shell\check" -Name "(default)" -Value "Check with HE"
Set-ItemProperty -Path "HKCU:\Software\Classes\HESourceFile\shell\check\command" -Name "(default)" -Value "cmd.exe /k `"$exeDest`" check `"%1`""

Write-Host "  ✓ Registered open + check right-click actions" -ForegroundColor Green

# ── Refresh icon cache ────────────────────────────────────────────────────────

Write-Host ""
Write-Host "  Refreshing icon cache..." -ForegroundColor Cyan

# Kill and restart Explorer to pick up the new icon immediately
try {
    # Touch the icon cache database
    $iconCache = "$env:LOCALAPPDATA\Microsoft\Windows\Explorer"
    if (Test-Path $iconCache) {
        Get-ChildItem "$iconCache\iconcache_*.db" | ForEach-Object {
            $_ | Remove-Item -Force -ErrorAction SilentlyContinue
        }
    }
    # Notify shell of the change
    $code = @"
using System;
using System.Runtime.InteropServices;
public class Shell {
    [DllImport("shell32.dll")] public static extern void SHChangeNotify(int wEventId, int uFlags, IntPtr dwItem1, IntPtr dwItem2);
}
"@
    Add-Type -TypeDefinition $code -ErrorAction SilentlyContinue
    [Shell]::SHChangeNotify(0x08000000, 0, [IntPtr]::Zero, [IntPtr]::Zero)
    Write-Host "  ✓ Icon cache refreshed" -ForegroundColor Green
} catch {
    Write-Host "  ✓ You may need to sign out and back in to see the new icon" -ForegroundColor Yellow
}

# ── Done ──────────────────────────────────────────────────────────────────────

Write-Host ""
Write-Host "  ─────────────────────────────────────────" -ForegroundColor DarkGray
Write-Host "  HE v5.0.0 installed successfully!" -ForegroundColor DarkYellow
Write-Host ""
Write-Host "  Try it:" -ForegroundColor White
Write-Host "    he run myapp.he" -ForegroundColor Gray
Write-Host "    he version" -ForegroundColor Gray
Write-Host ""
Write-Host "  .he files will now show the wolf icon in File Explorer." -ForegroundColor White
Write-Host "  Right-click any .he file to Run or Check with HE." -ForegroundColor White
Write-Host ""

pause

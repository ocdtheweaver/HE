@echo off
:: Silent Hunter's Engine (.he) Setup

:: 1. Create target folder for icon
set ICON_DIR=C:\Program Files\HELanguage
if not exist "%ICON_DIR%" mkdir "%ICON_DIR%"

:: 2. Copy the icon (overwrite if exists)
copy /Y "hefile.png" "%ICON_DIR%\hefile.png" >nul

:: 3. Create temporary .reg file
set REGFILE=%TEMP%\hefile-association.reg
(
echo Windows Registry Editor Version 5.00
echo.
echo [HKEY_CURRENT_USER\Software\Classes\.he]
echo @="HEFile"
echo.
echo [HKEY_CURRENT_USER\Software\Classes\HEFile]
echo @="Hunter's Engine Source File"
echo.
echo [HKEY_CURRENT_USER\Software\Classes\HEFile\DefaultIcon]
echo @="%ICON_DIR%\\hefile.png"
) > "%REGFILE%"

:: 4. Apply registry
reg import "%REGFILE%" >nul

:: 5. Create a sample .he file for testing if it doesn't exist
set SAMPLE_FILE=%CD%\Genesis.he
if not exist "%SAMPLE_FILE%" echo // Sample Hunter's Engine file > "%SAMPLE_FILE%"

:: 6. Refresh Explorer
taskkill /IM explorer.exe /F >nul
start explorer.exe

:: Done
exit
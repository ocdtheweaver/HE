@echo off
echo Setting up HE_New Compiler...

:: Create directory structure
if not exist "src\main\kotlin\ast" mkdir src\main\kotlin\ast
if not exist "src\main\kotlin\visitor" mkdir src\main\kotlin\visitor
if not exist "src\main\kotlin\compiler" mkdir src\main\kotlin\compiler
if not exist "src\main\antlr" mkdir src\main\antlr
if not exist "lib" mkdir lib

echo Directory structure created.
echo Please make sure you have:
echo - Java JDK 8+ installed
echo - Your HE_New.g4 file in src\main\antlr\
echo - antlr.jar in lib\ folder
echo.
echo Then run: gradlew build

pause
@echo off
echo Compiling HE_New Compiler...
echo.

:: Check for required files
if not exist "lib\antlr.jar" (
    echo Error: lib\antlr.jar not found!
    pause
    exit /b 1
)

if not exist "src\main\antlr\HE_New.g4" (
    echo Error: src\main\antlr\HE_New.g4 not found!
    pause
    exit /b 1
)

:: Create output directory  
if not exist "src\main\java\he_new\parser" (
    mkdir "src\main\java\he_new\parser" 2>nul
)

:: Display grammar for debugging
echo ===== GRAMMAR CHECK =====
echo First few lines of your grammar:
echo.
type "src\main\antlr\HE_New.g4" | find /n /v "" | findstr /b "[1-5]:"
echo.

:: Compile ANTLR grammar
echo ===== GENERATING PARSER =====
echo Running ANTLR...

java -cp "lib\antlr.jar" org.antlr.v4.Tool -visitor -package he_new.parser -o src\main\java\he_new\parser "src\main\antlr\HE_New.g4"

if %ERRORLEVEL% neq 0 (
    echo.
    echo ===== PARSER GENERATION FAILED =====
    echo.
    echo Trying simpler approach...
    
    :: Try without quotes around the path
    java -cp lib\antlr.jar org.antlr.v4.Tool -visitor -package he_new.parser -o src/main/java/he_new/parser src/main/antlr/HE_New.g4
    
    if %ERRORLEVEL% neq 0 (
        echo.
        echo Still failed. Common issues:
        echo 1. Grammar syntax errors
        echo 2. ANTLR version mismatch
        echo 3. File encoding issues
        echo.
        echo Try creating a minimal grammar first.
        pause
        exit /b 1
    )
)

echo.
echo SUCCESS: ANTLR parser generated!
echo.
echo Generated files:
dir "src\main\java\he_new\parser\*.java" /b
echo.
pause
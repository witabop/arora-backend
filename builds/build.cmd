@echo off
setlocal
REM Ensure you have zip installed (e.g., via Chocolatey: choco install zip)

REM Set lambda name from first argument
set "LAMBDA_NAME=%~1"
if "%LAMBDA_NAME%"=="" (
    echo Error: You must specify a lambda name
    exit /b 1
)

cd /d "lambdas\%LAMBDA_NAME%" || exit /b 1

REM Build Go binary
set GOOS=linux
set GOARCH=amd64
go build -o ..\..\builds\zips\%LAMBDA_NAME%\bootstrap main.go
if %errorlevel% neq 0 (
    echo bootstrap failed
    exit /b 1
)

cd /d "..\..\builds\zips\%LAMBDA_NAME%" || exit /b 1

if not exist bootstrap (
    echo Error: bootstrap not created in %LAMBDA_NAME%
    exit /b 1
)

REM Create zip package
if exist %LAMBDA_NAME%.zip del %LAMBDA_NAME%.zip
zip -j %LAMBDA_NAME%.zip bootstrap
if %errorlevel% neq 0 (
    echo Zip failed
    exit /b 1
)

if not exist %LAMBDA_NAME%.zip (
    echo Error: %LAMBDA_NAME%.zip not created from bootstrap
    exit /b 1
)

REM Update Lambda function code
aws lambda update-function-code --function-name "arora-search-%LAMBDA_NAME%" --zip-file "fileb://%LAMBDA_NAME%.zip"
if %errorlevel% neq 0 (
    echo Function code update failed
    exit /b 1
)

echo Waiting for function update to complete...
aws lambda wait function-updated --function-name "arora-search-%LAMBDA_NAME%"
if %errorlevel% neq 0 (
    echo Wait for function update timed out
    exit /b 1
)

REM Publish version
aws lambda publish-version --function-name "arora-search-%LAMBDA_NAME%" --description "DEV Release"
if %errorlevel% neq 0 (
    echo Function version publishing failed
    exit /b 1
)

endlocal
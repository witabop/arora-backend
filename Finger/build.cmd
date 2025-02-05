@echo off
setlocal enabledelayedexpansion

set GOOS=linux
set GOARCH=amd64
set CGO_ENABLED=0
set LAMBDA_FUNCTION_NAME=RUFinger

:: Build the binary
echo Building Lambda binary...
go build -mod=readonly -o bootstrap main.go
if %errorlevel% neq 0 (
    echo Build failed!
    del bootstrap
    exit /b 1
)

:: Create ZIP
echo Creating deployment package...
powershell -Command "Compress-Archive -Force -Path bootstrap -DestinationPath function.zip"
if not exist function.zip (
    echo ZIP creation failed!
    del bootstrap
    exit /b 1
)

:: Deploy to Lambda
echo Deploying to AWS Lambda...
aws lambda update-function-code ^
    --function-name %LAMBDA_FUNCTION_NAME% ^
    --zip-file fileb://function.zip ^
    --profile vscode-cli > nul 2>&1

if %errorlevel% neq 0 (
    echo Deployment failed!
    del bootstrap
    del function.zip
    exit /b 1
)

:: Cleanup
echo Cleaning up...
del bootstrap
del function.zip

echo Successfully deployed to %LAMBDA_FUNCTION_NAME%
endlocal
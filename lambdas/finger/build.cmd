@REM @echo off
@REM setlocal enabledelayedexpansion

@REM set GOOS=linux
@REM set GOARCH=amd64
@REM set CGO_ENABLED=0
@REM set LAMBDA_FUNCTION_NAME=RUFinger

@REM :: Build the binary
@REM echo Building Lambda binary...
@REM go build -mod=readonly -o bootstrap main.go
@REM if %errorlevel% neq 0 (
@REM     echo Build failed!
@REM     del bootstrap
@REM     exit /b 1
@REM )

@REM :: Create ZIP
@REM echo Creating deployment package...
@REM powershell -Command "Compress-Archive -Force -Path bootstrap -DestinationPath function.zip"
@REM if not exist function.zip (
@REM     echo ZIP creation failed!
@REM     del bootstrap
@REM     exit /b 1
@REM )

@REM :: Deploy to Lambda
@REM echo Deploying to AWS Lambda...
@REM aws lambda update-function-code ^
@REM     --function-name %LAMBDA_FUNCTION_NAME% ^
@REM     --zip-file fileb://function.zip ^
@REM     --profile vscode-cli > nul 2>&1

@REM if %errorlevel% neq 0 (
@REM     echo Deployment failed!
@REM     del bootstrap
@REM     del function.zip
@REM     exit /b 1
@REM )

@REM :: Cleanup
@REM echo Cleaning up...
@REM del bootstrap
@REM del function.zip

@REM echo Successfully deployed to %LAMBDA_FUNCTION_NAME%
@REM endlocal
@echo off

REM check if the number of arguments is correct
if "%~2"=="" (
    echo Usage: %0 ^<module_name^> ^<service_name^>
    echo Example: %0 MyTest my_test_service
    exit /b 1
)

set MODULE_NAME=%1
set SERVICE_NAME=%2

echo Generating project with Module: %MODULE_NAME%, Service: %SERVICE_NAME%

hz new --mod=%MODULE_NAME% ^
    --service=%SERVICE_NAME% ^
    --customize_layout=./template/layout.yaml ^
    --customize_package=./template/package.yaml ^
    --handler_dir=biz/handler ^
    --router_dir=biz/router ^
    --model_dir=hertz_gen ^
    -force

if %errorlevel% neq 0 (
    echo Failed to execute hz command.
    exit /b %errorlevel%
)

echo Running go mod tidy...
go mod tidy

if %errorlevel% neq 0 (
    echo Failed to execute go mod tidy.
    exit /b %errorlevel%
)

echo Script finished successfully.
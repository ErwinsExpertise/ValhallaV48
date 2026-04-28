@echo off
setlocal
powershell -NoProfile -ExecutionPolicy Bypass -File "%~dp0convert-wz-to-nx.ps1" %*
if errorlevel 1 (
  echo.
  echo Conversion helper failed.
  pause
  exit /b 1
)
echo.
echo Conversion helper finished.
pause

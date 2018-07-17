@echo off
start /B "" "%~dp0bin\telescope.exe" uninstall
choice /t 5 /d y /n >nul  
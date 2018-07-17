@echo off
start /B "" "%~dp0bin\telescope.exe" install
choice /t 5 /d y /n >nul   
 
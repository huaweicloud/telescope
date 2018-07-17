@echo off
echo Stopping agent...
sc stop telescoped
choice /t 5 /d y /n >nul  
 
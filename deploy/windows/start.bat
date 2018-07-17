@echo off
echo Agent is starting...
sc start telescoped
choice /t 5 /d y /n >nul  
 
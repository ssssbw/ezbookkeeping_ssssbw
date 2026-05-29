@echo off
cd /d "%~dp0package"
"%~dp0tmp\ezbookkeeping.exe" server run

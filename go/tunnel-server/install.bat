@echo off
rem run this script as admin

if not exist tunnel.exe (
    echo Build the example before installing by running "go build"
    goto :exit
)

sc create rps-local-tunnel binpath= "%CD%\tunnel.exe" start= auto DisplayName= "rps-local-tunnel"
sc description rps-local-tunnel "rps-local-tunnel"
net start rps-local-tunnel
sc query rps-local-tunnel

echo Check rps-local-tunnel.log

:exit
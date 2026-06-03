$ErrorActionPreference = "Stop"
$Root = Split-Path -Parent $MyInvocation.MyCommand.Path
$Exe = Join-Path $Root "bin\runner-monitor.exe"

if (-not (Test-Path -LiteralPath $Exe)) {
    & (Join-Path $Root "scripts\build.ps1") | Out-Null
}

& $Exe @args
exit $LASTEXITCODE


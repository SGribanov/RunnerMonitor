$ErrorActionPreference = "Stop"
$Root = Split-Path -Parent $MyInvocation.MyCommand.Path
$Exe = Join-Path $Root "bin\runner-monitor.exe"

$NeedsBuild = -not (Test-Path -LiteralPath $Exe)
if (-not $NeedsBuild) {
    $ExeTime = (Get-Item -LiteralPath $Exe).LastWriteTimeUtc
    $NewerSource = Get-ChildItem -LiteralPath $Root -Recurse -File |
        Where-Object { $_.Extension -eq ".go" -or $_.Name -in @("go.mod", "go.sum") } |
        Where-Object { $_.LastWriteTimeUtc -gt $ExeTime } |
        Select-Object -First 1
    $NeedsBuild = $null -ne $NewerSource
}

if ($NeedsBuild) {
    & (Join-Path $Root "scripts\build.ps1") | Out-Null
}

& $Exe @args
exit $LASTEXITCODE

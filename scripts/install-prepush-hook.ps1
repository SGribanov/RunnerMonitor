param(
    [string]$RepoPath = (Get-Location).Path
)

$ErrorActionPreference = "Stop"
$RepoPath = (Resolve-Path -LiteralPath $RepoPath).Path
$GitDir = (& git -C $RepoPath rev-parse --git-dir).Trim()
if (-not [System.IO.Path]::IsPathRooted($GitDir)) {
    $GitDir = Join-Path $RepoPath $GitDir
}

$HookDir = Join-Path $GitDir "hooks"
$HookPath = Join-Path $HookDir "pre-push"
New-Item -ItemType Directory -Force -Path $HookDir | Out-Null

$Hook = @'
#!/bin/sh
powershell.exe -NoProfile -ExecutionPolicy Bypass -File "D:/Repos/RunnerMonitor/runner-monitor.ps1" --start-current
'@

Set-Content -LiteralPath $HookPath -Value $Hook -NoNewline -Encoding UTF8
Write-Output "Installed pre-push hook: $HookPath"


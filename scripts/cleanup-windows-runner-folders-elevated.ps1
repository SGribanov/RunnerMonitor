$ErrorActionPreference = "Stop"

$LogPath = "C:\Repos\RunnerMonitor\reports\runner-cleanup-2026-06-04.log"
New-Item -ItemType Directory -Path (Split-Path -Parent $LogPath) -Force | Out-Null

function Write-Log {
    param([string]$Message)
    $line = "[{0}] {1}" -f (Get-Date -Format "yyyy-MM-dd HH:mm:ss"), $Message
    $line | Tee-Object -FilePath $LogPath -Append
}

function Assert-Admin {
    $principal = New-Object Security.Principal.WindowsPrincipal([Security.Principal.WindowsIdentity]::GetCurrent())
    if (-not $principal.IsInRole([Security.Principal.WindowsBuiltInRole]::Administrator)) {
        throw "This script must run from elevated PowerShell."
    }
}

function Get-GitHubRunner {
    param([string]$Repo, [string]$Runner)
    $response = gh api "repos/$Repo/actions/runners" | ConvertFrom-Json
    $match = $response.runners | Where-Object { $_.name -eq $Runner } | Select-Object -First 1
    if (-not $match) {
        throw "GitHub runner $Repo/$Runner not found."
    }
    return $match
}

function Assert-Idle {
    param([pscustomobject]$Runner)
    $ghRunner = Get-GitHubRunner -Repo $Runner.Repo -Runner $Runner.Name
    if ($ghRunner.busy) {
        throw "$($Runner.Repo) $($Runner.Name) is busy; aborting cleanup."
    }
    Write-Log "Idle check passed: $($Runner.Repo) $($Runner.Name) status=$($ghRunner.status) busy=$($ghRunner.busy)"
}

function Wait-ServiceState {
    param([string]$Name, [string]$State, [int]$TimeoutSeconds = 90)
    $deadline = (Get-Date).AddSeconds($TimeoutSeconds)
    while ((Get-Date) -lt $deadline) {
        $svc = Get-Service -Name $Name -ErrorAction Stop
        if ($svc.Status.ToString() -eq $State) {
            return
        }
        Start-Sleep -Seconds 1
    }
    throw "Service $Name did not reach $State."
}

function Stop-ManualRunner {
    param([string]$Root)
    $prefix = (Resolve-Path -LiteralPath $Root).Path.TrimEnd("\") + "\"
    $names = @("Runner.Listener.exe", "Runner.Worker.exe", "Runner.PluginHost.exe")
    $procs = Get-CimInstance Win32_Process | Where-Object {
        $_.ExecutablePath -and
        $_.ExecutablePath.StartsWith($prefix, [System.StringComparison]::OrdinalIgnoreCase) -and
        $names -contains $_.Name
    }
    foreach ($proc in $procs) {
        Stop-Process -Id $proc.ProcessId -Force
        Write-Log "Stopped manual runner process pid=$($proc.ProcessId) path=$($proc.ExecutablePath)"
    }
}

function Start-ManualRunner {
    param([string]$Root)
    $run = Join-Path $Root "run.cmd"
    if (-not (Test-Path -LiteralPath $run)) {
        throw "run.cmd not found at $run"
    }
    Start-Process -FilePath $run -WorkingDirectory $Root -WindowStyle Hidden
    Write-Log "Started manual runner from $Root"
}

function Clear-RunnerFolder {
    param([string]$Root)
    $work = Join-Path $Root "_work"
    if (Test-Path -LiteralPath $work) {
        Get-ChildItem -LiteralPath $work -Force -ErrorAction SilentlyContinue |
            Remove-Item -Recurse -Force -ErrorAction SilentlyContinue
        Write-Log "Cleared _work: $work"
    }
    Get-ChildItem -LiteralPath $Root -Force -File -Filter "actions-runner*.zip" -ErrorAction SilentlyContinue |
        Remove-Item -Force -ErrorAction SilentlyContinue
    Write-Log "Removed installer zip files under $Root"
}

function Wait-GitHubOnline {
    param([pscustomobject]$Runner, [int]$TimeoutSeconds = 120)
    $deadline = (Get-Date).AddSeconds($TimeoutSeconds)
    while ((Get-Date) -lt $deadline) {
        $ghRunner = Get-GitHubRunner -Repo $Runner.Repo -Runner $Runner.Name
        if ($ghRunner.status -eq "online" -and -not $ghRunner.busy) {
            Write-Log "Online check passed: $($Runner.Repo) $($Runner.Name)"
            return
        }
        Start-Sleep -Seconds 3
    }
    throw "$($Runner.Repo) $($Runner.Name) did not become online/busy=false."
}

Assert-Admin
Write-Log "Starting Windows runner cleanup."

$runners = @(
    [pscustomobject]@{ Repo = "SGribanov/AU"; Name = "windows-local"; Root = "C:\Runners\SGribanov-AU\windows-local"; Mode = "manual" },
    [pscustomobject]@{ Repo = "SGribanov/BackTester"; Name = "backtester-runner"; Root = "C:\Runners\SGribanov-BackTester\backtester-runner"; Mode = "manual" },
    [pscustomobject]@{ Repo = "SGribanov/MyCloneOsEngine"; Name = "mycloneosengine-windows-local"; Root = "C:\Runners\SGribanov-MyCloneOsEngine\mycloneosengine-windows-local"; Mode = "manual" },
    [pscustomobject]@{ Repo = "SGribanov/IdeaBox"; Name = "ideabox-runner"; Root = "C:\Runners\SGribanov-IdeaBox\ideabox-runner"; Mode = "service"; ServiceName = "actions.runner.SGribanov-IdeaBox.ideabox-runner" },
    [pscustomobject]@{ Repo = "SGribanov/DeltaG"; Name = "deltag-win"; Root = "C:\Runners\SGribanov-DeltaG\deltag-win"; Mode = "service"; ServiceName = "actions.runner.SGribanov-DeltaG.deltag-win" }
)

try {
    foreach ($runner in $runners) {
        Assert-Idle -Runner $runner
    }

    foreach ($runner in $runners) {
        if ($runner.Mode -eq "service") {
            Stop-Service -Name $runner.ServiceName -Force
            Wait-ServiceState -Name $runner.ServiceName -State "Stopped"
            Write-Log "Stopped service $($runner.ServiceName)"
        } else {
            Stop-ManualRunner -Root $runner.Root
        }
    }

    foreach ($runner in $runners) {
        Clear-RunnerFolder -Root $runner.Root
    }

    $oldDeltaG = "C:\github-runners\deltag"
    if (Test-Path -LiteralPath $oldDeltaG) {
        Remove-Item -LiteralPath $oldDeltaG -Recurse -Force
        Write-Log "Removed old DeltaG folder: $oldDeltaG"
    }

    if (Test-Path -LiteralPath "C:\Runners-backup") {
        Get-ChildItem -LiteralPath "C:\Runners-backup" -Force -File |
            Remove-Item -Force
        Write-Log "Removed Windows backup archives under C:\Runners-backup"
    }

    foreach ($runner in $runners) {
        if ($runner.Mode -eq "service") {
            & sc.exe config $runner.ServiceName start= demand | Out-Null
            Start-Service -Name $runner.ServiceName
            Wait-ServiceState -Name $runner.ServiceName -State "Running"
            Write-Log "Started service $($runner.ServiceName)"
        } else {
            Start-ManualRunner -Root $runner.Root
        }
    }

    Start-Sleep -Seconds 8
    foreach ($runner in $runners) {
        Wait-GitHubOnline -Runner $runner
    }

    Write-Log "Windows runner cleanup completed."
} catch {
    Write-Log "ERROR: $($_.Exception.Message)"
    throw
}

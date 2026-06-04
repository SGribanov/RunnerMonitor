$ErrorActionPreference = "Stop"

$roots = @(
    "C:\Runners\SGribanov-AU\windows-local",
    "C:\Runners\SGribanov-BackTester\backtester-runner",
    "C:\Runners\SGribanov-MyCloneOsEngine\mycloneosengine-windows-local"
)

$names = @("Runner.Listener.exe", "Runner.Worker.exe", "Runner.PluginHost.exe")
foreach ($root in $roots) {
    $prefix = (Resolve-Path -LiteralPath $root).Path.TrimEnd("\") + "\"
    $procs = Get-CimInstance Win32_Process | Where-Object {
        $_.Name -in $names -and
        ($_.ExecutablePath -and $_.ExecutablePath.StartsWith($prefix, [System.StringComparison]::OrdinalIgnoreCase) -or
         $_.CommandLine -and $_.CommandLine.Contains($root))
    }
    foreach ($proc in $procs) {
        Stop-Process -Id $proc.ProcessId -Force
    }
}

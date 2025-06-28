param(
    [ValidateSet("start", "stop")]
    [string]$Action = "start"
)

# 配置项：要控制的程序名称（必须唯一）
$ProgramName = "l4d2serverquery_windows_amd64"  # 替换为实际进程名（不带.exe）
$ProgramDir = "C:\zzk\app\l4d2serverquery_release"

pushd $ProgramDir

# 核心逻辑
switch ($Action.ToLower()) {
    "start" {
        # 先停止同名进程（确保唯一性）
        $existingProcess = Get-Process -Name $ProgramName -ErrorAction SilentlyContinue
        if ($existingProcess) {
            $existingProcess | Stop-Process -Force
            Write-Host "Stopped existing process：$ProgramName"
        }

        # 启动新进程（根据实际情况调整路径和参数）
        # 我没有添加参数, 所以去除 -ArgumentList @()
        Start-Process -FilePath "$ProgramDir\$ProgramName.exe" `
                     -WindowStyle Hidden
        Write-Host "Started new process：$ProgramName"
    }

    "stop" {
        $targetProcess = Get-Process -Name $ProgramName -ErrorAction SilentlyContinue
        if ($targetProcess) {
            $targetProcess | Stop-Process -Force
            Write-Host "Stopped process：$ProgramName"
        } else {
            Write-Host "Process not found：$ProgramName"
        }
    }
}

popd
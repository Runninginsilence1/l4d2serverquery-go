# 先停止本地的服务
# .\localrun.ps1 -Action stop

# 前端
#Set-Location "C:\zzk\code\vscode\l4d2"
#pnpm build
#Write-Host "frontend build success"
#Copy-Item -Path "dist" -Destination "C:\Users\zzk\GolandProjects\l4d2serverquery-go\resources" -Recurse -Force

# 后端
Set-Location "C:\Users\zzk\GolandProjects\l4d2serverquery-go"
go env -w GOOS=windows
go env -w GOARCH=amd64

$bin_name = "l4d2serverquery_windows_amd64.exe"
$target_dir = "C:\zzk\app\l4d2serverquery_release"

go build -tags embed_ui,myownpc -o $bin_name -v .

# 移动文件和数据库
Copy-Item -Path $bin_name -Destination $target_dir -Force

Remove-Item -Path $bin_name -Force

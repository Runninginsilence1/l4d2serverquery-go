# 前端
Set-Location "C:\zzk\code\vscode\l4d2"
pnpm build
Write-Host "frontend build success"
Copy-Item -Path "dist" -Destination "C:\Users\zzk\GolandProjects\l4d2serverquery-go\resources" -Recurse -Force

# 后端
Set-Location "C:\Users\zzk\GolandProjects\l4d2serverquery-go"
go env -w GOOS=linux
go env -w GOARCH=amd64

$bin_name = "l4d2serverquery_linux_amd64"
$target_dir = "C:\zzk\app\l4d2serverquery_release"

go mod tidy
go build -o $bin_name -v .

if (-not $?) {
    Write-Error "Go编译失败"
    exit 1
}

# 移动文件和数据库
Copy-Item -Path "l4d2serverquery_linux_amd64" -Destination $target_dir -Force

Remove-Item -Path "l4d2serverquery_linux_amd64" -Force

Write-Host "backend of linux-amd64 build success" -ForegroundColor Yellow
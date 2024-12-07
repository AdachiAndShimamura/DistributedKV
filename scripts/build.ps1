$Root=Split-Path -Path $(Split-Path -path $PSCommandPath -Parent) -Parent
Write-Host $Root
Set-Location $Root
$command="protoc -I=proto/proto --go_opt=paths=source_relative --go_out=proto/gen --go-grpc_opt=paths=source_relative --go-grpc_out=proto/gen ./proto/proto/*.proto"
Invoke-Expression $command
$VERSION = "1.1.0"

$FileName = "ip2locationio_$($VERSION)_windows_amd64"
$ZipFileName = "$($FileName).zip"
$FolderName = "ip2location-io-cli"
$ExeName = "ip2locationio.exe"

Invoke-WebRequest -Uri "https://github.com/ip2location/ip2location-io-cli/releases/download/v$VERSION/$FileName.zip" -OutFile ./$ZipFileName
Unblock-File ./$ZipFileName
Expand-Archive -Path ./$ZipFileName  -DestinationPath $env:LOCALAPPDATA\$FolderName -Force

if (Test-Path "$env:LOCALAPPDATA\$FolderName\$ExeName") {
  Remove-Item "$env:LOCALAPPDATA\$FolderName\$ExeName"
}
Rename-Item -Path "$env:LOCALAPPDATA\$FolderName\$FileName.exe" -NewName "$ExeName"

$PathContent = [Environment]::GetEnvironmentVariable('path', 'Machine')
$IP2LocationIOPath = "$env:LOCALAPPDATA\$FolderName"

if ($PathContent -ne $null) {
  if (-Not($PathContent -split ';' -contains $IP2LocationIOPath)) {
    [System.Environment]::SetEnvironmentVariable("PATH", $Env:Path + ";$env:LOCALAPPDATA\$FolderName", "Machine")
  }
}
else {
  [System.Environment]::SetEnvironmentVariable("PATH", $Env:Path + ";$env:LOCALAPPDATA\$FolderName", "Machine")
}

Remove-Item -Path ./$ZipFileName
"You can use ip2locationio now."

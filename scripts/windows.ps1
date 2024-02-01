$RESOURCE = "https://api.github.com/repos/ip2location/ip2location-io-cli/releases/latest"

$GitHub = Invoke-RestMethod -Method Get -URI $RESOURCE -ErrorAction SilentlyContinue

if ($null -eq $GitHub.tag_name) {
  "Error: Unable to get latest version."
}
else {
  $TagName = $GitHub.tag_name.Trim()

  if ($TagName -match 'v\d+\.\d+\.\d+') {
    $Version = $TagName.substring(1)

    $FileName = "ip2locationio_$($Version)_windows_amd64"
    $ZipFileName = "$($FileName).zip"
    $FolderName = "ip2location-io-cli"
    $ExeName = "ip2locationio.exe"

    Invoke-WebRequest -Uri "https://github.com/ip2location/ip2location-io-cli/releases/download/v$Version/$FileName.zip" -OutFile ./$ZipFileName
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
  }
  else {
    "Error: Invalid tag name found."
  }
}

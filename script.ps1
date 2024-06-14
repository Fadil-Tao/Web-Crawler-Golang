
$startTime = Get-Date


& .\main.exe https://pasangiconnet.com/


$endTime = Get-Date


$duration = $endTime - $startTime

Write-Output "Lama eksekusi: $($duration.TotalSeconds) detik"

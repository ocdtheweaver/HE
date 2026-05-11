# Download gradle-wrapper.jar
$url = "https://github.com/gradle/gradle/raw/master/gradle/wrapper/gradle-wrapper.jar"
$output = "gradle\wrapper\gradle-wrapper.jar"

# Create directory if it doesn't exist
if (!(Test-Path "gradle\wrapper")) {
    New-Item -ItemType Directory -Path "gradle\wrapper" -Force
}

# Download the file
Invoke-WebRequest -Uri $url -OutFile $output
Write-Host "Downloaded gradle-wrapper.jar"
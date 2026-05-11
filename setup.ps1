# Create directory structure
$dirs = @(
    "src/main/kotlin/ast",
    "src/main/kotlin/visitor", 
    "src/main/kotlin/compiler",
    "src/main/antlr",
    "gradle/wrapper",
    "lib"
)

foreach ($dir in $dirs) {
    if (!(Test-Path $dir)) {
        New-Item -ItemType Directory -Path $dir -Force
        Write-Host "Created directory: $dir"
    }
}

# Create gradle wrapper properties
$gradleProps = @"
distributionBase=GRADLE_USER_HOME
distributionPath=wrapper/dists
distributionUrl=https\://services.gradle.org/distributions/gradle-8.5-bin.zip
networkTimeout=10000
validateDistributionUrl=true
zipStoreBase=GRADLE_USER_HOME
zipStorePath=wrapper/dists
"@

$gradleProps | Out-File -FilePath "gradle/wrapper/gradle-wrapper.properties" -Encoding UTF8

Write-Host "Project structure created!"
Write-Host "Please download gradle-wrapper.jar and place it in gradle/wrapper/"
Write-Host "Or run the download script if you created it."
// build.gradle.kts
plugins {
    kotlin("jvm") version "1.9.0"
    antlr
}

repositories {
    mavenCentral()
}

dependencies {
    antlr("org.antlr:antlr4:4.13.0")
    implementation("org.antlr:antlr4-runtime:4.13.0")
    implementation(kotlin("stdlib"))
}

tasks {
    compileKotlin {
        dependsOn(generateGrammarSource)
    }
    
    generateGrammarSource {
        arguments = arguments + listOf("-visitor", "-package", "he_new.parser")
        outputDirectory = file("src/main/java/he_new/parser")
    }
}

sourceSets {
    main {
        java {
            srcDir("src/main/kotlin")
        }
    }
}

// For Windows compatibility
tasks.withType<JavaExec> {
    systemProperty("file.encoding", "UTF-8")
}

java {
    toolchain {
        languageVersion.set(JavaLanguageVersion.of(17)) // Or 21
    }
}
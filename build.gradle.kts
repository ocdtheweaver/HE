plugins {
    antlr
    kotlin("jvm") version "1.9.0"
}

repositories {
    mavenCentral()
}

dependencies {
    antlr("org.antlr:antlr4:4.13.0")
    implementation("org.antlr:antlr4-runtime:4.13.0")
}

tasks.generateGrammarSource {
    arguments = arguments + listOf("-visitor", "-package", "he.parser")
}
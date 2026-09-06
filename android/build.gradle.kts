plugins {
    id("com.android.application") version "9.3.2" apply false
    id("org.jetbrains.kotlin.plugin.compose") version "2.4.10" apply false
}

subprojects {
    if (tasks.findByName("prepareKotlinBuildScriptModel") == null) {
        tasks.register("prepareKotlinBuildScriptModel") {
            // Satisfies IDE Kotlin DSL sync on subprojects
        }
    }
    if (tasks.findByName("wrapper") == null) {
        tasks.register<Wrapper>("wrapper") {
            gradleVersion = "9.7.1"
        }
    }
}

plugins {
    id("com.android.application") version "8.2.2" apply false
    id("org.jetbrains.kotlin.android") version "1.9.22" apply false
}

subprojects {
    if (tasks.findByName("prepareKotlinBuildScriptModel") == null) {
        tasks.register("prepareKotlinBuildScriptModel") {
            // Satisfies IDE Kotlin DSL sync on subprojects
        }
    }
    if (tasks.findByName("wrapper") == null) {
        tasks.register<Wrapper>("wrapper") {
            gradleVersion = "8.10.2"
        }
    }
}

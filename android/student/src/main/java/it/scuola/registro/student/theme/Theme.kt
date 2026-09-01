package it.scuola.registro.student.theme

import androidx.compose.foundation.isSystemInDarkTheme
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.darkColorScheme
import androidx.compose.material3.lightColorScheme
import androidx.compose.runtime.Composable
import androidx.compose.ui.graphics.Color

// Youthful, Vibrant Palette for Student App
val StudentViolet = Color(0xFF7C3AED)
val StudentIndigo = Color(0xFF4F46E5)
val StudentEmerald = Color(0xFF10B981)
val StudentAmber = Color(0xFFF59E0B)

val DarkBackground = Color(0xFF0F172A)
val DarkSurface = Color(0xFF1E293B)
val LightBackground = Color(0xFFF8FAFC)
val LightSurface = Color(0xFFFFFFFF)

private val DarkColorScheme = darkColorScheme(
    primary = StudentViolet,
    secondary = StudentEmerald,
    tertiary = StudentAmber,
    background = DarkBackground,
    surface = DarkSurface
)

private val LightColorScheme = lightColorScheme(
    primary = StudentIndigo,
    secondary = StudentEmerald,
    tertiary = StudentAmber,
    background = LightBackground,
    surface = LightSurface
)

@Composable
fun RegistroStudentTheme(
    darkTheme: Boolean = isSystemInDarkTheme(),
    content: @Composable () -> Unit
) {
    val colorScheme = if (darkTheme) DarkColorScheme else LightColorScheme

    MaterialTheme(
        colorScheme = colorScheme,
        content = content
    )
}

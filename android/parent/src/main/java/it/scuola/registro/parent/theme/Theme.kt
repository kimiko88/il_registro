package it.scuola.registro.parent.theme

import androidx.compose.foundation.isSystemInDarkTheme
import androidx.compose.foundation.layout.defaultMinSize
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.Typography
import androidx.compose.material3.darkColorScheme
import androidx.compose.material3.lightColorScheme
import androidx.compose.runtime.Composable
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.text.TextStyle
import androidx.compose.ui.text.font.FontFamily
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp

// Warm, Reassuring Palette for Parent App
val ParentTeal = Color(0xFF0D9488)
val ParentAmber = Color(0xFFD97706)
val ParentRose = Color(0xFFE11D48)

val DarkBackground = Color(0xFF0F172A)
val DarkSurface = Color(0xFF1E293B)
val LightBackground = Color(0xFFF8FAFC)
val LightSurface = Color(0xFFFFFFFF)

// WCAG 2.2 AAA High Contrast Palette
val HighContrastBlack = Color(0xFF000000)
val HighContrastWhite = Color(0xFFFFFFFF)
val HighContrastYellow = Color(0xFFFFD700)
val HighContrastCyan = Color(0xFF00E5FF)

private val DarkColorScheme = darkColorScheme(
    primary = ParentTeal,
    secondary = ParentAmber,
    tertiary = ParentRose,
    background = DarkBackground,
    surface = DarkSurface
)

private val LightColorScheme = lightColorScheme(
    primary = ParentTeal,
    secondary = ParentAmber,
    tertiary = ParentRose,
    background = LightBackground,
    surface = LightSurface
)

private val HighContrastDarkColorScheme = darkColorScheme(
    primary = HighContrastYellow,
    secondary = HighContrastCyan,
    tertiary = HighContrastWhite,
    background = HighContrastBlack,
    surface = Color(0xFF121212),
    onPrimary = HighContrastBlack,
    onBackground = HighContrastWhite,
    onSurface = HighContrastWhite
)

private val HighContrastLightColorScheme = lightColorScheme(
    primary = HighContrastBlack,
    secondary = Color(0xFF004D40),
    tertiary = Color(0xFF880E4F),
    background = HighContrastWhite,
    surface = HighContrastWhite,
    onPrimary = HighContrastWhite,
    onBackground = HighContrastBlack,
    onSurface = HighContrastBlack
)

fun createAccessibleTypography(dyslexiaFriendly: Boolean): Typography {
    val letterSpacing = if (dyslexiaFriendly) 0.15.sp else 0.sp
    val lineHeightMultiplier = if (dyslexiaFriendly) 1.5 else 1.2

    return Typography(
        bodyLarge = TextStyle(
            fontFamily = FontFamily.SansSerif,
            fontWeight = FontWeight.Normal,
            fontSize = 16.sp,
            lineHeight = (16 * lineHeightMultiplier).sp,
            letterSpacing = letterSpacing
        ),
        bodyMedium = TextStyle(
            fontFamily = FontFamily.SansSerif,
            fontWeight = FontWeight.Normal,
            fontSize = 14.sp,
            lineHeight = (14 * lineHeightMultiplier).sp,
            letterSpacing = letterSpacing
        ),
        titleLarge = TextStyle(
            fontFamily = FontFamily.SansSerif,
            fontWeight = FontWeight.Bold,
            fontSize = 22.sp,
            lineHeight = (22 * lineHeightMultiplier).sp,
            letterSpacing = letterSpacing
        )
    )
}

fun Modifier.accessibleTouchTarget(): Modifier = this.defaultMinSize(minWidth = 48.dp, minHeight = 48.dp)

@Composable
fun RegistroParentTheme(
    darkTheme: Boolean = isSystemInDarkTheme(),
    highContrast: Boolean = false,
    dyslexiaFriendly: Boolean = false,
    content: @Composable () -> Unit
) {
    val colorScheme = when {
        highContrast && darkTheme -> HighContrastDarkColorScheme
        highContrast && !darkTheme -> HighContrastLightColorScheme
        darkTheme -> DarkColorScheme
        else -> LightColorScheme
    }

    val typography = createAccessibleTypography(dyslexiaFriendly)

    MaterialTheme(
        colorScheme = colorScheme,
        typography = typography,
        content = content
    )
}

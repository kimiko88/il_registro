package it.scuola.registro.teacher

import androidx.compose.ui.test.*
import androidx.compose.ui.test.junit4.createAndroidComposeRule
import it.scuola.registro.teacher.MainActivity
import org.junit.Rule
import org.junit.Test

class TeacherE2ETest {

    @get:Rule
    val composeTestRule = createAndroidComposeRule<MainActivity>()

    @Test
    fun teacherFullE2EFlow_loginSignRegisterTakeAttendanceAndScrutiny() {
        // 1. Check Login Screen
        composeTestRule.onNodeWithText("Registro Docente").assertIsDisplayed()

        // 2. Input credentials
        composeTestRule.onNodeWithText("Email Docente").performTextInput("docente.bianchi@scuola.it")
        composeTestRule.onNodeWithText("Password").performTextInput("password123")

        // 3. Login
        composeTestRule.onNodeWithText("ACCEDI AL REGISTRO").performClick()

        // 4. Verify Register Header
        composeTestRule.onNodeWithText("Firma Registro di Classe").assertIsDisplayed()

        // 5. Navigate to Appello Tab
        composeTestRule.onNodeWithText("Appello").performClick()
        composeTestRule.onNodeWithText("Rilevazione Presenze (Appello Giornaliero)").assertIsDisplayed()

        // 6. Navigate to Scrutini Tab
        composeTestRule.onNodeWithText("Scrutini").performClick()
        composeTestRule.onNodeWithText("Gestione Scrutini & Scrutini Differiti").assertIsDisplayed()
        composeTestRule.onNodeWithText("Apri Sessione Scrutinio Differito").assertIsDisplayed()
    }
}

package it.scuola.registro.student

import androidx.compose.ui.test.*
import androidx.compose.ui.test.junit4.createAndroidComposeRule
import it.scuola.registro.student.MainActivity
import org.junit.Rule
import org.junit.Test

class StudentE2ETest {

    @get:Rule
    val composeTestRule = createAndroidComposeRule<MainActivity>()

    @Test
    fun studentFullE2EFlow_loginToDashboardNavigation() {
        // 1. Verify Login Screen is present
        composeTestRule.onNodeWithText("Registro Studente").assertIsDisplayed()
        composeTestRule.onNodeWithText("Accedi al tuo account").assertIsDisplayed()

        // 2. Input Email and Password
        composeTestRule.onNodeWithText("Email Studente").performTextInput("mario.rossi@studenti.it")
        composeTestRule.onNodeWithText("Password").performTextInput("password123")

        // 3. Click ACCEDI button
        composeTestRule.onNodeWithText("ACCEDI").performClick()

        // 4. Verify Dashboard Screen rendered
        composeTestRule.onNodeWithText("La Mia Dashboard").assertIsDisplayed()
        composeTestRule.onNodeWithText("Ultimi Voti Inseriti").assertIsDisplayed()

        // 5. Navigate to Grades tab
        composeTestRule.onNodeWithText("Voti").performClick()
        composeTestRule.onNodeWithText("I Miei Voti").assertIsDisplayed()

        // 6. Navigate to Agenda tab
        composeTestRule.onNodeWithText("Agenda").performClick()
        composeTestRule.onNodeWithText("Compiti e Scadenze Imminenti").assertIsDisplayed()

        // 7. Navigate to Report Card tab
        composeTestRule.onNodeWithText("Pagella").performClick()
        composeTestRule.onNodeWithText("Documento di Valutazione Ufficiale").assertIsDisplayed()
        composeTestRule.onNodeWithText("Scarica Pagella Ufficiale (PDF)").assertIsDisplayed()
    }
}

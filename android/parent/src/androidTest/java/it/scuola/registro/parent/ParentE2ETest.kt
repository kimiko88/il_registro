package it.scuola.registro.parent

import androidx.compose.ui.test.*
import androidx.compose.ui.test.junit4.createAndroidComposeRule
import it.scuola.registro.parent.MainActivity
import org.junit.Rule
import org.junit.Test

class ParentE2ETest {

    @get:Rule
    val composeTestRule = createAndroidComposeRule<MainActivity>()

    @Test
    fun parentFullE2EFlow_loginChildSwitchingAndJustifications() {
        // 1. Check Login Screen
        composeTestRule.onNodeWithText("Registro Genitori").assertIsDisplayed()
        composeTestRule.onNodeWithText("Accedi al portale di supervisione").assertIsDisplayed()

        // 2. Input credentials
        composeTestRule.onNodeWithText("Email Genitore").performTextInput("genitore@famiglia.it")
        composeTestRule.onNodeWithText("Password").performTextInput("password123")

        // 3. Login
        composeTestRule.onNodeWithText("ACCEDI").performClick()

        // 4. Check Child Supervision Header
        composeTestRule.onNodeWithText("Marco Rossi").assertIsDisplayed()

        // 5. Navigate to Justifications Tab
        composeTestRule.onNodeWithText("Giustifiche").performClick()
        composeTestRule.onNodeWithText("Assenze e Ritardi da Giustificare").assertIsDisplayed()

        // 6. Justify first pending absence
        composeTestRule.onAllNodesWithText("Giustifica")[0].performClick()
        composeTestRule.onNodeWithText("Giustificata").assertIsDisplayed()

        // 7. Navigate to Colloqui Tab
        composeTestRule.onNodeWithText("Colloqui").performClick()
        composeTestRule.onNodeWithText("Prof. Bianchi").assertIsDisplayed()
    }
}

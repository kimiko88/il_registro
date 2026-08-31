package it.scuola.registro.student

import androidx.compose.ui.test.*
import androidx.compose.ui.test.junit4.createComposeRule
import it.scuola.registro.student.ui.StudentELearningScreen
import org.junit.Rule
import org.junit.Test

class StudentAdvancedFlowE2ETest {

    @get:Rule
    val composeTestRule = createComposeRule()

    @Test
    fun elearningScreen_rendersQuizAndStartButton() {
        composeTestRule.setContent {
            StudentELearningScreen()
        }

        composeTestRule.onNodeWithText("E-Learning & Quiz Interattivi").assertIsDisplayed()
        composeTestRule.onNodeWithText("Avvia Quiz Interattivo").assertIsDisplayed()
    }
}

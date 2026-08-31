package it.scuola.registro.teacher

import androidx.compose.ui.test.*
import androidx.compose.ui.test.junit4.createComposeRule
import it.scuola.registro.teacher.ui.TeacherNotesScreen
import org.junit.Rule
import org.junit.Test

class TeacherAdvancedFlowE2ETest {

    @get:Rule
    val composeTestRule = createComposeRule()

    @Test
    fun teacherNotesScreen_rendersForm() {
        composeTestRule.setContent {
            TeacherNotesScreen()
        }

        composeTestRule.onNodeWithText("Note & Provvedimenti Disciplinari").assertIsDisplayed()
        composeTestRule.onNodeWithText("Registra Nota sul Giornale di Classe").assertIsDisplayed()
    }
}

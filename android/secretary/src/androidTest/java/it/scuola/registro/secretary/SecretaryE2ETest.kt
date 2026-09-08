package it.scuola.registro.secretary

import androidx.compose.ui.test.*
import androidx.compose.ui.test.junit4.createAndroidComposeRule
import it.scuola.registro.secretary.MainActivity
import org.junit.Rule
import org.junit.Test

class SecretaryE2ETest {

    @get:Rule
    val composeTestRule = createAndroidComposeRule<MainActivity>()

    @Test
    fun secretaryFullE2EFlow_loginAdminDashboardUsersCertificatesAndAudit() {
        // 1. Check Login Screen
        composeTestRule.onNodeWithText("Registro Segreteria").assertIsDisplayed()

        // 2. Input credentials
        composeTestRule.onNodeWithText("Email Segreteria").performTextInput("segreteria@scuola.it")
        composeTestRule.onNodeWithText("Password").performTextInput("password123")

        // 3. Login
        composeTestRule.onNodeWithText("ACCEDI AL PANNELLO").performClick()

        // 4. Verify Executive Overview
        composeTestRule.onNodeWithText("Stato Sistema & Operazioni").assertIsDisplayed()

        // 5. Navigate to Users Tab
        composeTestRule.onNodeWithText("Utenti").performClick()
        composeTestRule.onNodeWithText("Anagrafica Utenti Sistema").assertIsDisplayed()

        // 6. Navigate to Scrutini Tab
        composeTestRule.onNodeWithText("Scrutini").performClick()
        composeTestRule.onNodeWithText("Stato Scrutini di Tutte le Classi").assertIsDisplayed()

        // 7. Navigate to Certificates Tab
        composeTestRule.onNodeWithText("Certificati").performClick()
        composeTestRule.onNodeWithText("Emissione Certificati Ufficiali").assertIsDisplayed()

        // 8. Navigate to Audit Tab
        composeTestRule.onNodeWithText("Audit").performClick()
        composeTestRule.onNodeWithText("Registro Audit & Eventi di Sistema").assertIsDisplayed()
    }
}

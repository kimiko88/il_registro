package it.scuola.registro.student

import android.os.Bundle
import androidx.activity.ComponentActivity
import androidx.activity.compose.setContent
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.Surface
import androidx.compose.runtime.*
import androidx.compose.ui.Modifier
import it.scuola.registro.student.theme.RegistroStudentTheme
import it.scuola.registro.student.ui.StudentDashboardScreen
import it.scuola.registro.student.ui.StudentLoginScreen

class MainActivity : ComponentActivity() {
    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContent {
            RegistroStudentTheme {
                Surface(
                    modifier = Modifier.fillMaxSize(),
                    color = MaterialTheme.colorScheme.background
                ) {
                    var isLoggedIn by remember { mutableStateOf(false) }
                    var authToken by remember { mutableStateOf<String?>(null) }
                    var loggedInStudentName by remember { mutableStateOf("Mario Rossi") }

                    if (isLoggedIn) {
                        StudentDashboardScreen(
                            token = authToken,
                            studentName = loggedInStudentName,
                            onLogout = {
                                isLoggedIn = false
                                authToken = null
                            }
                        )
                    } else {
                        StudentLoginScreen(
                            onLoginSuccess = { token, studentName ->
                                authToken = token
                                loggedInStudentName = studentName
                                isLoggedIn = true
                            }
                        )
                    }
                }
            }
        }
    }
}

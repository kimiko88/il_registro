package it.scuola.registro.student

import android.os.Bundle
import android.view.WindowManager
import androidx.activity.ComponentActivity
import androidx.activity.compose.setContent
import kotlinx.coroutines.CoroutineScope
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.launch
import it.scuola.registro.student.network.HttpStudentApiService
import it.scuola.registro.student.network.logout
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.Surface
import androidx.compose.runtime.*
import androidx.compose.ui.Modifier
import it.scuola.registro.student.theme.RegistroStudentTheme
import it.scuola.registro.student.security.MdmConfigManager
import it.scuola.registro.student.ui.StudentDashboardScreen
import it.scuola.registro.student.ui.StudentLoginScreen

class MainActivity : ComponentActivity() {
    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        MdmConfigManager(this).applyMdmPolicies()
        // Student records and administrative data must not appear in screenshots or recents previews.
        window.addFlags(WindowManager.LayoutParams.FLAG_SECURE)
        setContent {
            RegistroStudentTheme {
                Surface(
                    modifier = Modifier.fillMaxSize(),
                    color = MaterialTheme.colorScheme.background
                ) {
                    var isLoggedIn by remember { mutableStateOf(false) }
                    var authToken by remember { mutableStateOf<String?>(null) }
                    var loggedInStudentName by remember { mutableStateOf("Studente") }

                    if (isLoggedIn) {
                        StudentDashboardScreen(
                            token = authToken,
                            studentName = loggedInStudentName,
                            onLogout = {
                                val tokenToRevoke = authToken
                                if (!tokenToRevoke.isNullOrBlank()) {
                                    CoroutineScope(Dispatchers.IO).launch {
                                        HttpStudentApiService().logout(tokenToRevoke)
                                    }
                                }
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

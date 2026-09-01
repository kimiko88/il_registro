package it.scuola.registro.parent

import android.os.Bundle
import androidx.activity.ComponentActivity
import androidx.activity.compose.setContent
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.Surface
import androidx.compose.runtime.*
import androidx.compose.ui.Modifier
import it.scuola.registro.parent.ui.ParentDashboardScreen
import it.scuola.registro.parent.ui.ParentLoginScreen

class MainActivity : ComponentActivity() {
    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        setContent {
            MaterialTheme {
                Surface(modifier = Modifier.fillMaxSize()) {
                    var isLoggedIn by remember { mutableStateOf(false) }
                    var authToken by remember { mutableStateOf<String?>(null) }

                    if (isLoggedIn) {
                        ParentDashboardScreen(
                            token = authToken,
                            onLogout = {
                                isLoggedIn = false
                                authToken = null
                            }
                        )
                    } else {
                        ParentLoginScreen(
                            onLoginSuccess = { token, _ ->
                                authToken = token
                                isLoggedIn = true
                            }
                        )
                    }
                }
            }
        }
    }
}

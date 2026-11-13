package it.scuola.registro.secretary

import android.os.Bundle
import android.view.WindowManager
import androidx.activity.ComponentActivity
import androidx.activity.compose.setContent
import kotlinx.coroutines.CoroutineScope
import kotlinx.coroutines.Dispatchers
import kotlinx.coroutines.launch
it.scuola.registro.secretary.network.HttpSecretaryApiService
import androidx.compose.foundation.layout.fillMaxSize
import androidx.compose.material3.MaterialTheme
import androidx.compose.material3.Surface
import androidx.compose.runtime.*
import androidx.compose.ui.Modifier
import it.scuola.registro.secretary.security.MdmConfigManager
import it.scuola.registro.secretary.ui.SecretaryDashboardScreen
import it.scuola.registro.secretary.ui.SecretaryLoginScreen

class MainActivity : ComponentActivity() {
    override fun onCreate(savedInstanceState: Bundle?) {
        super.onCreate(savedInstanceState)
        MdmConfigManager(this).applyMdmPolicies()
        // Student records and administrative data must not appear in screenshots or recents previews.
        window.addFlags(WindowManager.LayoutParams.FLAG_SECURE)
        setContent {
            MaterialTheme {
                Surface(modifier = Modifier.fillMaxSize()) {
                    var isLoggedIn by remember { mutableStateOf(false) }
                    var authToken by remember { mutableStateOf<String?>(null) }

                    if (isLoggedIn) {
                        SecretaryDashboardScreen(
                            token = authToken,
                            onLogout = {
                                val tokenToRevoke = authToken
                                if (!tokenToRevoke.isNullOrBlank()) {
                                    CoroutineScope(Dispatchers.IO).launch {
                                        HttpSecretaryApiService().logout(tokenToRevoke)
                                    }
                                }
                                isLoggedIn = false
                                authToken = null
                            }
                        )
                    } else {
                        SecretaryLoginScreen(
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

package it.scuola.registro.secretary.ui

import androidx.compose.foundation.layout.*
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.*
import androidx.compose.material3.*
import androidx.compose.runtime.Composable
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.res.stringResource
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import it.scuola.registro.secretary.R

data class SecretaryFascicoloData(
    val studentNameAndClass: String = "",
    val fiscalCode: String = "",
    val sidiCode: String = "",
    val vaccineStatus: String = ""
)

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun SecretaryFascicoloScreen(
    fascicolo: SecretaryFascicoloData = SecretaryFascicoloData(),
    onBack: () -> Unit = {}
) {
    Scaffold(
        topBar = {
            TopAppBar(
                title = { Text(stringResource(R.string.secretary_dashboard_title), fontWeight = FontWeight.Bold) },
                navigationIcon = {
                    IconButton(onClick = onBack) {
                        Icon(Icons.Default.ArrowBack, contentDescription = "Back")
                    }
                }
            )
        }
    ) { padding ->
        LazyColumn(
            modifier = Modifier
                .fillMaxSize()
                .padding(padding)
                .padding(16.dp),
            verticalArrangement = Arrangement.spacedBy(14.dp)
        ) {
            item {
                Card(
                    modifier = Modifier.fillMaxWidth(),
                    shape = RoundedCornerShape(14.dp),
                    colors = CardDefaults.cardColors(containerColor = MaterialTheme.colorScheme.primaryContainer)
                ) {
                    Column(modifier = Modifier.padding(16.dp)) {
                        Text(stringResource(R.string.secretary_dashboard_title), fontWeight = FontWeight.Bold, fontSize = 16.sp)
                        Text(stringResource(R.string.stats_title), fontSize = 12.sp, color = Color.DarkGray)
                    }
                }
            }

            if (fascicolo.studentNameAndClass.isNotEmpty()) {
                item {
                    Text(fascicolo.studentNameAndClass, fontWeight = FontWeight.Bold, fontSize = 16.sp)
                }

                item {
                    Card(modifier = Modifier.fillMaxWidth(), shape = RoundedCornerShape(12.dp)) {
                        Column(modifier = Modifier.padding(14.dp), verticalArrangement = Arrangement.spacedBy(8.dp)) {
                            Row(modifier = Modifier.fillMaxWidth(), horizontalArrangement = Arrangement.SpaceBetween) {
                                Text(stringResource(R.string.profile_title), fontWeight = FontWeight.SemiBold)
                                Text(fascicolo.fiscalCode, fontWeight = FontWeight.Bold)
                            }
                            Divider()
                            Row(modifier = Modifier.fillMaxWidth(), horizontalArrangement = Arrangement.SpaceBetween) {
                                Text(stringResource(R.string.secretary_dashboard_title), fontWeight = FontWeight.SemiBold)
                                Text(fascicolo.sidiCode, fontWeight = FontWeight.Bold, color = MaterialTheme.colorScheme.primary)
                            }
                            Divider()
                            Row(modifier = Modifier.fillMaxWidth(), horizontalArrangement = Arrangement.SpaceBetween) {
                                Text(stringResource(R.string.classes_title), fontWeight = FontWeight.SemiBold)
                                Text(fascicolo.vaccineStatus, fontWeight = FontWeight.Bold, color = Color(0xFF10B981))
                            }
                        }
                    }
                }
            }
        }
    }
}

package it.scuola.registro.parent.ui

import androidx.compose.foundation.layout.*
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.*
import androidx.compose.material3.*
import androidx.compose.runtime.Composable
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun ParentAbsenceMonitoringScreen(onBack: () -> Unit = {}) {
    Scaffold(
        topBar = {
            TopAppBar(
                title = { Text("Monitoraggio Limite Assenze", fontWeight = FontWeight.Bold) },
                navigationIcon = {
                    IconButton(onClick = onBack) {
                        Icon(Icons.Default.ArrowBack, contentDescription = "Indietro")
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
                        Text("Validità Anno Scolastico (D.P.R. 122/2009)", fontWeight = FontWeight.Bold, fontSize = 16.sp)
                        Text("Per la validità dell'anno è richiesta la frequenza di almeno il 75% del monte ore personalizzato (max 25% di assenze).", fontSize = 12.sp, color = Color.DarkGray)
                    }
                }
            }

            item {
                Text("Stato Monte Ore: Mario Rossi", fontWeight = FontWeight.Bold, fontSize = 16.sp)
            }

            item {
                Card(modifier = Modifier.fillMaxWidth(), shape = RoundedCornerShape(12.dp)) {
                    Column(modifier = Modifier.padding(14.dp), verticalArrangement = Arrangement.spacedBy(8.dp)) {
                        Row(
                            modifier = Modifier.fillMaxWidth(),
                            horizontalArrangement = Arrangement.SpaceBetween,
                            verticalAlignment = Alignment.CenterVertically
                        ) {
                            Text("Ore di Assenza Effettuate:", fontWeight = FontWeight.SemiBold)
                            Text("54 Ore / 247 Ore Max (5.4%)", fontWeight = FontWeight.Bold, color = Color(0xFF10B981))
                        }
                        LinearProgressIndicator(
                            progress = 0.22f, // 54/247 = 22% of the max allowed absences
                            modifier = Modifier.fillMaxWidth().padding(vertical = 4.dp),
                            color = Color(0xFF10B981)
                        )
                        Divider()
                        Row(modifier = Modifier.fillMaxWidth(), horizontalArrangement = Arrangement.SpaceBetween) {
                            Text("Monte Ore Totale Annuale:", fontWeight = FontWeight.SemiBold)
                            Text("990 Ore", fontWeight = FontWeight.Bold)
                        }
                        Row(modifier = Modifier.fillMaxWidth(), horizontalArrangement = Arrangement.SpaceBetween) {
                            Text("Stato Frequenza:", fontWeight = FontWeight.SemiBold)
                            Text("Regolare (Anno Pienamente Valido)", fontWeight = FontWeight.Bold, color = Color(0xFF10B981))
                        }
                    }
                }
            }
        }
    }
}

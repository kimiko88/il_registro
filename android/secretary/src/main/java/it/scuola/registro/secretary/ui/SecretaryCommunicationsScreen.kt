package it.scuola.registro.secretary.ui

import androidx.compose.foundation.layout.*
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.*
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun SecretaryCommunicationsScreen(onBack: () -> Unit = {}) {
    Scaffold(
        topBar = {
            TopAppBar(
                title = { Text("Invio Circolari & Comunicazioni", fontWeight = FontWeight.Bold) },
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
                        Text("Pubblicazione Circolari con Tracciamento Lettura", fontWeight = FontWeight.Bold, fontSize = 16.sp)
                        Text("Invio comunicazioni a docenti, genitori o studenti con tracking in tempo reale delle conferme di lettura e firme digitali.", fontSize = 12.sp, color = Color.DarkGray)
                    }
                }
            }

            item {
                Text("Ultime Circolari Inviate", fontWeight = FontWeight.Bold, fontSize = 16.sp)
            }

            item {
                Card(modifier = Modifier.fillMaxWidth(), shape = RoundedCornerShape(12.dp)) {
                    Column(modifier = Modifier.padding(14.dp)) {
                        Row(
                            modifier = Modifier.fillMaxWidth(),
                            horizontalArrangement = Arrangement.SpaceBetween,
                            verticalAlignment = Alignment.CenterVertically
                        ) {
                            Text("Circolare N. 215: Chiusura per Elezioni", fontWeight = FontWeight.Bold)
                            Surface(color = Color(0xFF10B981), shape = RoundedCornerShape(6.dp)) {
                                Text("89% Letti (420/472)", color = Color.White, modifier = Modifier.padding(horizontal = 6.dp, vertical = 2.dp), fontSize = 10.sp)
                            }
                        }
                        Text("Destinatari: Tutte le Famiglie & Docenti • Data Invio: Oggi alle 09:00", fontSize = 12.sp, color = MaterialTheme.colorScheme.primary, fontWeight = FontWeight.SemiBold, modifier = Modifier.padding(vertical = 4.dp))
                        Text("Firma di Presa Visione Richiesta: Sì • Termine: 2 Giugno", fontSize = 12.sp, color = Color.Gray)
                    }
                }
            }
        }
    }
}

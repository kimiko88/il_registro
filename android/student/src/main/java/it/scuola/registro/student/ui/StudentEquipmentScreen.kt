package it.scuola.registro.student.ui

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
fun StudentEquipmentScreen(onBack: () -> Unit = {}) {
    Scaffold(
        topBar = {
            TopAppBar(
                title = { Text("Armadietto & Dotazioni", fontWeight = FontWeight.Bold) },
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
                        Text("Assegnazione Risorse Scolastiche", fontWeight = FontWeight.Bold, fontSize = 16.sp)
                        Text("Gestione dell'armadietto personale, badge NFC per l'accesso ai laboratori e dispositivi in comodato d'uso.", fontSize = 12.sp, color = Color.DarkGray)
                    }
                }
            }

            item {
                Text("Risorse Assegnate", fontWeight = FontWeight.Bold, fontSize = 16.sp)
            }

            item {
                Card(modifier = Modifier.fillMaxWidth(), shape = RoundedCornerShape(12.dp)) {
                    Column(modifier = Modifier.padding(14.dp), verticalArrangement = Arrangement.spacedBy(8.dp)) {
                        Row(modifier = Modifier.fillMaxWidth(), horizontalArrangement = Arrangement.SpaceBetween) {
                            Text("Numero Armadietto:", fontWeight = FontWeight.SemiBold)
                            Text("ARM-104 (Piano 1, Corridoio Ovest)", fontWeight = FontWeight.Bold, color = MaterialTheme.colorScheme.primary)
                        }
                        Divider()
                        Row(modifier = Modifier.fillMaxWidth(), horizontalArrangement = Arrangement.SpaceBetween) {
                            Text("Codice Badge NFC:", fontWeight = FontWeight.SemiBold)
                            Text("NFC-99120-STUD", fontWeight = FontWeight.Bold)
                        }
                        Divider()
                        Row(modifier = Modifier.fillMaxWidth(), horizontalArrangement = Arrangement.SpaceBetween) {
                            Text("Dispositivo Comodato:", fontWeight = FontWeight.SemiBold)
                            Text("Chromebook ASUS 14\" (Seriale: CB-8821)", fontWeight = FontWeight.Bold, color = Color(0xFF10B981))
                        }
                    }
                }
            }
        }
    }
}

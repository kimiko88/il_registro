package it.scuola.registro.secretary.ui

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
fun SecretaryFascicoloScreen(onBack: () -> Unit = {}) {
    Scaffold(
        topBar = {
            TopAppBar(
                title = { Text("Fascicolo Elettronico Studente", fontWeight = FontWeight.Bold) },
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
                        Text("Dossier Unico e Storico Carriera Scolastica", fontWeight = FontWeight.Bold, fontSize = 16.sp)
                        Text("Anagrafica ministeriale SIDI, storico iscrizioni, certificati vaccinali, intolleranze alimentari e percorsi BES/DSA.", fontSize = 12.sp, color = Color.DarkGray)
                    }
                }
            }

            item {
                Text("Dati Fascicolo: Mario Rossi (Classe 3ª A)", fontWeight = FontWeight.Bold, fontSize = 16.sp)
            }

            item {
                Card(modifier = Modifier.fillMaxWidth(), shape = RoundedCornerShape(12.dp)) {
                    Column(modifier = Modifier.padding(14.dp), verticalArrangement = Arrangement.spacedBy(8.dp)) {
                        Row(modifier = Modifier.fillMaxWidth(), horizontalArrangement = Arrangement.SpaceBetween) {
                            Text("Codice Fiscale:", fontWeight = FontWeight.SemiBold)
                            Text("RSSMRA08A01H501Z", fontWeight = FontWeight.Bold)
                        }
                        Divider()
                        Row(modifier = Modifier.fillMaxWidth(), horizontalArrangement = Arrangement.SpaceBetween) {
                            Text("Codice SIDI Ministeriale:", fontWeight = FontWeight.SemiBold)
                            Text("SIDI-1049281", fontWeight = FontWeight.Bold, color = MaterialTheme.colorScheme.primary)
                        }
                        Divider()
                        Row(modifier = Modifier.fillMaxWidth(), horizontalArrangement = Arrangement.SpaceBetween) {
                            Text("Stato Vaccinale:", fontWeight = FontWeight.SemiBold)
                            Text("Regolare (Conforme L. 119/2017)", fontWeight = FontWeight.Bold, color = Color(0xFF10B981))
                        }
                    }
                }
            }
        }
    }
}

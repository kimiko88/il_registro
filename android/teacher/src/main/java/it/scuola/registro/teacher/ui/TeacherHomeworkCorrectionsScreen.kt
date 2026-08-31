package it.scuola.registro.teacher.ui

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
fun TeacherHomeworkCorrectionsScreen(onBack: () -> Unit = {}) {
    Scaffold(
        topBar = {
            TopAppBar(
                title = { Text("Correzione Compiti Consegnati", fontWeight = FontWeight.Bold) },
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
                        Text("Revisione Elaborati e Compiti Digitali", fontWeight = FontWeight.Bold, fontSize = 16.sp)
                        Text("Visualizzazione compiti consegnati dagli studenti, annotazioni correttive e assegnazione feedback formativo.", fontSize = 12.sp, color = Color.DarkGray)
                    }
                }
            }

            item {
                Text("Elaborati da Valutare - Classe 3ª A", fontWeight = FontWeight.Bold, fontSize = 16.sp)
            }

            item {
                Card(modifier = Modifier.fillMaxWidth(), shape = RoundedCornerShape(12.dp)) {
                    Column(modifier = Modifier.padding(14.dp)) {
                        Row(
                            modifier = Modifier.fillMaxWidth(),
                            horizontalArrangement = Arrangement.SpaceBetween,
                            verticalAlignment = Alignment.CenterVertically
                        ) {
                            Text("Compito: Relazione Laboratorio Fisica", fontWeight = FontWeight.Bold)
                            Surface(color = Color(0xFF3B82F6), shape = RoundedCornerShape(6.dp)) {
                                Text("22 Consegnati / 24", color = Color.White, modifier = Modifier.padding(horizontal = 6.dp, vertical = 2.dp), fontSize = 10.sp)
                            }
                        }
                        Text("Studente: Mario Rossi • File: Relazione_Fisica_Rossi.pdf", fontSize = 12.sp, color = MaterialTheme.colorScheme.primary, fontWeight = FontWeight.SemiBold, modifier = Modifier.padding(vertical = 4.dp))
                        Text("Data Consegna: Ieri alle 18:20 (Puntuale) • Stato: In Attesa di Valutazione", fontSize = 12.sp, color = Color.Gray)
                        Button(
                            onClick = {},
                            shape = RoundedCornerShape(8.dp),
                            modifier = Modifier.fillMaxWidth().padding(top = 6.dp)
                        ) {
                            Icon(Icons.Default.Edit, contentDescription = null, modifier = Modifier.size(16.dp))
                            Spacer(modifier = Modifier.width(6.dp))
                            Text("Apri Elaborato e Inserisci Valutazione")
                        }
                    }
                }
            }
        }
    }
}

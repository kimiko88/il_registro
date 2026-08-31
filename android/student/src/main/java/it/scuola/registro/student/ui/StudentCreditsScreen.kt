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
fun StudentCreditsScreen(onBack: () -> Unit = {}) {
    Scaffold(
        topBar = {
            TopAppBar(
                title = { Text("Credito Scolastico Triennio", fontWeight = FontWeight.Bold) },
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
                        Text("Attribuzione Credito Scolastico (Esame di Stato)", fontWeight = FontWeight.Bold, fontSize = 16.sp)
                        Text("Tabella ministeriale D.Lgs. 62/2017: massimo 40 punti accumulabili nel triennio (12 punti in 3ª, 13 punti in 4ª, 15 punti in 5ª).", fontSize = 12.sp, color = Color.DarkGray)
                    }
                }
            }

            item {
                Text("Riepilogo Crediti Acquisiti", fontWeight = FontWeight.Bold, fontSize = 16.sp)
            }

            item {
                Card(modifier = Modifier.fillMaxWidth(), shape = RoundedCornerShape(12.dp)) {
                    Column(modifier = Modifier.padding(14.dp), verticalArrangement = Arrangement.spacedBy(8.dp)) {
                        Row(modifier = Modifier.fillMaxWidth(), horizontalArrangement = Arrangement.SpaceBetween) {
                            Text("Classe 3ª (Media 7.6):", fontWeight = FontWeight.SemiBold)
                            Text("11 / 12 Punti", fontWeight = FontWeight.Bold, color = MaterialTheme.colorScheme.primary)
                        }
                        Divider()
                        Row(modifier = Modifier.fillMaxWidth(), horizontalArrangement = Arrangement.SpaceBetween) {
                            Text("Classe 4ª (Media 8.2):", fontWeight = FontWeight.SemiBold)
                            Text("12 / 13 Punti", fontWeight = FontWeight.Bold, color = MaterialTheme.colorScheme.primary)
                        }
                        Divider()
                        Row(modifier = Modifier.fillMaxWidth(), horizontalArrangement = Arrangement.SpaceBetween) {
                            Text("Totale Credito Attuale:", fontWeight = FontWeight.Bold)
                            Text("23 Punti", fontWeight = FontWeight.ExtraBold, color = Color(0xFF10B981), fontSize = 16.sp)
                        }
                    }
                }
            }
        }
    }
}

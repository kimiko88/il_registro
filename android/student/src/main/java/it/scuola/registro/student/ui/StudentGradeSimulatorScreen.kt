package it.scuola.registro.student.ui

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
fun StudentGradeSimulatorScreen(onBack: () -> Unit = {}) {
    var simulatedGrade by remember { mutableStateOf("8.0") }

    Scaffold(
        topBar = {
            TopAppBar(
                title = { Text("Simulatore Media Voti", fontWeight = FontWeight.Bold) },
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
                        Text("Calcolo What-If & Previsione Media", fontWeight = FontWeight.Bold, fontSize = 16.sp)
                        Text("Simula l'impatto delle prossime verifiche sulla tua media finale ponderata per ogni disciplina.", fontSize = 12.sp, color = Color.DarkGray)
                    }
                }
            }

            item {
                Text("Simulazione: Matematica (Media Attuale: 7.2)", fontWeight = FontWeight.Bold, fontSize = 16.sp)
            }

            item {
                Card(modifier = Modifier.fillMaxWidth(), shape = RoundedCornerShape(12.dp)) {
                    Column(modifier = Modifier.padding(14.dp), verticalArrangement = Arrangement.spacedBy(10.dp)) {
                        Text("Inserisci voto ipotetico per la prossima prova scritta:")
                        OutlinedTextField(
                            value = simulatedGrade,
                            onValueChange = { simulatedGrade = it },
                            label = { Text("Voto Ipotetico (es. 8.0)") },
                            modifier = Modifier.fillMaxWidth()
                        )

                        Card(
                            modifier = Modifier.fillMaxWidth(),
                            colors = CardDefaults.cardColors(containerColor = Color(0xFFE0F2FE))
                        ) {
                            Row(
                                modifier = Modifier.fillMaxWidth().padding(12.dp),
                                horizontalArrangement = Arrangement.SpaceBetween,
                                verticalAlignment = Alignment.CenterVertically
                            ) {
                                Text("Nuova Media Stimata:", fontWeight = FontWeight.Bold, color = Color(0xFF0369A1))
                                Text("7.45", fontWeight = FontWeight.ExtraBold, fontSize = 18.sp, color = Color(0xFF0284C7))
                            }
                        }
                    }
                }
            }
        }
    }
}

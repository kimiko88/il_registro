package it.scuola.registro.student.ui

import androidx.compose.foundation.layout.*
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.*
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun StudentAccessibilityFeedbackScreen(onBack: () -> Unit = {}) {
    var description by remember { mutableStateOf("") }
    var barrierType by remember { mutableStateOf("Contrasto / Visibilità") }
    var submitted by remember { mutableStateOf(false) }

    Scaffold(
        topBar = {
            TopAppBar(
                title = { Text("Segnalazione Barriere AgID", fontWeight = FontWeight.Bold) },
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
                        Text("Meccanismo di Feedback AgID (Direttiva UE 2016/2102)", fontWeight = FontWeight.Bold, fontSize = 16.sp)
                        Text("Segnala all'RTD scolastico eventuali difficoltà di accesso, incompatibilità con screen reader o problemi di contrasto.", fontSize = 12.sp, color = Color.DarkGray)
                    }
                }
            }

            if (submitted) {
                item {
                    Card(
                        modifier = Modifier.fillMaxWidth(),
                        shape = RoundedCornerShape(12.dp),
                        colors = CardDefaults.cardColors(containerColor = Color(0xFFD1FAE5))
                    ) {
                        Column(modifier = Modifier.padding(16.dp)) {
                            Text("Segnalazione Inviata con Successo!", fontWeight = FontWeight.Bold, color = Color(0xFF065F46))
                            Text("Codice Protocollo Assegnato: A11Y-2026-0528-0912", fontSize = 12.sp, color = Color(0xFF047857), fontWeight = FontWeight.SemiBold, modifier = Modifier.padding(top = 4.dp))
                        }
                    }
                }
            } else {
                item {
                    Text("Dettagli Segnalazione", fontWeight = FontWeight.Bold, fontSize = 16.sp)
                }

                item {
                    Card(modifier = Modifier.fillMaxWidth(), shape = RoundedCornerShape(12.dp)) {
                        Column(modifier = Modifier.padding(14.dp), verticalArrangement = Arrangement.spacedBy(10.dp)) {
                            Text("Tipologia di Barriera Riscontrata:", fontWeight = FontWeight.SemiBold)
                            OutlinedTextField(
                                value = barrierType,
                                onValueChange = { barrierType = it },
                                label = { Text("Categoria") },
                                modifier = Modifier.fillMaxWidth()
                            )

                            OutlinedTextField(
                                value = description,
                                onValueChange = { description = it },
                                label = { Text("Descrizione dettagliata della difficoltà riscontrata") },
                                minLines = 3,
                                modifier = Modifier.fillMaxWidth()
                            )

                            Button(
                                onClick = { submitted = true },
                                shape = RoundedCornerShape(8.dp),
                                modifier = Modifier.fillMaxWidth().padding(top = 4.dp),
                                enabled = description.isNotBlank()
                            ) {
                                Text("Invia Segnalazione all'RTD")
                            }
                        }
                    }
                }
            }
        }
    }
}

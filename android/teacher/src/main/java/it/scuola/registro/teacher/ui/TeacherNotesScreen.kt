package it.scuola.registro.teacher.ui

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
fun TeacherNotesScreen(onBack: () -> Unit = {}) {
    var studentName by remember { mutableStateOf("") }
    var noteText by remember { mutableStateOf("") }
    var noteType by remember { mutableStateOf("Disciplinare") }
    var isSubmitted by remember { mutableStateOf(false) }

    Scaffold(
        topBar = {
            TopAppBar(
                title = { Text("Note & Provvedimenti Disciplinari", fontWeight = FontWeight.Bold) },
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
                        Text("Inserimento Note di Classe o Individuali", fontWeight = FontWeight.Bold, fontSize = 16.sp)
                        Text("Le note disciplinari inviano una notifica push immediata ai genitori per la richiesta di presa visione con firma.", fontSize = 12.sp, color = Color.DarkGray)
                    }
                }
            }

            if (isSubmitted) {
                item {
                    Card(
                        modifier = Modifier.fillMaxWidth(),
                        shape = RoundedCornerShape(12.dp),
                        colors = CardDefaults.cardColors(containerColor = Color(0xFFD1FAE5))
                    ) {
                        Column(modifier = Modifier.padding(16.dp)) {
                            Text("Nota Inserita sul Registro con Successo!", fontWeight = FontWeight.Bold, color = Color(0xFF065F46))
                            Text("Notifica di presa visione inviata alla famiglia.", fontSize = 12.sp, color = Color(0xFF047857), modifier = Modifier.padding(top = 4.dp))
                        }
                    }
                }
            } else {
                item {
                    Text("Nuova Nota Disciplinare", fontWeight = FontWeight.Bold, fontSize = 16.sp)
                }

                item {
                    Card(modifier = Modifier.fillMaxWidth(), shape = RoundedCornerShape(12.dp)) {
                        Column(modifier = Modifier.padding(14.dp), verticalArrangement = Arrangement.spacedBy(10.dp)) {
                            OutlinedTextField(
                                value = studentName,
                                onValueChange = { studentName = it },
                                label = { Text("Studente (es. Mario Rossi)") },
                                modifier = Modifier.fillMaxWidth()
                            )

                            OutlinedTextField(
                                value = noteType,
                                onValueChange = { noteType = it },
                                label = { Text("Tipologia (Disciplinare / Richiamo / Mancanza compiti)") },
                                modifier = Modifier.fillMaxWidth()
                            )

                            OutlinedTextField(
                                value = noteText,
                                onValueChange = { noteText = it },
                                label = { Text("Descrizione del comportamento riscontrato") },
                                minLines = 3,
                                modifier = Modifier.fillMaxWidth()
                            )

                            Button(
                                onClick = { isSubmitted = true },
                                shape = RoundedCornerShape(8.dp),
                                modifier = Modifier.fillMaxWidth().padding(top = 4.dp),
                                enabled = studentName.isNotBlank() && noteText.isNotBlank()
                            ) {
                                Text("Registra Nota sul Giornale di Classe")
                            }
                        }
                    }
                }
            }
        }
    }
}

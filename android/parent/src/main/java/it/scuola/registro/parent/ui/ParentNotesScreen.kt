package it.scuola.registro.parent.ui

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
fun ParentNotesScreen(onBack: () -> Unit = {}) {
    var acknowledged by remember { mutableStateOf(false) }

    Scaffold(
        topBar = {
            TopAppBar(
                title = { Text("Note & Richiami Disciplinari", fontWeight = FontWeight.Bold) },
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
                        Text("Registro Note e Annotazioni Scolastiche", fontWeight = FontWeight.Bold, fontSize = 16.sp)
                        Text("Visualizza le note disciplinari, i richiami e apponi la firma digitale di presa visione.", fontSize = 12.sp, color = Color.DarkGray)
                    }
                }
            }

            item {
                Text("Note sul Registro di Classe", fontWeight = FontWeight.Bold, fontSize = 16.sp)
            }

            item {
                Card(modifier = Modifier.fillMaxWidth(), shape = RoundedCornerShape(12.dp)) {
                    Column(modifier = Modifier.padding(14.dp)) {
                        Row(
                            modifier = Modifier.fillMaxWidth(),
                            horizontalArrangement = Arrangement.SpaceBetween,
                            verticalAlignment = Alignment.CenterVertically
                        ) {
                            Text("Nota Disciplinare Individuale", fontWeight = FontWeight.Bold)
                            Surface(
                                color = if (acknowledged) Color(0xFF10B981) else Color(0xFFEF4444),
                                shape = RoundedCornerShape(6.dp)
                            ) {
                                Text(
                                    if (acknowledged) "Presa Visione Confermata" else "Firma Richiesta",
                                    color = Color.White,
                                    modifier = Modifier.padding(horizontal = 6.dp, vertical = 2.dp),
                                    fontSize = 10.sp
                                )
                            }
                        }
                        Text("Docente: Prof. Rossi • Data: 22 Maggio 2026 • 2ª Ora", fontSize = 12.sp, color = MaterialTheme.colorScheme.primary, fontWeight = FontWeight.SemiBold, modifier = Modifier.padding(vertical = 4.dp))
                        Text("Motivazione: Lo studente disturba ripetutamente lo svolgimento della lezione nonostante i richiami.", fontSize = 12.sp, color = Color.Gray)

                        if (!acknowledged) {
                            Button(
                                onClick = { acknowledged = true },
                                shape = RoundedCornerShape(8.dp),
                                modifier = Modifier.fillMaxWidth().padding(top = 8.dp)
                            ) {
                                Icon(Icons.Default.Check, contentDescription = null, modifier = Modifier.size(16.dp))
                                Spacer(modifier = Modifier.width(6.dp))
                                Text("Apponi Firma di Presa Visione")
                            }
                        }
                    }
                }
            }
        }
    }
}

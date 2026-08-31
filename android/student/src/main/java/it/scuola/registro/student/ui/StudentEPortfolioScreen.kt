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
fun StudentEPortfolioScreen(onBack: () -> Unit = {}) {
    var showDialog by remember { mutableStateOf(false) }
    var title by remember { mutableStateOf("") }
    var year by remember { mutableStateOf("2025/2026") }
    var description by remember { mutableStateOf("") }
    var reflection by remember { mutableStateOf("") }

    Scaffold(
        topBar = {
            TopAppBar(
                title = { Text("E-Portfolio & Capolavori (MIM)", fontWeight = FontWeight.Bold) },
                navigationIcon = {
                    IconButton(onClick = onBack) {
                        Icon(Icons.Default.ArrowBack, contentDescription = "Indietro")
                    }
                }
            )
        },
        floatingActionButton = {
            FloatingActionButton(
                onClick = { showDialog = true },
                containerColor = MaterialTheme.colorScheme.primary,
                contentColor = Color.White
            ) {
                Icon(Icons.Default.Add, contentDescription = "Carica Capolavoro")
            }
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
                    shape = RoundedCornerShape(16.dp),
                    colors = CardDefaults.cardColors(containerColor = MaterialTheme.colorScheme.tertiaryContainer)
                ) {
                    Column(modifier = Modifier.padding(16.dp)) {
                        Text("Orientamento & E-Portfolio", fontWeight = FontWeight.Bold, fontSize = 16.sp)
                        Text("Linee Guida Nazionali per l'Orientamento (D.M. 328/2022)", fontSize = 12.sp, color = Color.DarkGray)
                        Spacer(modifier = Modifier.height(10.dp))
                        Row(
                            modifier = Modifier.fillMaxWidth(),
                            horizontalArrangement = Arrangement.SpaceBetween
                        ) {
                            Text("Ore Orientamento Svolte:", fontSize = 13.sp)
                            Text("30 / 30 h (Target Raggiunto)", fontWeight = FontWeight.Bold, color = Color(0xFF0D9488))
                        }
                    }
                }
            }

            item {
                Text("I Tuoi Capolavori Selezionati", fontWeight = FontWeight.Bold, fontSize = 16.sp)
            }

            item {
                Card(modifier = Modifier.fillMaxWidth(), shape = RoundedCornerShape(12.dp)) {
                    Column(modifier = Modifier.padding(14.dp)) {
                        Row(
                            modifier = Modifier.fillMaxWidth(),
                            horizontalArrangement = Arrangement.SpaceBetween,
                            verticalAlignment = Alignment.CenterVertically
                        ) {
                            Text("Progetto Registro Elettronico Open Source", fontWeight = FontWeight.Bold)
                            Surface(color = Color(0xFF8B5CF6), shape = RoundedCornerShape(6.dp)) {
                                Text("A.S. 2025/2026", color = Color.White, modifier = Modifier.padding(horizontal = 6.dp, vertical = 2.dp), fontSize = 10.sp)
                            }
                        }
                        Text("Sviluppo di un modulo concorrente in Go e applicazione mobile nativa per la scuola pubblica.", fontSize = 13.sp, color = Color.DarkGray, modifier = Modifier.padding(vertical = 4.dp))
                        Surface(color = Color(0xFFF3F4F6), shape = RoundedCornerShape(8.dp), modifier = Modifier.fillMaxWidth().padding(top = 4.dp)) {
                            Text("Autovalutazione: Ho accresciuto la capacità di problem solving e lavoro in team.", fontSize = 12.sp, color = Color.Gray, modifier = Modifier.padding(8.dp))
                        }
                    }
                }
            }
        }
    }

    if (showDialog) {
        AlertDialog(
            onDismissRequest = { showDialog = false },
            title = { Text("Nuovo Capolavoro E-Portfolio") },
            text = {
                Column(verticalArrangement = Arrangement.spacedBy(8.dp)) {
                    OutlinedTextField(
                        value = title,
                        onValueChange = { title = it },
                        label = { Text("Titolo Capolavoro *") },
                        modifier = Modifier.fillMaxWidth()
                    )
                    OutlinedTextField(
                        value = description,
                        onValueChange = { description = it },
                        label = { Text("Descrizione") },
                        modifier = Modifier.fillMaxWidth()
                    )
                    OutlinedTextField(
                        value = reflection,
                        onValueChange = { reflection = it },
                        label = { Text("Riflessione critica") },
                        modifier = Modifier.fillMaxWidth()
                    )
                }
            },
            confirmButton = {
                Button(onClick = { showDialog = false }) {
                    Text("Salva")
                }
            },
            dismissButton = {
                TextButton(onClick = { showDialog = false }) {
                    Text("Annulla")
                }
            }
        )
    }
}

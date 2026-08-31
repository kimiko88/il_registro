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
fun StudentGoalsScreen(onBack: () -> Unit = {}) {
    Scaffold(
        topBar = {
            TopAppBar(
                title = { Text("Obiettivi di Apprendimento", fontWeight = FontWeight.Bold) },
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
                    colors = CardDefaults.cardColors(containerColor = MaterialTheme.colorScheme.secondaryContainer)
                ) {
                    Column(modifier = Modifier.padding(16.dp)) {
                        Text("Traguardi di Miglioramento Personale", fontWeight = FontWeight.Bold, fontSize = 16.sp)
                        Text("Obiettivi concordati con i docenti per il potenziamento delle competenze disciplinari e trasversali.", fontSize = 12.sp, color = Color.DarkGray)
                    }
                }
            }

            item {
                Text("Obiettivi del Quadrimestre", fontWeight = FontWeight.Bold, fontSize = 16.sp)
            }

            item {
                Card(modifier = Modifier.fillMaxWidth(), shape = RoundedCornerShape(12.dp)) {
                    Column(modifier = Modifier.padding(14.dp)) {
                        Row(
                            modifier = Modifier.fillMaxWidth(),
                            horizontalArrangement = Arrangement.SpaceBetween,
                            verticalAlignment = Alignment.CenterVertically
                        ) {
                            Text("Consolidamento Studio di Funzione", fontWeight = FontWeight.Bold)
                            Surface(color = Color(0xFF3B82F6), shape = RoundedCornerShape(6.dp)) {
                                Text("In Corso (Avanzamento 75%)", color = Color.White, modifier = Modifier.padding(horizontal = 6.dp, vertical = 2.dp), fontSize = 10.sp)
                            }
                        }
                        Text("Materia: Matematica • Target: Raggiungere media 7.0", fontSize = 12.sp, color = MaterialTheme.colorScheme.primary, fontWeight = FontWeight.SemiBold, modifier = Modifier.padding(vertical = 4.dp))
                        LinearProgressIndicator(
                            progress = 0.75f,
                            modifier = Modifier.fillMaxWidth().padding(vertical = 8.dp),
                            color = MaterialTheme.colorScheme.primary
                        )
                        Text("Feedback Docente: Ottimo impegno nei compiti a casa e nei laboratori", fontSize = 12.sp, color = Color.Gray)
                    }
                }
            }
        }
    }
}

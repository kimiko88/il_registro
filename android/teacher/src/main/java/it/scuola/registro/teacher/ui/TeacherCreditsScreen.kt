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
fun TeacherCreditsScreen(onBack: () -> Unit = {}) {
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
                        Text("Calcolo e Delibera Credito Scolastico (D.Lgs. 62/2017)", fontWeight = FontWeight.Bold, fontSize = 16.sp)
                        Text("Assegnazione automatica della banda di oscillazione in base alla media dei voti dello scrutinio finale e delle attività integrative.", fontSize = 12.sp, color = Color.DarkGray)
                    }
                }
            }

            item {
                Text("Proposte di Credito per la Classe 3ª A", fontWeight = FontWeight.Bold, fontSize = 16.sp)
            }

            item {
                Card(modifier = Modifier.fillMaxWidth(), shape = RoundedCornerShape(12.dp)) {
                    Column(modifier = Modifier.padding(14.dp)) {
                        Row(
                            modifier = Modifier.fillMaxWidth(),
                            horizontalArrangement = Arrangement.SpaceBetween,
                            verticalAlignment = Alignment.CenterVertically
                        ) {
                            Text("Mario Rossi • Media Finale: 7.60", fontWeight = FontWeight.Bold)
                            Surface(color = Color(0xFF10B981), shape = RoundedCornerShape(6.dp)) {
                                Text("Banda: 10 - 11 Punti", color = Color.White, modifier = Modifier.padding(horizontal = 6.dp, vertical = 2.dp), fontSize = 10.sp)
                            }
                        }
                        Text("Proposta Consiglio: 11 Punti (Fascia Superiore per crediti formativi e PCTO)", fontSize = 12.sp, color = MaterialTheme.colorScheme.primary, fontWeight = FontWeight.SemiBold, modifier = Modifier.padding(vertical = 4.dp))
                    }
                }
            }
        }
    }
}

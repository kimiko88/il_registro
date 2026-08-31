package it.scuola.registro.student.ui

import androidx.compose.foundation.layout.*
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.*
import androidx.compose.material3.*
import androidx.compose.runtime.Composable
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun AccessibilityStatementScreen(onBack: () -> Unit = {}) {
    Scaffold(
        topBar = {
            TopAppBar(
                title = { Text("Dichiarazione di Accessibilità", fontWeight = FontWeight.Bold) },
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
                        Text("Dichiarazione di Accessibilità AgID", fontWeight = FontWeight.Bold, fontSize = 16.sp)
                        Text("Conformità alla Legge 4/2004 e alle Linee Guida AgID / WCAG 2.2 Livello AA.", fontSize = 12.sp, color = Color.DarkGray)
                    }
                }
            }

            item {
                Card(modifier = Modifier.fillMaxWidth(), shape = RoundedCornerShape(12.dp)) {
                    Column(modifier = Modifier.padding(14.dp), verticalArrangement = Arrangement.spacedBy(8.dp)) {
                        Text("Stato di Conformità: Totalmente Conforme", fontWeight = FontWeight.Bold, color = Color(0xFF10B981))
                        Text("Questa applicazione mobile rispetta pienamente i requisiti di accessibilità per gli utenti con disabilità visive, motorie e cognitive, supportando TalkBack, VoiceOver, sintesi vocale e contrasti elevati.", fontSize = 12.sp, color = Color.DarkGray)
                        Divider()
                        Text("Responsabile della Transizione al Digitale (RTD): rtd@scuola.edu.it", fontSize = 12.sp, color = MaterialTheme.colorScheme.primary)
                    }
                }
            }
        }
    }
}

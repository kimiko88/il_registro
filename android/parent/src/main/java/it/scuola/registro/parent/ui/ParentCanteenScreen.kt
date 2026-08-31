package it.scuola.registro.parent.ui

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
fun ParentCanteenScreen(onBack: () -> Unit = {}) {
    Scaffold(
        topBar = {
            TopAppBar(
                title = { Text("Mensa & Diete Speciali", fontWeight = FontWeight.Bold) },
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
                        Text("Servizio Ristorazione Scolastica", fontWeight = FontWeight.Bold, fontSize = 16.sp)
                        Text("Consulta il menu settimanale, monitora la presenza al pasto e gestisci le certificazioni per diete etico-religiose o sanitarie.", fontSize = 12.sp, color = Color.DarkGray)
                    }
                }
            }

            item {
                Text("Menu del Giorno - Mercoledì", fontWeight = FontWeight.Bold, fontSize = 16.sp)
            }

            item {
                Card(modifier = Modifier.fillMaxWidth(), shape = RoundedCornerShape(12.dp)) {
                    Column(modifier = Modifier.padding(14.dp), verticalArrangement = Arrangement.spacedBy(6.dp)) {
                        Row(
                            modifier = Modifier.fillMaxWidth(),
                            horizontalArrangement = Arrangement.SpaceBetween,
                            verticalAlignment = Alignment.CenterVertically
                        ) {
                            Text("Menu Primaverile (Settimana 3)", fontWeight = FontWeight.Bold)
                            Surface(color = Color(0xFF10B981), shape = RoundedCornerShape(6.dp)) {
                                Text("Pasto Consumato", color = Color.White, modifier = Modifier.padding(horizontal = 6.dp, vertical = 2.dp), fontSize = 10.sp)
                            }
                        }
                        Text("Primo: Risotto con verdure di stagione BIO", fontSize = 12.sp, color = MaterialTheme.colorScheme.primary, fontWeight = FontWeight.SemiBold)
                        Text("Secondo: Filetto di platessa al forno con patate", fontSize = 12.sp, color = MaterialTheme.colorScheme.primary, fontWeight = FontWeight.SemiBold)
                        Text("Frutta: Mela biologica italiana • Acqua microfiltrata", fontSize = 12.sp, color = Color.Gray)
                    }
                }
            }
        }
    }
}

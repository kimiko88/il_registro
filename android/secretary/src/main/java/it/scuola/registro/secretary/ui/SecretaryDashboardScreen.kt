package it.scuola.registro.secretary.ui

import androidx.compose.foundation.layout.*
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.items
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.*
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.res.stringResource
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import it.scuola.registro.secretary.R

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun SecretaryDashboardScreen(onLogout: () -> Unit = {}) {
    var selectedTab by remember { mutableStateOf(0) }

    Scaffold(
        topBar = {
            TopAppBar(
                title = {
                    Text(
                        when (selectedTab) {
                            0 -> stringResource(R.string.secretary_dashboard_title)
                            1 -> stringResource(R.string.user_management)
                            2 -> stringResource(R.string.scrutiny_supervision)
                            3 -> stringResource(R.string.generate_certificates)
                            else -> stringResource(R.string.audit_logs)
                        },
                        fontWeight = FontWeight.Bold
                    )
                },
                actions = {
                    IconButton(onClick = onLogout) {
                        Icon(Icons.Default.Logout, contentDescription = "Esci", tint = Color.White)
                    }
                },
                colors = TopAppBarDefaults.topAppBarColors(
                    containerColor = Color(0xFF581C87),
                    titleContentColor = Color.White
                )
            )
        },
        bottomBar = {
            NavigationBar(containerColor = MaterialTheme.colorScheme.surface, tonalElevation = 8.dp) {
                NavigationBarItem(
                    selected = selectedTab == 0,
                    onClick = { selectedTab = 0 },
                    icon = { Icon(Icons.Default.Dashboard, contentDescription = "Pannello") },
                    label = { Text("Pannello", fontSize = 11.sp) }
                )
                NavigationBarItem(
                    selected = selectedTab == 1,
                    onClick = { selectedTab = 1 },
                    icon = { Icon(Icons.Default.People, contentDescription = "Utenti") },
                    label = { Text("Utenti", fontSize = 11.sp) }
                )
                NavigationBarItem(
                    selected = selectedTab == 2,
                    onClick = { selectedTab = 2 },
                    icon = { Icon(Icons.Default.Class, contentDescription = "Scrutini") },
                    label = { Text("Scrutini", fontSize = 11.sp) }
                )
                NavigationBarItem(
                    selected = selectedTab == 3,
                    onClick = { selectedTab = 3 },
                    icon = { Icon(Icons.Default.Description, contentDescription = "Certificati") },
                    label = { Text("Certificati", fontSize = 11.sp) }
                )
                NavigationBarItem(
                    selected = selectedTab == 4,
                    onClick = { selectedTab = 4 },
                    icon = { Icon(Icons.Default.Security, contentDescription = "Audit") },
                    label = { Text("Audit", fontSize = 11.sp) }
                )
            }
        }
    ) { padding ->
        Box(modifier = Modifier.fillMaxSize().padding(padding).padding(16.dp)) {
            when (selectedTab) {
                0 -> SecretaryOverviewTab()
                1 -> SecretaryUsersTab()
                2 -> SecretaryScrutinyTab()
                3 -> SecretaryCertificatesTab()
                4 -> SecretaryAuditTab()
            }
        }
    }
}

@Composable
fun SecretaryOverviewTab() {
    LazyColumn(verticalArrangement = Arrangement.spacedBy(14.dp)) {
        item {
            Row(horizontalArrangement = Arrangement.spacedBy(12.dp)) {
                MetricCard(modifier = Modifier.weight(1f), title = stringResource(R.string.total_students), value = "1,248", icon = Icons.Default.People)
                MetricCard(modifier = Modifier.weight(1f), title = stringResource(R.string.total_teachers), value = "94", icon = Icons.Default.School)
            }
        }
        item {
            Card(modifier = Modifier.fillMaxWidth(), shape = RoundedCornerShape(16.dp)) {
                Column(modifier = Modifier.padding(16.dp)) {
                    Text("Stato Sistema & Operazioni", fontWeight = FontWeight.Bold, fontSize = 16.sp)
                    Spacer(modifier = Modifier.height(8.dp))
                    Text("Anno Scolastico Attivo: 2025/2026", fontSize = 13.sp, color = Color.Gray)
                    Text("Server API & Database: Operativi (Latenza 14ms)", fontSize = 13.sp, color = Color(0xFF10B981))
                }
            }
        }
    }
}

@Composable
fun SecretaryUsersTab() {
    val sampleUsers = listOf(
        "Prof.ssa Maria Rossi" to "Docente • Lettere",
        "Prof. Marco Bianchi" to "Docente • Matematica",
        "Mario Rossi (2B)" to "Studente",
        "Giuseppe Rossi" to "Genitore"
    )
    LazyColumn(verticalArrangement = Arrangement.spacedBy(12.dp)) {
        item {
            Row(modifier = Modifier.fillMaxWidth(), horizontalArrangement = Arrangement.SpaceBetween) {
                Text("Anagrafica Utenti Sistema", fontWeight = FontWeight.Bold, fontSize = 16.sp)
                Button(onClick = { }, shape = RoundedCornerShape(8.dp)) {
                    Icon(Icons.Default.Add, contentDescription = null, modifier = Modifier.size(16.dp))
                    Spacer(modifier = Modifier.width(4.dp))
                    Text("Nuovo Utente", fontSize = 12.sp)
                }
            }
        }
        items(sampleUsers) { (name, role) ->
            Card(modifier = Modifier.fillMaxWidth(), shape = RoundedCornerShape(12.dp)) {
                Row(
                    modifier = Modifier.fillMaxWidth().padding(14.dp),
                    horizontalArrangement = Arrangement.SpaceBetween,
                    verticalAlignment = Alignment.CenterVertically
                ) {
                    Column {
                        Text(name, fontWeight = FontWeight.Bold)
                        Text(role, fontSize = 12.sp, color = Color.Gray)
                    }
                    IconButton(onClick = { }) {
                        Icon(Icons.Default.Edit, contentDescription = null, tint = Color(0xFF581C87))
                    }
                }
            }
        }
    }
}

@Composable
fun SecretaryScrutinyTab() {
    val classes = listOf("1A" to "Completato", "2A" to "In corso", "3A" to "Completato", "4B" to "Differito")
    LazyColumn(verticalArrangement = Arrangement.spacedBy(12.dp)) {
        item {
            Text("Stato Scrutini di Tutte le Classi", fontWeight = FontWeight.Bold, fontSize = 16.sp)
        }
        items(classes) { (cls, status) ->
            Card(modifier = Modifier.fillMaxWidth(), shape = RoundedCornerShape(12.dp)) {
                Row(
                    modifier = Modifier.fillMaxWidth().padding(14.dp),
                    horizontalArrangement = Arrangement.SpaceBetween,
                    verticalAlignment = Alignment.CenterVertically
                ) {
                    Text("Classe $cls", fontWeight = FontWeight.Bold)
                    Surface(
                        color = if (status == "Completato") Color(0xFF10B981) else Color(0xFFEA580C),
                        shape = RoundedCornerShape(6.dp)
                    ) {
                        Text(status, color = Color.White, modifier = Modifier.padding(horizontal = 8.dp, vertical = 4.dp), fontSize = 12.sp)
                    }
                }
            }
        }
    }
}

@Composable
fun SecretaryCertificatesTab() {
    val certTypes = listOf(
        "Certificato di Iscrizione e Frequenza",
        "Certificato con Voti e Valutazioni",
        "Certificato di Diploma / Esame di Stato"
    )
    LazyColumn(verticalArrangement = Arrangement.spacedBy(12.dp)) {
        item {
            Text("Emissione Certificati Ufficiali", fontWeight = FontWeight.Bold, fontSize = 16.sp)
        }
        items(certTypes) { cert ->
            Card(modifier = Modifier.fillMaxWidth(), shape = RoundedCornerShape(12.dp)) {
                Row(
                    modifier = Modifier.fillMaxWidth().padding(14.dp),
                    horizontalArrangement = Arrangement.SpaceBetween,
                    verticalAlignment = Alignment.CenterVertically
                ) {
                    Text(cert, fontWeight = FontWeight.Medium, fontSize = 13.sp)
                    Button(onClick = { }, shape = RoundedCornerShape(8.dp)) {
                        Icon(Icons.Default.Download, contentDescription = null, modifier = Modifier.size(16.dp))
                        Spacer(modifier = Modifier.width(4.dp))
                        Text("PDF", fontSize = 12.sp)
                    }
                }
            }
        }
    }
}

@Composable
fun SecretaryAuditTab() {
    val logs = listOf(
        "31 Ago 11:20 • Login SuperAdmin da IP sicuro",
        "31 Ago 10:45 • Generazione certificato iscrizione matricola 412",
        "31 Ago 09:30 • Chiusura scrutinio classe 3A"
    )
    LazyColumn(verticalArrangement = Arrangement.spacedBy(10.dp)) {
        item {
            Text("Registro Audit & Eventi di Sistema", fontWeight = FontWeight.Bold, fontSize = 16.sp)
        }
        items(logs) { log ->
            Card(modifier = Modifier.fillMaxWidth(), shape = RoundedCornerShape(10.dp)) {
                Row(modifier = Modifier.fillMaxWidth().padding(12.dp), verticalAlignment = Alignment.CenterVertically) {
                    Icon(Icons.Default.Security, contentDescription = null, tint = Color(0xFF581C87))
                    Spacer(modifier = Modifier.width(8.dp))
                    Text(log, fontSize = 12.sp)
                }
            }
        }
    }
}

@Composable
fun MetricCard(modifier: Modifier = Modifier, title: String, value: String, icon: androidx.compose.ui.graphics.vector.ImageVector) {
    Card(
        modifier = modifier,
        shape = RoundedCornerShape(14.dp),
        colors = CardDefaults.cardColors(containerColor = MaterialTheme.colorScheme.surfaceVariant)
    ) {
        Column(modifier = Modifier.padding(16.dp)) {
            Icon(icon, contentDescription = null, tint = Color(0xFF581C87))
            Spacer(modifier = Modifier.height(8.dp))
            Text(value, fontSize = 22.sp, fontWeight = FontWeight.Bold)
            Text(title, fontSize = 12.sp, color = Color.Gray)
        }
    }
}

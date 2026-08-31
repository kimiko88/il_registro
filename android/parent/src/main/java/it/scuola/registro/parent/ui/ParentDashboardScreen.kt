package it.scuola.registro.parent.ui

import androidx.compose.foundation.background
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
import it.scuola.registro.parent.R

data class ChildItem(val id: String, val name: String, val className: String)
data class AbsenceToJustify(val date: String, val reason: String, var isJustified: Boolean = false)
data class ColloquioSlot(val teacher: String, val subject: String, val slotTime: String, val isBooked: Boolean)

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun ParentDashboardScreen(onLogout: () -> Unit = {}) {
    val children = remember {
        listOf(
            ChildItem("1", "Marco Rossi", "Classe 2A"),
            ChildItem("2", "Giulia Rossi", "Classe 4B")
        )
    }
    var selectedChild by remember { mutableStateOf(children.first()) }
    var selectedTab by remember { mutableStateOf(0) }

    Scaffold(
        topBar = {
            TopAppBar(
                title = {
                    Text(
                        when (selectedTab) {
                            0 -> stringResource(R.string.parent_dashboard_title)
                            1 -> stringResource(R.string.child_grades)
                            2 -> stringResource(R.string.child_attendance)
                            3 -> stringResource(R.string.upcoming_colloqui)
                            else -> "Avvisi & Circolari"
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
                    containerColor = Color(0xFF0F172A),
                    titleContentColor = Color.White
                )
            )
        },
        bottomBar = {
            NavigationBar(containerColor = MaterialTheme.colorScheme.surface, tonalElevation = 8.dp) {
                NavigationBarItem(
                    selected = selectedTab == 0,
                    onClick = { selectedTab = 0 },
                    icon = { Icon(Icons.Default.FamilyRestroom, contentDescription = "Figli") },
                    label = { Text("Figli", fontSize = 11.sp) }
                )
                NavigationBarItem(
                    selected = selectedTab == 1,
                    onClick = { selectedTab = 1 },
                    icon = { Icon(Icons.Default.Grade, contentDescription = "Voti") },
                    label = { Text("Voti", fontSize = 11.sp) }
                )
                NavigationBarItem(
                    selected = selectedTab == 2,
                    onClick = { selectedTab = 2 },
                    icon = { Icon(Icons.Default.CheckCircle, contentDescription = "Giustifiche") },
                    label = { Text("Giustifiche", fontSize = 11.sp) }
                )
                NavigationBarItem(
                    selected = selectedTab == 3,
                    onClick = { selectedTab = 3 },
                    icon = { Icon(Icons.Default.CalendarMonth, contentDescription = "Colloqui") },
                    label = { Text("Colloqui", fontSize = 11.sp) }
                )
                NavigationBarItem(
                    selected = selectedTab == 4,
                    onClick = { selectedTab = 4 },
                    icon = { Icon(Icons.Default.Campaign, contentDescription = "Avvisi") },
                    label = { Text("Avvisi", fontSize = 11.sp) }
                )
            }
        }
    ) { padding ->
        Column(
            modifier = Modifier
                .fillMaxSize()
                .padding(padding)
        ) {
            // Child Selector Bar (Persistent across tabs)
            Card(
                modifier = Modifier
                    .fillMaxWidth()
                    .padding(12.dp),
                shape = RoundedCornerShape(12.dp),
                colors = CardDefaults.cardColors(containerColor = MaterialTheme.colorScheme.surfaceVariant)
            ) {
                Row(
                    modifier = Modifier
                        .fillMaxWidth()
                        .padding(horizontal = 14.dp, vertical = 10.dp),
                    horizontalArrangement = Arrangement.SpaceBetween,
                    verticalAlignment = Alignment.CenterVertically
                ) {
                    Text("Figlio attivo:", fontSize = 13.sp, color = Color.Gray)
                    Row(horizontalArrangement = Arrangement.spacedBy(8.dp)) {
                        children.forEach { child ->
                            FilterChip(
                                selected = selectedChild.id == child.id,
                                onClick = { selectedChild = child },
                                label = { Text(child.name, fontSize = 12.sp) }
                            )
                        }
                    }
                }
            }

            Box(modifier = Modifier.fillMaxSize().padding(horizontal = 16.dp)) {
                when (selectedTab) {
                    0 -> ParentOverviewTab(selectedChild)
                    1 -> ParentGradesTab(selectedChild)
                    2 -> ParentJustificationsTab(selectedChild)
                    3 -> ParentColloquiTab()
                    4 -> ParentCommunicationsTab()
                }
            }
        }
    }
}

@Composable
fun ParentOverviewTab(child: ChildItem) {
    LazyColumn(verticalArrangement = Arrangement.spacedBy(14.dp)) {
        item {
            Card(
                modifier = Modifier.fillMaxWidth(),
                shape = RoundedCornerShape(16.dp),
                colors = CardDefaults.cardColors(containerColor = Color(0xFF0F172A))
            ) {
                Column(modifier = Modifier.padding(18.dp)) {
                    Text(child.name, color = Color.White, fontWeight = FontWeight.Bold, fontSize = 20.sp)
                    Text(child.className, color = Color.White.copy(alpha = 0.8f), fontSize = 14.sp)
                    Spacer(modifier = Modifier.height(12.dp))
                    Row(horizontalArrangement = Arrangement.spacedBy(16.dp)) {
                        Surface(color = Color.White.copy(alpha = 0.2f), shape = RoundedCornerShape(8.dp)) {
                            Text("Media: 8.1", color = Color.White, modifier = Modifier.padding(6.dp), fontWeight = FontWeight.Bold)
                        }
                        Surface(color = Color.White.copy(alpha = 0.2f), shape = RoundedCornerShape(8.dp)) {
                            Text("Assenze: 2", color = Color.White, modifier = Modifier.padding(6.dp), fontWeight = FontWeight.Bold)
                        }
                    }
                }
            }
        }
    }
}

@Composable
fun ParentGradesTab(child: ChildItem) {
    val sampleGrades = listOf(
        "Matematica" to "8.5 (Scritto - 30 Ago)",
        "Italiano" to "8.0 (Orale - 28 Ago)",
        "Fisica" to "7.5 (Scritto - 22 Ago)",
        "Inglese" to "9.0 (Pratico - 20 Ago)"
    )
    LazyColumn(verticalArrangement = Arrangement.spacedBy(12.dp)) {
        items(sampleGrades) { (sub, grade) ->
            Card(modifier = Modifier.fillMaxWidth(), shape = RoundedCornerShape(12.dp)) {
                Row(
                    modifier = Modifier.fillMaxWidth().padding(14.dp),
                    horizontalArrangement = Arrangement.SpaceBetween,
                    verticalAlignment = Alignment.CenterVertically
                ) {
                    Text(sub, fontWeight = FontWeight.Bold)
                    Text(grade, color = Color(0xFF0D9488), fontWeight = FontWeight.SemiBold)
                }
            }
        }
    }
}

@Composable
fun ParentJustificationsTab(child: ChildItem) {
    val absences = remember {
        mutableStateListOf(
            AbsenceToJustify("26 Ago 2026", "Assenza per motivi di salute"),
            AbsenceToJustify("18 Ago 2026", "Ritardo 10 minuti 1a ora")
        )
    }

    LazyColumn(verticalArrangement = Arrangement.spacedBy(12.dp)) {
        item {
            Text("Assenze e Ritardi da Giustificare", fontWeight = FontWeight.Bold, fontSize = 16.sp)
        }
        items(absences) { item ->
            Card(modifier = Modifier.fillMaxWidth(), shape = RoundedCornerShape(12.dp)) {
                Row(
                    modifier = Modifier.fillMaxWidth().padding(14.dp),
                    horizontalArrangement = Arrangement.SpaceBetween,
                    verticalAlignment = Alignment.CenterVertically
                ) {
                    Column {
                        Text(item.date, fontWeight = FontWeight.Bold)
                        Text(item.reason, fontSize = 12.sp, color = Color.Gray)
                    }
                    if (item.isJustified) {
                        Surface(color = Color(0xFF10B981), shape = RoundedCornerShape(6.dp)) {
                            Text("Giustificata", color = Color.White, modifier = Modifier.padding(horizontal = 8.dp, vertical = 4.dp), fontSize = 12.sp)
                        }
                    } else {
                        Button(onClick = { item.isJustified = true }, shape = RoundedCornerShape(8.dp)) {
                            Text("Giustifica", fontSize = 12.sp)
                        }
                    }
                }
            }
        }
    }
}

@Composable
fun ParentColloquiTab() {
    val slots = listOf(
        ColloquioSlot("Prof. Bianchi", "Matematica", "Giovedì 15:30 - 15:45", false),
        ColloquioSlot("Prof.ssa Rossi", "Italiano", "Venerdì 16:00 - 16:15", true)
    )
    LazyColumn(verticalArrangement = Arrangement.spacedBy(12.dp)) {
        items(slots) { slot ->
            Card(modifier = Modifier.fillMaxWidth(), shape = RoundedCornerShape(12.dp)) {
                Row(
                    modifier = Modifier.fillMaxWidth().padding(14.dp),
                    horizontalArrangement = Arrangement.SpaceBetween,
                    verticalAlignment = Alignment.CenterVertically
                ) {
                    Column {
                        Text(slot.teacher, fontWeight = FontWeight.Bold)
                        Text("${slot.subject} • ${slot.slotTime}", fontSize = 12.sp, color = Color.Gray)
                    }
                    Button(
                        onClick = { },
                        colors = ButtonDefaults.buttonColors(containerColor = if (slot.isBooked) Color(0xFF10B981) else Color(0xFF0D9488)),
                        shape = RoundedCornerShape(8.dp)
                    ) {
                        Text(if (slot.isBooked) "Prenotato" else "Prenota", fontSize = 12.sp)
                    }
                }
            }
        }
    }
}

@Composable
fun ParentCommunicationsTab() {
    val circulars = listOf(
        "Circolare n. 42: Calendario incontri scuola-famiglia di settembre",
        "Circolare n. 41: Protocollo uscite didattiche e viaggi di istruzione",
        "Circolare n. 40: Attivazione registro elettronico e credenziali"
    )
    LazyColumn(verticalArrangement = Arrangement.spacedBy(12.dp)) {
        items(circulars) { c ->
            Card(modifier = Modifier.fillMaxWidth(), shape = RoundedCornerShape(12.dp)) {
                Row(modifier = Modifier.fillMaxWidth().padding(14.dp), verticalAlignment = Alignment.CenterVertically) {
                    Icon(Icons.Default.Description, contentDescription = null, tint = Color(0xFF0F172A))
                    Spacer(modifier = Modifier.width(10.dp))
                    Text(c, fontSize = 13.sp, fontWeight = FontWeight.Medium)
                }
            }
        }
    }
}

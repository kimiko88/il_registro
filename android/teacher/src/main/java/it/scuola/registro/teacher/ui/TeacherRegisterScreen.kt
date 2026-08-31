package it.scuola.registro.teacher.ui

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
import it.scuola.registro.teacher.R

data class StudentAttendance(val name: String, var isPresent: Boolean, var isLate: Boolean, var grade: String = "")

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun TeacherRegisterScreen(onLogout: () -> Unit = {}) {
    var selectedTab by remember { mutableStateOf(0) }
    val students = remember {
        mutableStateListOf(
            StudentAttendance("Banchi Andrea", isPresent = true, isLate = false, grade = "8"),
            StudentAttendance("Bianchi Elena", isPresent = true, isLate = false, grade = "7.5"),
            StudentAttendance("Ferrari Matteo", isPresent = false, isLate = false, grade = ""),
            StudentAttendance("Rossi Sofia", isPresent = true, isLate = true, grade = "8.5")
        )
    }

    Scaffold(
        topBar = {
            TopAppBar(
                title = {
                    Text(
                        when (selectedTab) {
                            0 -> stringResource(R.string.teacher_dashboard_title)
                            1 -> stringResource(R.string.take_attendance)
                            2 -> "Inserimento Voti"
                            3 -> "Agenda di Classe"
                            else -> "Tabellone Scrutini"
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
                    containerColor = Color(0xFF1E293B),
                    titleContentColor = Color.White
                )
            )
        },
        bottomBar = {
            NavigationBar(containerColor = MaterialTheme.colorScheme.surface, tonalElevation = 8.dp) {
                NavigationBarItem(
                    selected = selectedTab == 0,
                    onClick = { selectedTab = 0 },
                    icon = { Icon(Icons.Default.MenuBook, contentDescription = "Registro") },
                    label = { Text("Registro", fontSize = 11.sp) }
                )
                NavigationBarItem(
                    selected = selectedTab == 1,
                    onClick = { selectedTab = 1 },
                    icon = { Icon(Icons.Default.Checklist, contentDescription = "Appello") },
                    label = { Text("Appello", fontSize = 11.sp) }
                )
                NavigationBarItem(
                    selected = selectedTab == 2,
                    onClick = { selectedTab = 2 },
                    icon = { Icon(Icons.Default.Grade, contentDescription = "Voti") },
                    label = { Text("Voti", fontSize = 11.sp) }
                )
                NavigationBarItem(
                    selected = selectedTab == 3,
                    onClick = { selectedTab = 3 },
                    icon = { Icon(Icons.Default.EditCalendar, contentDescription = "Agenda") },
                    label = { Text("Agenda", fontSize = 11.sp) }
                )
                NavigationBarItem(
                    selected = selectedTab == 4,
                    onClick = { selectedTab = 4 },
                    icon = { Icon(Icons.Default.Gavel, contentDescription = "Scrutini") },
                    label = { Text("Scrutini", fontSize = 11.sp) }
                )
            }
        }
    ) { padding ->
        Box(modifier = Modifier.fillMaxSize().padding(padding).padding(16.dp)) {
            when (selectedTab) {
                0 -> TeacherFirmaTab()
                1 -> TeacherAppelloTab(students)
                2 -> TeacherGradesEntryTab(students)
                3 -> TeacherAgendaTab()
                4 -> TeacherScrutinyTab()
            }
        }
    }
}

@Composable
fun TeacherFirmaTab() {
    var topic by remember { mutableStateOf("Equazioni e disequazioni esponenziali") }
    var isSigned by remember { mutableStateOf(true) }

    LazyColumn(verticalArrangement = Arrangement.spacedBy(14.dp)) {
        item {
            Card(
                modifier = Modifier.fillMaxWidth(),
                shape = RoundedCornerShape(16.dp),
                colors = CardDefaults.cardColors(containerColor = MaterialTheme.colorScheme.surfaceVariant)
            ) {
                Column(modifier = Modifier.padding(16.dp)) {
                    Text("Firma Registro di Classe", fontWeight = FontWeight.Bold, fontSize = 18.sp)
                    Spacer(modifier = Modifier.height(4.dp))
                    Text("Classe 3A • Matematica (1a e 2a ora)", fontSize = 13.sp, color = Color.Gray)
                    Spacer(modifier = Modifier.height(12.dp))
                    OutlinedTextField(
                        value = topic,
                        onValueChange = { topic = it },
                        label = { Text("Argomento Lezione") },
                        modifier = Modifier.fillMaxWidth()
                    )
                    Spacer(modifier = Modifier.height(12.dp))
                    Button(
                        onClick = { isSigned = true },
                        shape = RoundedCornerShape(10.dp),
                        colors = ButtonDefaults.buttonColors(containerColor = if (isSigned) Color(0xFF10B981) else Color(0xFF2563EB)),
                        modifier = Modifier.fillMaxWidth()
                    ) {
                        Icon(Icons.Default.Edit, contentDescription = null, modifier = Modifier.size(18.dp))
                        Spacer(modifier = Modifier.width(8.dp))
                        Text(if (isSigned) "Ora Lezione Firmata" else "Firma Ora di Lezione")
                    }
                }
            }
        }
    }
}

@Composable
fun TeacherAppelloTab(students: List<StudentAttendance>) {
    LazyColumn(verticalArrangement = Arrangement.spacedBy(12.dp)) {
        item {
            Text("Rilevazione Presenze (Appello Giornaliero)", fontWeight = FontWeight.Bold, fontSize = 16.sp)
        }
        items(students) { student ->
            Card(modifier = Modifier.fillMaxWidth(), shape = RoundedCornerShape(12.dp)) {
                Row(
                    modifier = Modifier.fillMaxWidth().padding(14.dp),
                    horizontalArrangement = Arrangement.SpaceBetween,
                    verticalAlignment = Alignment.CenterVertically
                ) {
                    Text(text = student.name, fontWeight = FontWeight.SemiBold, fontSize = 15.sp)
                    Row(horizontalArrangement = Arrangement.spacedBy(6.dp)) {
                        FilterChip(
                            selected = student.isPresent,
                            onClick = { student.isPresent = true; student.isLate = false },
                            label = { Text("P") }
                        )
                        FilterChip(
                            selected = !student.isPresent,
                            onClick = { student.isPresent = false; student.isLate = false },
                            label = { Text("A") }
                        )
                        FilterChip(
                            selected = student.isLate,
                            onClick = { student.isPresent = true; student.isLate = true },
                            label = { Text("R") }
                        )
                    }
                }
            }
        }
    }
}

@Composable
fun TeacherGradesEntryTab(students: List<StudentAttendance>) {
    LazyColumn(verticalArrangement = Arrangement.spacedBy(12.dp)) {
        item {
            Text("Inserimento Valutazioni & Voti Pesati", fontWeight = FontWeight.Bold, fontSize = 16.sp)
        }
        items(students) { student ->
            Card(modifier = Modifier.fillMaxWidth(), shape = RoundedCornerShape(12.dp)) {
                Row(
                    modifier = Modifier.fillMaxWidth().padding(14.dp),
                    horizontalArrangement = Arrangement.SpaceBetween,
                    verticalAlignment = Alignment.CenterVertically
                ) {
                    Text(student.name, fontWeight = FontWeight.SemiBold)
                    Row(verticalAlignment = Alignment.CenterVertically, horizontalArrangement = Arrangement.spacedBy(8.dp)) {
                        Text(if (student.grade.isNotEmpty()) "Voto: ${student.grade}" else "N/A", fontWeight = FontWeight.Bold, color = Color(0xFF2563EB))
                        Button(onClick = { }, shape = RoundedCornerShape(6.dp)) {
                            Text("Assegna", fontSize = 12.sp)
                        }
                    }
                }
            }
        }
    }
}

@Composable
fun TeacherAgendaTab() {
    LazyColumn(verticalArrangement = Arrangement.spacedBy(12.dp)) {
        item {
            Card(modifier = Modifier.fillMaxWidth(), shape = RoundedCornerShape(14.dp)) {
                Column(modifier = Modifier.padding(16.dp)) {
                    Text("Assegna Compito / Verifica", fontWeight = FontWeight.Bold, fontSize = 16.sp)
                    Spacer(modifier = Modifier.height(8.dp))
                    Button(onClick = { }, modifier = Modifier.fillMaxWidth(), shape = RoundedCornerShape(8.dp)) {
                        Icon(Icons.Default.Add, contentDescription = null)
                        Spacer(modifier = Modifier.width(6.dp))
                        Text("Nuovo Compito in Agenda")
                    }
                }
            }
        }
    }
}

@Composable
fun TeacherScrutinyTab() {
    LazyColumn(verticalArrangement = Arrangement.spacedBy(12.dp)) {
        item {
            Card(
                modifier = Modifier.fillMaxWidth(),
                shape = RoundedCornerShape(14.dp),
                colors = CardDefaults.cardColors(containerColor = Color(0xFF1E293B))
            ) {
                Column(modifier = Modifier.padding(16.dp)) {
                    Text("Gestione Scrutini & Scrutini Differiti", color = Color.White, fontWeight = FontWeight.Bold, fontSize = 16.sp)
                    Spacer(modifier = Modifier.height(4.dp))
                    Text("Delibera voti finali, condotta e saldo debiti formativi", color = Color.White.copy(alpha = 0.8f), fontSize = 13.sp)
                    Spacer(modifier = Modifier.height(12.dp))
                    Button(
                        onClick = { },
                        colors = ButtonDefaults.buttonColors(containerColor = Color(0xFFEA580C)),
                        shape = RoundedCornerShape(8.dp)
                    ) {
                        Icon(Icons.Default.EventRepeat, contentDescription = null)
                        Spacer(modifier = Modifier.width(6.dp))
                        Text("Apri Sessione Scrutinio Differito")
                    }
                }
            }
        }
    }
}

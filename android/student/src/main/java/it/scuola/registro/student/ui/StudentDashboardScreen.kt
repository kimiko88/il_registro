package it.scuola.registro.student.ui

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
import androidx.compose.ui.graphics.Brush
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.res.stringResource
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import it.scuola.registro.student.R

data class GradeItem(val subject: String, val grade: Double, val date: String, val type: String)
data class LessonItem(val hour: String, val subject: String, val room: String)
data class HomeworkItem(val subject: String, val description: String, val dueDate: String, val isCompleted: Boolean)

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun StudentDashboardScreen(
    studentName: String = "Mario Rossi",
    gpaAverage: Double = 7.8,
    onLogout: () -> Unit = {}
) {
    var selectedTab by remember { mutableStateOf(0) }

    Scaffold(
        topBar = {
            TopAppBar(
                title = {
                    Text(
                        when (selectedTab) {
                            0 -> stringResource(R.string.dashboard_title)
                            1 -> stringResource(R.string.grades_title)
                            2 -> stringResource(R.string.agenda_title)
                            3 -> stringResource(R.string.attendance_title)
                            else -> stringResource(R.string.report_card_title)
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
                    containerColor = MaterialTheme.colorScheme.primary,
                    titleContentColor = Color.White
                )
            )
        },
        bottomBar = {
            NavigationBar(
                containerColor = MaterialTheme.colorScheme.surface,
                tonalElevation = 8.dp
            ) {
                NavigationBarItem(
                    selected = selectedTab == 0,
                    onClick = { selectedTab = 0 },
                    icon = { Icon(Icons.Default.Dashboard, contentDescription = "Dashboard") },
                    label = { Text("Home", fontSize = 11.sp) }
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
                    icon = { Icon(Icons.Default.EventNote, contentDescription = "Agenda") },
                    label = { Text("Agenda", fontSize = 11.sp) }
                )
                NavigationBarItem(
                    selected = selectedTab == 3,
                    onClick = { selectedTab = 3 },
                    icon = { Icon(Icons.Default.CheckCircleOutline, contentDescription = "Presenze") },
                    label = { Text("Presenze", fontSize = 11.sp) }
                )
                NavigationBarItem(
                    selected = selectedTab == 4,
                    onClick = { selectedTab = 4 },
                    icon = { Icon(Icons.Default.Description, contentDescription = "Pagella") },
                    label = { Text("Pagella", fontSize = 11.sp) }
                )
            }
        }
    ) { paddingValues ->
        Box(
            modifier = Modifier
                .fillMaxSize()
                .padding(paddingValues)
        ) {
            when (selectedTab) {
                0 -> StudentHomeTab(studentName, gpaAverage)
                1 -> StudentGradesTab()
                2 -> StudentAgendaTab()
                3 -> StudentAttendanceTab()
                4 -> StudentReportCardTab()
            }
        }
    }
}

@Composable
fun StudentHomeTab(studentName: String, gpaAverage: Double) {
    val sampleGrades = remember {
        listOf(
            GradeItem("Matematica", 8.5, "30 Ago", "Scritto"),
            GradeItem("Italiano", 7.5, "28 Ago", "Orale"),
            GradeItem("Inglese", 8.0, "25 Ago", "Pratico"),
            GradeItem("Fisica", 7.0, "22 Ago", "Scritto")
        )
    }
    val sampleLessons = remember {
        listOf(
            LessonItem("08:00 - 09:00", "Matematica", "Aula 2B"),
            LessonItem("09:00 - 10:00", "Fisica", "Laboratorio"),
            LessonItem("10:00 - 11:00", "Italiano", "Aula 2B"),
            LessonItem("11:15 - 12:15", "Inglese", "Aula Lingue")
        )
    }

    LazyColumn(
        modifier = Modifier
            .fillMaxSize()
            .padding(16.dp),
        verticalArrangement = Arrangement.spacedBy(16.dp)
    ) {
        item {
            Card(
                modifier = Modifier.fillMaxWidth(),
                shape = RoundedCornerShape(20.dp),
                elevation = CardDefaults.cardElevation(defaultElevation = 6.dp)
            ) {
                Box(
                    modifier = Modifier
                        .fillMaxWidth()
                        .background(
                            Brush.horizontalGradient(
                                colors = listOf(
                                    MaterialTheme.colorScheme.primary,
                                    MaterialTheme.colorScheme.secondary
                                )
                            )
                        )
                        .padding(20.dp)
                ) {
                    Column {
                        Text(
                            text = stringResource(R.string.welcome_student, studentName),
                            color = Color.White,
                            fontSize = 22.sp,
                            fontWeight = FontWeight.Bold
                        )
                        Spacer(modifier = Modifier.height(8.dp))
                        Row(verticalAlignment = Alignment.CenterVertically) {
                            Surface(
                                color = Color.White.copy(alpha = 0.2f),
                                shape = RoundedCornerShape(12.dp)
                            ) {
                                Text(
                                    text = "${stringResource(R.string.gpa_average)}: ${"%.1f".format(gpaAverage)}",
                                    color = Color.White,
                                    modifier = Modifier.padding(horizontal = 12.dp, vertical = 6.dp),
                                    fontWeight = FontWeight.SemiBold
                                )
                            }
                        }
                    }
                }
            }
        }

        item {
            Row(
                modifier = Modifier.fillMaxWidth(),
                horizontalArrangement = Arrangement.spacedBy(12.dp)
            ) {
                StatCard(
                    modifier = Modifier.weight(1f),
                    title = stringResource(R.string.absences_count),
                    value = "3",
                    icon = Icons.Default.Warning,
                    color = MaterialTheme.colorScheme.tertiary
                )
                StatCard(
                    modifier = Modifier.weight(1f),
                    title = stringResource(R.string.lates_count),
                    value = "1",
                    icon = Icons.Default.Schedule,
                    color = MaterialTheme.colorScheme.secondary
                )
            }
        }

        item {
            Text(
                text = stringResource(R.string.recent_grades),
                fontSize = 18.sp,
                fontWeight = FontWeight.Bold
            )
        }

        items(sampleGrades) { grade ->
            GradeCard(grade = grade)
        }

        item {
            Spacer(modifier = Modifier.height(8.dp))
            Text(
                text = stringResource(R.string.today_lessons),
                fontSize = 18.sp,
                fontWeight = FontWeight.Bold
            )
        }

        items(sampleLessons) { lesson ->
            LessonCard(lesson = lesson)
        }
    }
}

@Composable
fun StudentGradesTab() {
    val allGrades = remember {
        listOf(
            GradeItem("Matematica", 8.5, "30 Ago", "Scritto"),
            GradeItem("Matematica", 7.0, "15 Ago", "Orale"),
            GradeItem("Italiano", 8.0, "28 Ago", "Tema"),
            GradeItem("Italiano", 7.5, "10 Ago", "Orale"),
            GradeItem("Inglese", 9.0, "25 Ago", "Pratico"),
            GradeItem("Fisica", 7.0, "22 Ago", "Scritto"),
            GradeItem("Storia", 8.0, "18 Ago", "Orale"),
            GradeItem("Filosofia", 8.5, "12 Ago", "Orale")
        )
    }

    LazyColumn(
        modifier = Modifier
            .fillMaxSize()
            .padding(16.dp),
        verticalArrangement = Arrangement.spacedBy(12.dp)
    ) {
        item {
            Card(
                modifier = Modifier.fillMaxWidth(),
                colors = CardDefaults.cardColors(containerColor = MaterialTheme.colorScheme.surfaceVariant),
                shape = RoundedCornerShape(16.dp)
            ) {
                Row(
                    modifier = Modifier
                        .fillMaxWidth()
                        .padding(16.dp),
                    horizontalArrangement = Arrangement.SpaceBetween,
                    verticalAlignment = Alignment.CenterVertically
                ) {
                    Column {
                        Text("Riepilogo Voti e Medie", fontWeight = FontWeight.Bold, fontSize = 18.sp)
                        Text("1° Quadrimestre • 8 Valutazioni", fontSize = 13.sp, color = Color.Gray)
                    }
                    Surface(
                        color = MaterialTheme.colorScheme.primary,
                        shape = RoundedCornerShape(12.dp)
                    ) {
                        Text("Media 7.9", color = Color.White, fontWeight = FontWeight.Bold, modifier = Modifier.padding(8.dp))
                    }
                }
            }
        }

        items(allGrades) { grade ->
            GradeCard(grade = grade)
        }
    }
}

@Composable
fun StudentAgendaTab() {
    val homeworkList = remember {
        listOf(
            HomeworkItem("Matematica", "Esercizi pag. 142 n. 12, 15, 18 sulle disequazioni esponenziali", "Domani", false),
            HomeworkItem("Fisica", "Relazione di laboratorio sul moto rettilineo uniforme", "Tra 2 giorni", false),
            HomeworkItem("Italiano", "Lettura Capitolo 8 I Promessi Sposi con riassunto", "Tra 3 giorni", true),
            HomeworkItem("Inglese", "Unit 4 vocabulary exercises & Grammar workbook", "Settimana prossima", true)
        )
    }

    LazyColumn(
        modifier = Modifier
            .fillMaxSize()
            .padding(16.dp),
        verticalArrangement = Arrangement.spacedBy(12.dp)
    ) {
        item {
            Text("Compiti e Scadenze Imminenti", fontSize = 18.sp, fontWeight = FontWeight.Bold)
        }

        items(homeworkList) { item ->
            Card(
                modifier = Modifier.fillMaxWidth(),
                shape = RoundedCornerShape(14.dp),
                colors = CardDefaults.cardColors(
                    containerColor = if (item.isCompleted) Color(0xFFF1F5F9) else MaterialTheme.colorScheme.surface
                )
            ) {
                Row(
                    modifier = Modifier
                        .fillMaxWidth()
                        .padding(16.dp),
                    verticalAlignment = Alignment.CenterVertically
                ) {
                    Icon(
                        imageVector = if (item.isCompleted) Icons.Default.CheckCircle else Icons.Default.RadioButtonUnchecked,
                        contentDescription = null,
                        tint = if (item.isCompleted) Color(0xFF10B981) else Color.Gray,
                        modifier = Modifier.size(24.dp)
                    )
                    Spacer(modifier = Modifier.width(12.dp))
                    Column(modifier = Modifier.weight(1f)) {
                        Text(text = item.subject, fontWeight = FontWeight.Bold, fontSize = 16.sp)
                        Text(text = item.description, fontSize = 13.sp, color = Color.DarkGray)
                        Spacer(modifier = Modifier.height(4.dp))
                        Text(text = "Scadenza: ${item.dueDate}", fontSize = 12.sp, color = MaterialTheme.colorScheme.primary, fontWeight = FontWeight.Medium)
                    }
                }
            }
        }
    }
}

@Composable
fun StudentAttendanceTab() {
    LazyColumn(
        modifier = Modifier
            .fillMaxSize()
            .padding(16.dp),
        verticalArrangement = Arrangement.spacedBy(14.dp)
    ) {
        item {
            Card(
                modifier = Modifier.fillMaxWidth(),
                shape = RoundedCornerShape(16.dp),
                colors = CardDefaults.cardColors(containerColor = MaterialTheme.colorScheme.surfaceVariant)
            ) {
                Column(modifier = Modifier.padding(16.dp)) {
                    Text("Statistiche Annuali Presenze", fontWeight = FontWeight.Bold, fontSize = 18.sp)
                    Spacer(modifier = Modifier.height(12.dp))
                    Row(
                        modifier = Modifier.fillMaxWidth(),
                        horizontalArrangement = Arrangement.SpaceAround
                    ) {
                        Column(horizontalAlignment = Alignment.CenterHorizontally) {
                            Text("184", fontSize = 22.sp, fontWeight = FontWeight.Bold, color = Color(0xFF10B981))
                            Text("Presenze", fontSize = 12.sp, color = Color.Gray)
                        }
                        Column(horizontalAlignment = Alignment.CenterHorizontally) {
                            Text("3", fontSize = 22.sp, fontWeight = FontWeight.Bold, color = Color(0xFFEF4444))
                            Text("Assenze", fontSize = 12.sp, color = Color.Gray)
                        }
                        Column(horizontalAlignment = Alignment.CenterHorizontally) {
                            Text("1", fontSize = 22.sp, fontWeight = FontWeight.Bold, color = Color(0xFFF59E0B))
                            Text("Ritardi", fontSize = 12.sp, color = Color.Gray)
                        }
                    }
                }
            }
        }

        item {
            Text("Registro Giustificazioni", fontWeight = FontWeight.Bold, fontSize = 18.sp)
        }

        items(listOf("26 Ago 2026 • Assenza (Giustificata dal genitore)", "18 Ago 2026 • Ritardo 10 min (Giustificato)", "04 Ago 2026 • Assenza (Giustificata)")) { log ->
            Card(
                modifier = Modifier.fillMaxWidth(),
                shape = RoundedCornerShape(12.dp)
            ) {
                Row(
                    modifier = Modifier
                        .fillMaxWidth()
                        .padding(14.dp),
                    horizontalArrangement = Arrangement.SpaceBetween,
                    verticalAlignment = Alignment.CenterVertically
                ) {
                    Text(log, fontSize = 13.sp, fontWeight = FontWeight.Medium)
                    Icon(Icons.Default.Verified, contentDescription = null, tint = Color(0xFF10B981), modifier = Modifier.size(20.dp))
                }
            }
        }
    }
}

@Composable
fun StudentReportCardTab() {
    val reportCardGrades = remember {
        listOf(
            "Italiano" to 8,
            "Latino" to 7,
            "Lingua Inglese" to 8,
            "Storia e Geografia" to 8,
            "Matematica" to 8,
            "Fisica" to 7,
            "Scienze Naturali" to 8,
            "Disegno e Storia dell'Arte" to 9,
            "Scienze Motorie" to 9,
            "Religione Cattolica" to 9
        )
    }

    LazyColumn(
        modifier = Modifier
            .fillMaxSize()
            .padding(16.dp),
        verticalArrangement = Arrangement.spacedBy(14.dp)
    ) {
        item {
            Card(
                modifier = Modifier.fillMaxWidth(),
                shape = RoundedCornerShape(16.dp),
                colors = CardDefaults.cardColors(containerColor = MaterialTheme.colorScheme.primary)
            ) {
                Column(modifier = Modifier.padding(16.dp)) {
                    Text("Documento di Valutazione Ufficiale", color = Color.White, fontWeight = FontWeight.Bold, fontSize = 18.sp)
                    Text("Esito Finale Scrutinio: AMMESSO / PROMOSSO", color = Color.White.copy(alpha = 0.9f), fontSize = 14.sp)
                    Spacer(modifier = Modifier.height(8.dp))
                    Text("Voto di Condotta: 9 • Crediti: 8", color = Color.White.copy(alpha = 0.9f), fontSize = 13.sp)
                }
            }
        }

        items(reportCardGrades) { (subject, grade) ->
            Card(
                modifier = Modifier.fillMaxWidth(),
                shape = RoundedCornerShape(12.dp)
            ) {
                Row(
                    modifier = Modifier
                        .fillMaxWidth()
                        .padding(14.dp),
                    horizontalArrangement = Arrangement.SpaceBetween,
                    verticalAlignment = Alignment.CenterVertically
                ) {
                    Text(subject, fontWeight = FontWeight.SemiBold, fontSize = 15.sp)
                    Surface(
                        color = if (grade >= 6) Color(0xFF10B981) else Color(0xFFEF4444),
                        shape = RoundedCornerShape(8.dp)
                    ) {
                        Text(
                            text = grade.toString(),
                            color = Color.White,
                            fontWeight = FontWeight.Bold,
                            modifier = Modifier.padding(horizontal = 12.dp, vertical = 4.dp)
                        )
                    }
                }
            }
        }

        item {
            Button(
                onClick = { },
                modifier = Modifier
                    .fillMaxWidth()
                    .height(48.dp),
                shape = RoundedCornerShape(12.dp)
            ) {
                Icon(Icons.Default.Download, contentDescription = null)
                Spacer(modifier = Modifier.width(8.dp))
                Text("Scarica Pagella Ufficiale (PDF)")
            }
        }
    }
}

@Composable
fun StatCard(
    modifier: Modifier = Modifier,
    title: String,
    value: String,
    icon: androidx.compose.ui.graphics.vector.ImageVector,
    color: Color
) {
    Card(
        modifier = modifier,
        shape = RoundedCornerShape(16.dp),
        colors = CardDefaults.cardColors(containerColor = MaterialTheme.colorScheme.surfaceVariant)
    ) {
        Column(modifier = Modifier.padding(16.dp)) {
            Icon(imageVector = icon, contentDescription = null, tint = color)
            Spacer(modifier = Modifier.height(8.dp))
            Text(text = value, fontSize = 24.sp, fontWeight = FontWeight.Bold)
            Text(text = title, fontSize = 12.sp, color = MaterialTheme.colorScheme.onSurfaceVariant)
        }
    }
}

@Composable
fun GradeCard(grade: GradeItem) {
    Card(
        modifier = Modifier.fillMaxWidth(),
        shape = RoundedCornerShape(14.dp),
        colors = CardDefaults.cardColors(containerColor = MaterialTheme.colorScheme.surface)
    ) {
        Row(
            modifier = Modifier
                .fillMaxWidth()
                .padding(16.dp),
            horizontalArrangement = Arrangement.SpaceBetween,
            verticalAlignment = Alignment.CenterVertically
        ) {
            Column {
                Text(text = grade.subject, fontWeight = FontWeight.Bold, fontSize = 16.sp)
                Text(text = "${grade.type} • ${grade.date}", fontSize = 12.sp, color = Color.Gray)
            }
            Surface(
                color = if (grade.grade >= 6.0) Color(0xFF10B981) else Color(0xFFEF4444),
                shape = RoundedCornerShape(10.dp)
            ) {
                Text(
                    text = grade.grade.toString(),
                    color = Color.White,
                    fontWeight = FontWeight.Bold,
                    fontSize = 16.sp,
                    modifier = Modifier.padding(horizontal = 14.dp, vertical = 6.dp)
                )
            }
        }
    }
}

@Composable
fun LessonCard(lesson: LessonItem) {
    Card(
        modifier = Modifier.fillMaxWidth(),
        shape = RoundedCornerShape(12.dp),
        colors = CardDefaults.cardColors(containerColor = MaterialTheme.colorScheme.surfaceVariant)
    ) {
        Row(
            modifier = Modifier
                .fillMaxWidth()
                .padding(14.dp),
            verticalAlignment = Alignment.CenterVertically,
            horizontalArrangement = Arrangement.SpaceBetween
        ) {
            Row(verticalAlignment = Alignment.CenterVertically) {
                Icon(imageVector = Icons.Default.Book, contentDescription = null, tint = MaterialTheme.colorScheme.primary)
                Spacer(modifier = Modifier.width(12.dp))
                Column {
                    Text(text = lesson.subject, fontWeight = FontWeight.SemiBold)
                    Text(text = lesson.room, fontSize = 12.sp, color = Color.Gray)
                }
            }
            Text(text = lesson.hour, fontSize = 13.sp, fontWeight = FontWeight.Medium)
        }
    }
}

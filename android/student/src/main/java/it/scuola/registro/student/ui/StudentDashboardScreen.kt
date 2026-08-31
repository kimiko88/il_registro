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
import it.scuola.registro.student.data.AttendanceRecord
import it.scuola.registro.student.data.GradeEntry
import it.scuola.registro.student.data.HomeworkAssignment
import it.scuola.registro.student.viewmodel.StudentViewModel

data class LessonItem(val hour: String, val subject: String, val room: String)

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun StudentDashboardScreen(
    studentName: String = "Mario Rossi",
    viewModel: StudentViewModel = remember { StudentViewModel() },
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
                        Icon(Icons.Default.Logout, contentDescription = "Logout", tint = Color.White)
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
                    icon = { Icon(Icons.Default.Dashboard, contentDescription = null) },
                    label = { Text(stringResource(R.string.dashboard_title), fontSize = 10.sp) }
                )
                NavigationBarItem(
                    selected = selectedTab == 1,
                    onClick = { selectedTab = 1 },
                    icon = { Icon(Icons.Default.Grade, contentDescription = null) },
                    label = { Text(stringResource(R.string.grades_title), fontSize = 10.sp) }
                )
                NavigationBarItem(
                    selected = selectedTab == 2,
                    onClick = { selectedTab = 2 },
                    icon = { Icon(Icons.Default.EventNote, contentDescription = null) },
                    label = { Text(stringResource(R.string.agenda_title), fontSize = 10.sp) }
                )
                NavigationBarItem(
                    selected = selectedTab == 3,
                    onClick = { selectedTab = 3 },
                    icon = { Icon(Icons.Default.CheckCircleOutline, contentDescription = null) },
                    label = { Text(stringResource(R.string.attendance_title), fontSize = 10.sp) }
                )
                NavigationBarItem(
                    selected = selectedTab == 4,
                    onClick = { selectedTab = 4 },
                    icon = { Icon(Icons.Default.Description, contentDescription = null) },
                    label = { Text(stringResource(R.string.report_card_title), fontSize = 10.sp) }
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
                0 -> StudentHomeTab(studentName, viewModel)
                1 -> StudentGradesTab(viewModel)
                2 -> StudentAgendaTab(viewModel)
                3 -> StudentAttendanceTab(viewModel)
                4 -> StudentReportCardTab(viewModel)
            }
        }
    }
}

@Composable
fun StudentHomeTab(studentName: String, viewModel: StudentViewModel) {
    val grades = viewModel.grades
    val gpa = viewModel.calculateGPA()

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
                                    text = "${stringResource(R.string.gpa_average)}: ${"%.1f".format(gpa)}",
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
                    value = viewModel.getTotalAbsences().toString(),
                    icon = Icons.Default.Warning,
                    color = MaterialTheme.colorScheme.tertiary
                )
                StatCard(
                    modifier = Modifier.weight(1f),
                    title = stringResource(R.string.lates_count),
                    value = viewModel.getTotalLates().toString(),
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

        items(grades) { grade ->
            GradeCard(grade = grade)
        }
    }
}

@Composable
fun StudentGradesTab(viewModel: StudentViewModel) {
    val allGrades = viewModel.grades
    val gpa = viewModel.calculateGPA()

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
                        Text(stringResource(R.string.grades_title), fontWeight = FontWeight.Bold, fontSize = 18.sp)
                        Text("${allGrades.size} valutazioni caricate dal database", fontSize = 13.sp, color = Color.Gray)
                    }
                    Surface(
                        color = MaterialTheme.colorScheme.primary,
                        shape = RoundedCornerShape(12.dp)
                    ) {
                        Text("${stringResource(R.string.gpa_average)}: ${"%.1f".format(gpa)}", color = Color.White, fontWeight = FontWeight.Bold, modifier = Modifier.padding(8.dp))
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
fun StudentAgendaTab(viewModel: StudentViewModel) {
    val homeworkList = viewModel.homeworkList

    LazyColumn(
        modifier = Modifier
            .fillMaxSize()
            .padding(16.dp),
        verticalArrangement = Arrangement.spacedBy(12.dp)
    ) {
        item {
            Text(stringResource(R.string.agenda_title), fontSize = 18.sp, fontWeight = FontWeight.Bold)
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
                    IconButton(onClick = { viewModel.toggleHomeworkCompletion(item.id) }) {
                        Icon(
                            imageVector = if (item.isCompleted) Icons.Default.CheckCircle else Icons.Default.RadioButtonUnchecked,
                            contentDescription = null,
                            tint = if (item.isCompleted) Color(0xFF10B981) else Color.Gray,
                            modifier = Modifier.size(24.dp)
                        )
                    }
                    Spacer(modifier = Modifier.width(8.dp))
                    Column(modifier = Modifier.weight(1f)) {
                        Text(text = item.subject, fontWeight = FontWeight.Bold, fontSize = 16.sp)
                        Text(text = item.description.ifEmpty { item.title }, fontSize = 13.sp, color = Color.DarkGray)
                        Spacer(modifier = Modifier.height(4.dp))
                        Text(text = "Scadenza: ${item.dueDate}", fontSize = 12.sp, color = MaterialTheme.colorScheme.primary, fontWeight = FontWeight.Medium)
                    }
                }
            }
        }
    }
}

@Composable
fun StudentAttendanceTab(viewModel: StudentViewModel) {
    val records = viewModel.attendanceRecords

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
                    Text(stringResource(R.string.attendance_title), fontWeight = FontWeight.Bold, fontSize = 18.sp)
                    Spacer(modifier = Modifier.height(12.dp))
                    Row(
                        modifier = Modifier.fillMaxWidth(),
                        horizontalArrangement = Arrangement.SpaceAround
                    ) {
                        Column(horizontalAlignment = Alignment.CenterHorizontally) {
                            Text("${viewModel.getTotalAbsences()}", fontSize = 22.sp, fontWeight = FontWeight.Bold, color = Color(0xFFEF4444))
                            Text(stringResource(R.string.absences_count), fontSize = 12.sp, color = Color.Gray)
                        }
                        Column(horizontalAlignment = Alignment.CenterHorizontally) {
                            Text("${viewModel.getTotalLates()}", fontSize = 22.sp, fontWeight = FontWeight.Bold, color = Color(0xFFF59E0B))
                            Text(stringResource(R.string.lates_count), fontSize = 12.sp, color = Color.Gray)
                        }
                    }
                }
            }
        }

        items(records) { record ->
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
                    Column {
                        Text("${record.date} • ${record.type}", fontSize = 14.sp, fontWeight = FontWeight.SemiBold)
                        if (record.reason.isNotBlank()) {
                            Text(record.reason, fontSize = 12.sp, color = Color.Gray)
                        }
                    }
                    if (record.isJustified) {
                        Icon(Icons.Default.Verified, contentDescription = null, tint = Color(0xFF10B981), modifier = Modifier.size(20.dp))
                    }
                }
            }
        }
    }
}

@Composable
fun StudentReportCardTab(viewModel: StudentViewModel) {
    val reportCard = viewModel.reportCard

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
                    Text(stringResource(R.string.report_card_title), color = Color.White, fontWeight = FontWeight.Bold, fontSize = 18.sp)
                    Text("Esito Scrutinio: ${reportCard?.finalDecision ?: "In Elaborazione"}", color = Color.White.copy(alpha = 0.9f), fontSize = 14.sp)
                    Spacer(modifier = Modifier.height(8.dp))
                    Text("Condotta: ${reportCard?.conductGrade ?: 8} • Crediti: ${reportCard?.credits ?: 0}", color = Color.White.copy(alpha = 0.9f), fontSize = 13.sp)
                }
            }
        }

        reportCard?.grades?.let { gradesMap ->
            items(gradesMap.toList()) { (subject, grade) ->
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
                            color = if (grade >= 6.0) Color(0xFF10B981) else Color(0xFFEF4444),
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
fun GradeCard(grade: GradeEntry) {
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

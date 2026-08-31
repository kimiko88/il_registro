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
import it.scuola.registro.teacher.data.StudentRollCall
import it.scuola.registro.teacher.viewmodel.TeacherViewModel

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun TeacherRegisterScreen(
    viewModel: TeacherViewModel = remember { TeacherViewModel() },
    onLogout: () -> Unit = {}
) {
    var selectedTab by remember { mutableStateOf(0) }

    Scaffold(
        topBar = {
            TopAppBar(
                title = {
                    Text(
                        when (selectedTab) {
                            0 -> stringResource(R.string.teacher_dashboard_title)
                            1 -> stringResource(R.string.take_attendance)
                            2 -> stringResource(R.string.add_grade)
                            3 -> stringResource(R.string.teacher_dashboard_title)
                            else -> stringResource(R.string.scrutiny_board)
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
                    icon = { Icon(Icons.Default.MenuBook, contentDescription = null) },
                    label = { Text(stringResource(R.string.teacher_dashboard_title), fontSize = 10.sp) }
                )
                NavigationBarItem(
                    selected = selectedTab == 1,
                    onClick = { selectedTab = 1 },
                    icon = { Icon(Icons.Default.Checklist, contentDescription = null) },
                    label = { Text(stringResource(R.string.take_attendance), fontSize = 10.sp) }
                )
                NavigationBarItem(
                    selected = selectedTab == 2,
                    onClick = { selectedTab = 2 },
                    icon = { Icon(Icons.Default.Grade, contentDescription = null) },
                    label = { Text(stringResource(R.string.add_grade), fontSize = 10.sp) }
                )
                NavigationBarItem(
                    selected = selectedTab == 3,
                    onClick = { selectedTab = 3 },
                    icon = { Icon(Icons.Default.EditCalendar, contentDescription = null) },
                    label = { Text(stringResource(R.string.teacher_dashboard_title), fontSize = 10.sp) }
                )
                NavigationBarItem(
                    selected = selectedTab == 4,
                    onClick = { selectedTab = 4 },
                    icon = { Icon(Icons.Default.Gavel, contentDescription = null) },
                    label = { Text(stringResource(R.string.scrutiny_board), fontSize = 10.sp) }
                )
            }
        }
    ) { padding ->
        Box(modifier = Modifier.fillMaxSize().padding(padding).padding(16.dp)) {
            when (selectedTab) {
                0 -> TeacherFirmaTab(viewModel)
                1 -> TeacherAppelloTab(viewModel)
                2 -> TeacherGradesEntryTab(viewModel)
                3 -> TeacherAgendaTab(viewModel)
                4 -> TeacherScrutinyTab(viewModel)
            }
        }
    }
}

@Composable
fun TeacherFirmaTab(viewModel: TeacherViewModel) {
    val session = viewModel.activeSession
    var topic by remember(session) { mutableStateOf(session?.lessonTopic ?: "") }
    var isSigned by remember(session) { mutableStateOf(session?.isSigned ?: false) }

    LazyColumn(verticalArrangement = Arrangement.spacedBy(14.dp)) {
        item {
            Card(
                modifier = Modifier.fillMaxWidth(),
                shape = RoundedCornerShape(16.dp),
                colors = CardDefaults.cardColors(containerColor = MaterialTheme.colorScheme.surfaceVariant)
            ) {
                Column(modifier = Modifier.padding(16.dp)) {
                    Text(stringResource(R.string.sign_hour), fontWeight = FontWeight.Bold, fontSize = 18.sp)
                    Spacer(modifier = Modifier.height(4.dp))
                    Text("${stringResource(R.string.select_class)} ${session?.className ?: ""} • ${session?.subject ?: ""}", fontSize = 13.sp, color = Color.Gray)
                    Spacer(modifier = Modifier.height(12.dp))
                    OutlinedTextField(
                        value = topic,
                        onValueChange = { topic = it },
                        label = { Text(stringResource(R.string.sign_hour)) },
                        modifier = Modifier.fillMaxWidth()
                    )
                    Spacer(modifier = Modifier.height(12.dp))
                    Button(
                        onClick = {
                            viewModel.signLessonHour(topic)
                            isSigned = true
                        },
                        shape = RoundedCornerShape(10.dp),
                        colors = ButtonDefaults.buttonColors(containerColor = if (isSigned) Color(0xFF10B981) else Color(0xFF2563EB)),
                        modifier = Modifier.fillMaxWidth()
                    ) {
                        Icon(Icons.Default.Edit, contentDescription = null, modifier = Modifier.size(18.dp))
                        Spacer(modifier = Modifier.width(8.dp))
                        Text(if (isSigned) stringResource(R.string.lesson_signed) else stringResource(R.string.sign_hour))
                    }
                }
            }
        }
    }
}

@Composable
fun TeacherAppelloTab(viewModel: TeacherViewModel) {
    val students = viewModel.rollCallList

    LazyColumn(verticalArrangement = Arrangement.spacedBy(12.dp)) {
        item {
            Text(stringResource(R.string.take_attendance), fontWeight = FontWeight.Bold, fontSize = 16.sp)
        }
        if (students.isEmpty()) {
            item {
                Text(stringResource(R.string.select_class), color = Color.Gray, fontSize = 13.sp)
            }
        } else {
            items(students) { student ->
                Card(modifier = Modifier.fillMaxWidth(), shape = RoundedCornerShape(12.dp)) {
                    Row(
                        modifier = Modifier.fillMaxWidth().padding(14.dp),
                        horizontalArrangement = Arrangement.SpaceBetween,
                        verticalAlignment = Alignment.CenterVertically
                    ) {
                        Text(text = student.studentName, fontWeight = FontWeight.SemiBold, fontSize = 15.sp)
                        Row(horizontalArrangement = Arrangement.spacedBy(6.dp)) {
                            FilterChip(
                                selected = student.status == "presente",
                                onClick = { viewModel.updateStudentStatus(student.studentId, "presente") },
                                label = { Text("P") }
                            )
                            FilterChip(
                                selected = student.status == "assente",
                                onClick = { viewModel.updateStudentStatus(student.studentId, "assente") },
                                label = { Text("A") }
                            )
                            FilterChip(
                                selected = student.status == "ritardo",
                                onClick = { viewModel.updateStudentStatus(student.studentId, "ritardo") },
                                label = { Text("R") }
                            )
                        }
                    }
                }
            }
        }
    }
}

@Composable
fun TeacherGradesEntryTab(viewModel: TeacherViewModel) {
    val students = viewModel.rollCallList

    LazyColumn(verticalArrangement = Arrangement.spacedBy(12.dp)) {
        item {
            Text(stringResource(R.string.add_grade), fontWeight = FontWeight.Bold, fontSize = 16.sp)
        }
        if (students.isEmpty()) {
            item {
                Text(stringResource(R.string.select_class), color = Color.Gray, fontSize = 13.sp)
            }
        } else {
            items(students) { student ->
                Card(modifier = Modifier.fillMaxWidth(), shape = RoundedCornerShape(12.dp)) {
                    Row(
                        modifier = Modifier.fillMaxWidth().padding(14.dp),
                        horizontalArrangement = Arrangement.SpaceBetween,
                        verticalAlignment = Alignment.CenterVertically
                    ) {
                        Text(student.studentName, fontWeight = FontWeight.SemiBold)
                        Button(onClick = { viewModel.submitGradeForStudent(student.studentId, 8.0, 1.0, "Scritto", "") }, shape = RoundedCornerShape(6.dp)) {
                            Text(stringResource(R.string.add_grade), fontSize = 12.sp)
                        }
                    }
                }
            }
        }
    }
}

@Composable
fun TeacherAgendaTab(viewModel: TeacherViewModel) {
    LazyColumn(verticalArrangement = Arrangement.spacedBy(12.dp)) {
        item {
            Card(modifier = Modifier.fillMaxWidth(), shape = RoundedCornerShape(14.dp)) {
                Column(modifier = Modifier.padding(16.dp)) {
                    Text(stringResource(R.string.teacher_dashboard_title), fontWeight = FontWeight.Bold, fontSize = 16.sp)
                    Spacer(modifier = Modifier.height(8.dp))
                    Button(onClick = { }, modifier = Modifier.fillMaxWidth(), shape = RoundedCornerShape(8.dp)) {
                        Icon(Icons.Default.Add, contentDescription = null)
                        Spacer(modifier = Modifier.width(6.dp))
                        Text(stringResource(R.string.teacher_dashboard_title))
                    }
                }
            }
        }
    }
}

@Composable
fun TeacherScrutinyTab(viewModel: TeacherViewModel) {
    LazyColumn(verticalArrangement = Arrangement.spacedBy(12.dp)) {
        item {
            Card(
                modifier = Modifier.fillMaxWidth(),
                shape = RoundedCornerShape(14.dp),
                colors = CardDefaults.cardColors(containerColor = Color(0xFF1E293B))
            ) {
                Column(modifier = Modifier.padding(16.dp)) {
                    Text(stringResource(R.string.scrutiny_board), color = Color.White, fontWeight = FontWeight.Bold, fontSize = 16.sp)
                    Spacer(modifier = Modifier.height(4.dp))
                    Text(stringResource(R.string.scrutiny_board), color = Color.White.copy(alpha = 0.8f), fontSize = 13.sp)
                    Spacer(modifier = Modifier.height(12.dp))
                    Button(
                        onClick = { },
                        colors = ButtonDefaults.buttonColors(containerColor = Color(0xFFEA580C)),
                        shape = RoundedCornerShape(8.dp)
                    ) {
                        Icon(Icons.Default.EventRepeat, contentDescription = null)
                        Spacer(modifier = Modifier.width(6.dp))
                        Text(stringResource(R.string.scrutiny_board))
                    }
                }
            }
        }
    }
}

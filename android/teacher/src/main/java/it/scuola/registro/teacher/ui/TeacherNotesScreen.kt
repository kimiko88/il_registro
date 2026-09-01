package it.scuola.registro.teacher.ui

import androidx.compose.foundation.layout.*
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.*
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.res.stringResource
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import it.scuola.registro.teacher.R

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun TeacherNotesScreen(
    onBack: () -> Unit = {},
    onSubmitNote: (studentName: String, noteType: String, noteText: String) -> Unit = { _, _, _ -> }
) {
    var studentName by remember { mutableStateOf("") }
    var noteText by remember { mutableStateOf("") }
    var noteType by remember { mutableStateOf("Disciplinare") }
    var isSubmitted by remember { mutableStateOf(false) }

    Scaffold(
        topBar = {
            TopAppBar(
                title = { Text(stringResource(R.string.teacher_dashboard_title), fontWeight = FontWeight.Bold) },
                navigationIcon = {
                    IconButton(onClick = onBack) {
                        Icon(Icons.Default.ArrowBack, contentDescription = "Back")
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
                        Text(stringResource(R.string.select_class), fontWeight = FontWeight.Bold, fontSize = 16.sp)
                        Text(stringResource(R.string.take_attendance), fontSize = 12.sp, color = Color.DarkGray)
                    }
                }
            }

            if (isSubmitted) {
                item {
                    Card(
                        modifier = Modifier.fillMaxWidth(),
                        shape = RoundedCornerShape(12.dp),
                        colors = CardDefaults.cardColors(containerColor = Color(0xFFD1FAE5))
                    ) {
                        Column(modifier = Modifier.padding(16.dp)) {
                            Text("Nota Inserita con Successo!", fontWeight = FontWeight.Bold, color = Color(0xFF065F46))
                        }
                    }
                }
            } else {
                item {
                    Card(modifier = Modifier.fillMaxWidth(), shape = RoundedCornerShape(12.dp)) {
                        Column(modifier = Modifier.padding(14.dp), verticalArrangement = Arrangement.spacedBy(10.dp)) {
                            OutlinedTextField(
                                value = studentName,
                                onValueChange = { studentName = it },
                                label = { Text("Studente") },
                                modifier = Modifier.fillMaxWidth()
                            )

                            OutlinedTextField(
                                value = noteType,
                                onValueChange = { noteType = it },
                                label = { Text("Tipologia") },
                                modifier = Modifier.fillMaxWidth()
                            )

                            OutlinedTextField(
                                value = noteText,
                                onValueChange = { noteText = it },
                                label = { Text("Descrizione") },
                                minLines = 3,
                                modifier = Modifier.fillMaxWidth()
                            )

                            Button(
                                onClick = {
                                    onSubmitNote(studentName, noteType, noteText)
                                    isSubmitted = true
                                },
                                shape = RoundedCornerShape(8.dp),
                                modifier = Modifier.fillMaxWidth().padding(top = 4.dp),
                                enabled = studentName.isNotBlank() && noteText.isNotBlank()
                            ) {
                                Text(stringResource(R.string.add_grade))
                            }
                        }
                    }
                }
            }
        }
    }
}

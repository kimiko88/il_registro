package it.scuola.registro.parent.ui

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

data class ParentDisciplinaryNoteItem(
    val id: String,
    val teacherAndDate: String,
    val description: String,
    var isAcknowledged: Boolean = false
)

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun ParentNotesScreen(
    notes: List<ParentDisciplinaryNoteItem> = emptyList(),
    onAcknowledge: (String) -> Unit = {},
    onBack: () -> Unit = {}
) {
    val acknowledgedSet = remember { mutableStateListOf<String>() }

    Scaffold(
        topBar = {
            TopAppBar(
                title = { Text(stringResource(R.string.parent_dashboard_title), fontWeight = FontWeight.Bold) },
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
                        Text(stringResource(R.string.child_grades), fontWeight = FontWeight.Bold, fontSize = 16.sp)
                        Text(stringResource(R.string.pending_justifications), fontSize = 12.sp, color = Color.DarkGray)
                    }
                }
            }

            if (notes.isNotEmpty()) {
                items(notes) { note ->
                    val isAck = note.isAcknowledged || acknowledgedSet.contains(note.id)
                    Card(modifier = Modifier.fillMaxWidth(), shape = RoundedCornerShape(12.dp)) {
                        Column(modifier = Modifier.padding(14.dp)) {
                            Row(
                                modifier = Modifier.fillMaxWidth(),
                                horizontalArrangement = Arrangement.SpaceBetween,
                                verticalAlignment = Alignment.CenterVertically
                            ) {
                                Text(stringResource(R.string.pending_justifications), fontWeight = FontWeight.Bold)
                                Surface(
                                    color = if (isAck) Color(0xFF10B981) else Color(0xFFEF4444),
                                    shape = RoundedCornerShape(6.dp)
                                ) {
                                    Text(
                                        if (isAck) "Firmata" else "In Attesa",
                                        color = Color.White,
                                        modifier = Modifier.padding(horizontal = 6.dp, vertical = 2.dp),
                                        fontSize = 10.sp
                                    )
                                }
                            }
                            Text(note.teacherAndDate, fontSize = 12.sp, color = MaterialTheme.colorScheme.primary, fontWeight = FontWeight.SemiBold, modifier = Modifier.padding(vertical = 4.dp))
                            Text(note.description, fontSize = 12.sp, color = Color.Gray)

                            if (!isAck) {
                                Button(
                                    onClick = {
                                        acknowledgedSet.add(note.id)
                                        onAcknowledge(note.id)
                                    },
                                    shape = RoundedCornerShape(8.dp),
                                    modifier = Modifier.fillMaxWidth().padding(top = 8.dp)
                                ) {
                                    Icon(Icons.Default.Check, contentDescription = null, modifier = Modifier.size(16.dp))
                                    Spacer(modifier = Modifier.width(6.dp))
                                    Text(stringResource(R.string.sign_hour))
                                }
                            }
                        }
                    }
                }
            }
        }
    }
}

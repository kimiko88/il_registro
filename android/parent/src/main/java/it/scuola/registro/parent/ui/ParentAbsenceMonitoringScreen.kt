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
import androidx.compose.ui.res.stringResource
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import it.scuola.registro.parent.R

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun ParentAbsenceMonitoringScreen(
    studentName: String = "",
    absentHours: Int = 0,
    maxAllowedAbsentHours: Int = 247,
    totalYearHours: Int = 990,
    onBack: () -> Unit = {}
) {
    val progress = if (maxAllowedAbsentHours > 0) (absentHours.toFloat() / maxAllowedAbsentHours.toFloat()).coerceIn(0f, 1f) else 0f
    val isCritical = progress > 0.8f

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
                        Text(stringResource(R.string.parent_dashboard_title), fontWeight = FontWeight.Bold, fontSize = 16.sp)
                        Text(stringResource(R.string.pending_justifications), fontSize = 12.sp, color = Color.DarkGray)
                    }
                }
            }

            if (studentName.isNotEmpty()) {
                item {
                    Text(studentName, fontWeight = FontWeight.Bold, fontSize = 16.sp)
                }
            }

            item {
                Card(modifier = Modifier.fillMaxWidth(), shape = RoundedCornerShape(12.dp)) {
                    Column(modifier = Modifier.padding(14.dp), verticalArrangement = Arrangement.spacedBy(8.dp)) {
                        Row(
                            modifier = Modifier.fillMaxWidth(),
                            horizontalArrangement = Arrangement.SpaceBetween,
                            verticalAlignment = Alignment.CenterVertically
                        ) {
                            Text(stringResource(R.string.pending_justifications), fontWeight = FontWeight.SemiBold)
                            Text(
                                "$absentHours / $maxAllowedAbsentHours h",
                                fontWeight = FontWeight.Bold,
                                color = if (isCritical) Color(0xFFEF4444) else Color(0xFF10B981)
                            )
                        }
                        LinearProgressIndicator(
                            progress = progress,
                            modifier = Modifier.fillMaxWidth().padding(vertical = 4.dp),
                            color = if (isCritical) Color(0xFFEF4444) else Color(0xFF10B981)
                        )
                        Divider()
                        Row(modifier = Modifier.fillMaxWidth(), horizontalArrangement = Arrangement.SpaceBetween) {
                            Text("Monte Ore Totale:", fontWeight = FontWeight.SemiBold)
                            Text("$totalYearHours h", fontWeight = FontWeight.Bold)
                        }
                    }
                }
            }
        }
    }
}

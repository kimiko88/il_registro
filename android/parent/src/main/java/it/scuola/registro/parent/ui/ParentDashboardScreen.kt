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
import it.scuola.registro.parent.data.ColloquioBooking
import it.scuola.registro.parent.data.ParentChild
import it.scuola.registro.parent.data.PendingAbsence
import it.scuola.registro.parent.viewmodel.ParentViewModel

data class ParentGradeDisplayItem(
    val subject: String,
    val gradeWithDetails: String
)

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun ParentDashboardScreen(
    token: String? = null,
    viewModel: ParentViewModel = remember { ParentViewModel() },
    gradesList: List<ParentGradeDisplayItem> = emptyList(),
    circularsList: List<String> = emptyList(),
    onLogout: () -> Unit = {}
) {
    LaunchedEffect(token) {
        if (!token.isNullOrBlank()) {
            viewModel.loadFromDatabase(token)
        }
    }

    val children = viewModel.children
    val selectedChild = children.find { it.id == viewModel.selectedChildId } ?: children.firstOrNull()
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
                            else -> stringResource(R.string.parent_dashboard_title)
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
                    icon = { Icon(Icons.Default.FamilyRestroom, contentDescription = null) },
                    label = { Text(stringResource(R.string.parent_dashboard_title), fontSize = 10.sp) }
                )
                NavigationBarItem(
                    selected = selectedTab == 1,
                    onClick = { selectedTab = 1 },
                    icon = { Icon(Icons.Default.Grade, contentDescription = null) },
                    label = { Text(stringResource(R.string.child_grades), fontSize = 10.sp) }
                )
                NavigationBarItem(
                    selected = selectedTab == 2,
                    onClick = { selectedTab = 2 },
                    icon = { Icon(Icons.Default.CheckCircle, contentDescription = null) },
                    label = { Text(stringResource(R.string.pending_justifications), fontSize = 10.sp) }
                )
                NavigationBarItem(
                    selected = selectedTab == 3,
                    onClick = { selectedTab = 3 },
                    icon = { Icon(Icons.Default.CalendarMonth, contentDescription = null) },
                    label = { Text(stringResource(R.string.upcoming_colloqui), fontSize = 10.sp) }
                )
                NavigationBarItem(
                    selected = selectedTab == 4,
                    onClick = { selectedTab = 4 },
                    icon = { Icon(Icons.Default.Campaign, contentDescription = null) },
                    label = { Text(stringResource(R.string.parent_dashboard_title), fontSize = 10.sp) }
                )
            }
        }
    ) { padding ->
        Column(
            modifier = Modifier
                .fillMaxSize()
                .padding(padding)
        ) {
            if (children.isNotEmpty()) {
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
                        Text(stringResource(R.string.switch_child), fontSize = 13.sp, color = Color.Gray)
                        Row(horizontalArrangement = Arrangement.spacedBy(8.dp)) {
                            children.forEach { child ->
                                FilterChip(
                                    selected = viewModel.selectedChildId == child.id,
                                    onClick = { viewModel.selectChild(child.id) },
                                    label = { Text("${child.firstName} ${child.lastName}", fontSize = 12.sp) }
                                )
                            }
                        }
                    }
                }
            }

            Box(modifier = Modifier.fillMaxSize().padding(horizontal = 16.dp)) {
                selectedChild?.let { child ->
                    when (selectedTab) {
                        0 -> ParentOverviewTab(child)
                        1 -> ParentGradesTab(gradesList)
                        2 -> ParentJustificationsTab(viewModel)
                        3 -> ParentColloquiTab(viewModel)
                        4 -> ParentCommunicationsTab(circularsList)
                    }
                } ?: run {
                    Box(modifier = Modifier.fillMaxSize(), contentAlignment = Alignment.Center) {
                        Text(stringResource(R.string.no_pending_justifications))
                    }
                }
            }
        }
    }
}

@Composable
fun ParentOverviewTab(child: ParentChild) {
    LazyColumn(verticalArrangement = Arrangement.spacedBy(14.dp)) {
        item {
            Card(
                modifier = Modifier.fillMaxWidth(),
                shape = RoundedCornerShape(16.dp),
                colors = CardDefaults.cardColors(containerColor = Color(0xFF0F172A))
            ) {
                Column(modifier = Modifier.padding(18.dp)) {
                    Text("${child.firstName} ${child.lastName}", color = Color.White, fontWeight = FontWeight.Bold, fontSize = 20.sp)
                    Text(child.className, color = Color.White.copy(alpha = 0.8f), fontSize = 14.sp)
                    Spacer(modifier = Modifier.height(12.dp))
                    Row(horizontalArrangement = Arrangement.spacedBy(16.dp)) {
                        Surface(color = Color.White.copy(alpha = 0.2f), shape = RoundedCornerShape(8.dp)) {
                            Text(stringResource(R.string.parent_dashboard_title), color = Color.White, modifier = Modifier.padding(6.dp), fontWeight = FontWeight.Bold)
                        }
                    }
                }
            }
        }
    }
}

@Composable
fun ParentGradesTab(grades: List<ParentGradeDisplayItem>) {
    LazyColumn(verticalArrangement = Arrangement.spacedBy(12.dp)) {
        if (grades.isEmpty()) {
            item {
                Text(stringResource(R.string.child_grades), fontSize = 13.sp, color = Color.Gray)
            }
        } else {
            items(grades) { item ->
                Card(modifier = Modifier.fillMaxWidth(), shape = RoundedCornerShape(12.dp)) {
                    Row(
                        modifier = Modifier.fillMaxWidth().padding(14.dp),
                        horizontalArrangement = Arrangement.SpaceBetween,
                        verticalAlignment = Alignment.CenterVertically
                    ) {
                        Text(item.subject, fontWeight = FontWeight.Bold)
                        Text(item.gradeWithDetails, color = Color(0xFF0D9488), fontWeight = FontWeight.SemiBold)
                    }
                }
            }
        }
    }
}

@Composable
fun ParentJustificationsTab(viewModel: ParentViewModel) {
    val absences = viewModel.getAbsencesForSelectedChild()

    LazyColumn(verticalArrangement = Arrangement.spacedBy(12.dp)) {
        item {
            Text(stringResource(R.string.pending_justifications), fontWeight = FontWeight.Bold, fontSize = 16.sp)
        }
        if (absences.isEmpty()) {
            item {
                Text(stringResource(R.string.no_pending_justifications), fontSize = 13.sp, color = Color.Gray)
            }
        } else {
            items(absences) { item ->
                Card(modifier = Modifier.fillMaxWidth(), shape = RoundedCornerShape(12.dp)) {
                    Row(
                        modifier = Modifier.fillMaxWidth().padding(14.dp),
                        horizontalArrangement = Arrangement.SpaceBetween,
                        verticalAlignment = Alignment.CenterVertically
                    ) {
                        Column {
                            Text(item.date, fontWeight = FontWeight.Bold)
                            Text(item.reason.ifEmpty { item.type }, fontSize = 12.sp, color = Color.Gray)
                        }
                        if (item.isJustified) {
                            Surface(color = Color(0xFF10B981), shape = RoundedCornerShape(6.dp)) {
                                Text(stringResource(R.string.pending_justifications), color = Color.White, modifier = Modifier.padding(horizontal = 8.dp, vertical = 4.dp), fontSize = 12.sp)
                            }
                        } else {
                            Button(onClick = { viewModel.justifyAbsence(item.id, "Motivata") }, shape = RoundedCornerShape(8.dp)) {
                                Text(stringResource(R.string.justify_absence), fontSize = 12.sp)
                            }
                        }
                    }
                }
            }
        }
    }
}

@Composable
fun ParentColloquiTab(viewModel: ParentViewModel) {
    val slots = viewModel.availableColloqui

    LazyColumn(verticalArrangement = Arrangement.spacedBy(12.dp)) {
        item {
            Text(stringResource(R.string.upcoming_colloqui), fontWeight = FontWeight.Bold, fontSize = 16.sp)
        }
        items(slots) { slot ->
            Card(modifier = Modifier.fillMaxWidth(), shape = RoundedCornerShape(12.dp)) {
                Row(
                    modifier = Modifier.fillMaxWidth().padding(14.dp),
                    horizontalArrangement = Arrangement.SpaceBetween,
                    verticalAlignment = Alignment.CenterVertically
                ) {
                    Column {
                        Text(slot.teacherName, fontWeight = FontWeight.Bold)
                        Text("${slot.subject} • ${slot.dateTimeSlot}", fontSize = 12.sp, color = Color.Gray)
                    }
                    Button(
                        onClick = { viewModel.bookColloquio(slot.id) },
                        colors = ButtonDefaults.buttonColors(containerColor = if (slot.isConfirmed) Color(0xFF10B981) else Color(0xFF0D9488)),
                        shape = RoundedCornerShape(8.dp)
                    ) {
                        Text(if (slot.isConfirmed) stringResource(R.string.upcoming_colloqui) else stringResource(R.string.book_colloquio), fontSize = 12.sp)
                    }
                }
            }
        }
    }
}

@Composable
fun ParentCommunicationsTab(circulars: List<String> = emptyList()) {
    LazyColumn(verticalArrangement = Arrangement.spacedBy(12.dp)) {
        if (circulars.isEmpty()) {
            item {
                Text(stringResource(R.string.parent_dashboard_title), fontSize = 13.sp, color = Color.Gray)
            }
        } else {
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
}

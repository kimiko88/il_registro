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
import it.scuola.registro.secretary.viewmodel.SecretaryViewModel

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun SecretaryDashboardScreen(
    viewModel: SecretaryViewModel = remember { SecretaryViewModel() },
    auditLogs: List<String> = emptyList(),
    onLogout: () -> Unit = {}
) {
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
                    icon = { Icon(Icons.Default.Dashboard, contentDescription = null) },
                    label = { Text(stringResource(R.string.secretary_dashboard_title), fontSize = 10.sp) }
                )
                NavigationBarItem(
                    selected = selectedTab == 1,
                    onClick = { selectedTab = 1 },
                    icon = { Icon(Icons.Default.People, contentDescription = null) },
                    label = { Text(stringResource(R.string.user_management), fontSize = 10.sp) }
                )
                NavigationBarItem(
                    selected = selectedTab == 2,
                    onClick = { selectedTab = 2 },
                    icon = { Icon(Icons.Default.Class, contentDescription = null) },
                    label = { Text(stringResource(R.string.scrutiny_supervision), fontSize = 10.sp) }
                )
                NavigationBarItem(
                    selected = selectedTab == 3,
                    onClick = { selectedTab = 3 },
                    icon = { Icon(Icons.Default.Description, contentDescription = null) },
                    label = { Text(stringResource(R.string.generate_certificates), fontSize = 10.sp) }
                )
                NavigationBarItem(
                    selected = selectedTab == 4,
                    onClick = { selectedTab = 4 },
                    icon = { Icon(Icons.Default.Security, contentDescription = null) },
                    label = { Text(stringResource(R.string.audit_logs), fontSize = 10.sp) }
                )
            }
        }
    ) { padding ->
        Box(modifier = Modifier.fillMaxSize().padding(padding).padding(16.dp)) {
            when (selectedTab) {
                0 -> SecretaryOverviewTab(viewModel)
                1 -> SecretaryUsersTab(viewModel)
                2 -> SecretaryScrutinyTab(viewModel)
                3 -> SecretaryCertificatesTab(viewModel)
                4 -> SecretaryAuditTab(auditLogs)
            }
        }
    }
}

@Composable
fun SecretaryOverviewTab(viewModel: SecretaryViewModel) {
    val studentCount = viewModel.getUsersByRole("student").size
    val teacherCount = viewModel.getUsersByRole("teacher").size

    LazyColumn(verticalArrangement = Arrangement.spacedBy(14.dp)) {
        item {
            Row(horizontalArrangement = Arrangement.spacedBy(12.dp)) {
                MetricCard(modifier = Modifier.weight(1f), title = stringResource(R.string.total_students), value = studentCount.toString(), icon = Icons.Default.People)
                MetricCard(modifier = Modifier.weight(1f), title = stringResource(R.string.total_teachers), value = teacherCount.toString(), icon = Icons.Default.School)
            }
        }
        item {
            Card(modifier = Modifier.fillMaxWidth(), shape = RoundedCornerShape(16.dp)) {
                Column(modifier = Modifier.padding(16.dp)) {
                    Text(stringResource(R.string.secretary_dashboard_title), fontWeight = FontWeight.Bold, fontSize = 16.sp)
                    Spacer(modifier = Modifier.height(8.dp))
                    Text(stringResource(R.string.system_logs), fontSize = 13.sp, color = Color.Gray)
                }
            }
        }
    }
}

@Composable
fun SecretaryUsersTab(viewModel: SecretaryViewModel) {
    val users = viewModel.usersList

    LazyColumn(verticalArrangement = Arrangement.spacedBy(12.dp)) {
        item {
            Row(modifier = Modifier.fillMaxWidth(), horizontalArrangement = Arrangement.SpaceBetween) {
                Text(stringResource(R.string.user_management), fontWeight = FontWeight.Bold, fontSize = 16.sp)
                Button(onClick = { }, shape = RoundedCornerShape(8.dp)) {
                    Icon(Icons.Default.Add, contentDescription = null, modifier = Modifier.size(16.dp))
                    Spacer(modifier = Modifier.width(4.dp))
                    Text(stringResource(R.string.add_user_button), fontSize = 12.sp)
                }
            }
        }
        if (users.isEmpty()) {
            item {
                Text(stringResource(R.string.user_management), color = Color.Gray, fontSize = 13.sp)
            }
        } else {
            items(users) { user ->
                Card(modifier = Modifier.fillMaxWidth(), shape = RoundedCornerShape(12.dp)) {
                    Row(
                        modifier = Modifier.fillMaxWidth().padding(14.dp),
                        horizontalArrangement = Arrangement.SpaceBetween,
                        verticalAlignment = Alignment.CenterVertically
                    ) {
                        Column {
                            Text("${user.firstName} ${user.lastName}", fontWeight = FontWeight.Bold)
                            Text("${user.role} • ${user.email}", fontSize = 12.sp, color = Color.Gray)
                        }
                        IconButton(onClick = { }) {
                            Icon(Icons.Default.Edit, contentDescription = null, tint = Color(0xFF581C87))
                        }
                    }
                }
            }
        }
    }
}

@Composable
fun SecretaryScrutinyTab(viewModel: SecretaryViewModel) {
    val classes = viewModel.scrutinyClasses

    LazyColumn(verticalArrangement = Arrangement.spacedBy(12.dp)) {
        item {
            Text(stringResource(R.string.scrutiny_supervision), fontWeight = FontWeight.Bold, fontSize = 16.sp)
        }
        if (classes.isEmpty()) {
            item {
                Text(stringResource(R.string.scrutiny_supervision), color = Color.Gray, fontSize = 13.sp)
            }
        } else {
            items(classes) { cls ->
                Card(modifier = Modifier.fillMaxWidth(), shape = RoundedCornerShape(12.dp)) {
                    Row(
                        modifier = Modifier.fillMaxWidth().padding(14.dp),
                        horizontalArrangement = Arrangement.SpaceBetween,
                        verticalAlignment = Alignment.CenterVertically
                    ) {
                        Text(cls.className, fontWeight = FontWeight.Bold)
                        Surface(
                            color = if (cls.isLocked) Color(0xFF10B981) else Color(0xFFEA580C),
                            shape = RoundedCornerShape(6.dp)
                        ) {
                            Text(if (cls.isLocked) "Bloccato/Chiuso" else "Aperto/In corso", color = Color.White, modifier = Modifier.padding(horizontal = 8.dp, vertical = 4.dp), fontSize = 12.sp)
                        }
                    }
                }
            }
        }
    }
}

@Composable
fun SecretaryCertificatesTab(viewModel: SecretaryViewModel) {
    val certs = viewModel.certificateRequests

    LazyColumn(verticalArrangement = Arrangement.spacedBy(12.dp)) {
        item {
            Text(stringResource(R.string.generate_certificates), fontWeight = FontWeight.Bold, fontSize = 16.sp)
        }
        if (certs.isEmpty()) {
            item {
                Text(stringResource(R.string.generate_certificates), color = Color.Gray, fontSize = 13.sp)
            }
        } else {
            items(certs) { cert ->
                Card(modifier = Modifier.fillMaxWidth(), shape = RoundedCornerShape(12.dp)) {
                    Row(
                        modifier = Modifier.fillMaxWidth().padding(14.dp),
                        horizontalArrangement = Arrangement.SpaceBetween,
                        verticalAlignment = Alignment.CenterVertically
                    ) {
                        Text("Certificato: ${cert.certificateType}", fontWeight = FontWeight.Medium, fontSize = 13.sp)
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
}

@Composable
fun SecretaryAuditTab(auditLogs: List<String> = emptyList()) {
    LazyColumn(verticalArrangement = Arrangement.spacedBy(10.dp)) {
        item {
            Text(stringResource(R.string.audit_logs), fontWeight = FontWeight.Bold, fontSize = 16.sp)
        }
        if (auditLogs.isEmpty()) {
            item {
                Text(stringResource(R.string.audit_logs), color = Color.Gray, fontSize = 13.sp)
            }
        } else {
            items(auditLogs) { log ->
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

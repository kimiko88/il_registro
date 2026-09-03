package it.scuola.registro.secretary.ui

import androidx.compose.foundation.background
import androidx.compose.foundation.clickable
import androidx.compose.foundation.layout.*
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.lazy.LazyRow
import androidx.compose.foundation.lazy.items
import androidx.compose.foundation.shape.CircleShape
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.*
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.graphics.vector.ImageVector
import androidx.compose.ui.res.stringResource
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.text.input.PasswordVisualTransformation
import androidx.compose.ui.text.input.VisualTransformation
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp
import it.scuola.registro.secretary.R
import it.scuola.registro.secretary.data.CreateUserPayload
import it.scuola.registro.secretary.data.ManagedUser
import it.scuola.registro.secretary.data.UpdateUserPayload
import it.scuola.registro.secretary.viewmodel.SecretaryViewModel
import kotlinx.coroutines.launch

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun SecretaryDashboardScreen(
    token: String? = null,
    viewModel: SecretaryViewModel = remember { SecretaryViewModel() },
    onLogout: () -> Unit = {}
) {
    var selectedTab by remember { mutableIntStateOf(0) }
    val coroutineScope = rememberCoroutineScope()
    val snackbarHostState = remember { SnackbarHostState() }

    // Dialog States
    var showAddUserDialog by remember { mutableStateOf(false) }
    var userToEdit by remember { mutableStateOf<ManagedUser?>(null) }
    var userToDelete by remember { mutableStateOf<ManagedUser?>(null) }
    var showGenerateCertDialog by remember { mutableStateOf(false) }
    var downloadedCertProtocol by remember { mutableStateOf<String?>(null) }

    // Initial load
    LaunchedEffect(token) {
        if (!token.isNullOrBlank()) {
            viewModel.loadFromDatabase(token)
        }
    }

    // Feedback Snackbar
    LaunchedEffect(viewModel.errorMessage, viewModel.successMessage) {
        viewModel.errorMessage?.let {
            snackbarHostState.showSnackbar(it, duration = SnackbarDuration.Short)
            viewModel.clearMessages()
        }
        viewModel.successMessage?.let {
            snackbarHostState.showSnackbar(it, duration = SnackbarDuration.Short)
            viewModel.clearMessages()
        }
    }

    Scaffold(
        snackbarHost = { SnackbarHost(snackbarHostState) },
        topBar = {
            TopAppBar(
                title = {
                    Text(
                        when (selectedTab) {
                            0 -> stringResource(R.string.secretary_dashboard_title)
                            1 -> stringResource(R.string.user_management)
                            2 -> stringResource(R.string.scrutiny_supervision)
                            3 -> stringResource(R.string.generate_certificates)
                            4 -> stringResource(R.string.audit_logs)
                            else -> stringResource(R.string.profile_title)
                        },
                        fontWeight = FontWeight.Bold
                    )
                },
                actions = {
                    IconButton(onClick = {
                        if (!token.isNullOrBlank()) {
                            coroutineScope.launch { viewModel.loadFromDatabase(token) }
                        }
                    }) {
                        Icon(Icons.Default.Refresh, contentDescription = stringResource(R.string.btn_refresh), tint = Color.White)
                    }
                    IconButton(onClick = { selectedTab = 5 }) {
                        Icon(Icons.Default.ManageAccounts, contentDescription = stringResource(R.string.profile_title), tint = Color.White)
                    }
                    IconButton(onClick = onLogout) {
                        Icon(Icons.Default.Logout, contentDescription = stringResource(R.string.logout), tint = Color.White)
                    }
                },
                colors = TopAppBarDefaults.topAppBarColors(
                    containerColor = Color(0xFF581C87),
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
                    label = { Text(stringResource(R.string.secretary_dashboard_title), fontSize = 10.sp, maxLines = 1) }
                )
                NavigationBarItem(
                    selected = selectedTab == 1,
                    onClick = { selectedTab = 1 },
                    icon = { Icon(Icons.Default.People, contentDescription = null) },
                    label = { Text(stringResource(R.string.user_management), fontSize = 10.sp, maxLines = 1) }
                )
                NavigationBarItem(
                    selected = selectedTab == 2,
                    onClick = { selectedTab = 2 },
                    icon = { Icon(Icons.Default.Class, contentDescription = null) },
                    label = { Text(stringResource(R.string.scrutiny_supervision), fontSize = 10.sp, maxLines = 1) }
                )
                NavigationBarItem(
                    selected = selectedTab == 3,
                    onClick = { selectedTab = 3 },
                    icon = { Icon(Icons.Default.Description, contentDescription = null) },
                    label = { Text(stringResource(R.string.generate_certificates), fontSize = 10.sp, maxLines = 1) }
                )
                NavigationBarItem(
                    selected = selectedTab == 4,
                    onClick = { selectedTab = 4 },
                    icon = { Icon(Icons.Default.Security, contentDescription = null) },
                    label = { Text(stringResource(R.string.audit_logs), fontSize = 10.sp, maxLines = 1) }
                )
                NavigationBarItem(
                    selected = selectedTab == 5,
                    onClick = { selectedTab = 5 },
                    icon = { Icon(Icons.Default.Settings, contentDescription = null) },
                    label = { Text(stringResource(R.string.profile_title), fontSize = 10.sp, maxLines = 1) }
                )
            }
        }
    ) { padding ->
        Box(
            modifier = Modifier
                .fillMaxSize()
                .padding(padding)
                .padding(horizontal = 16.dp, vertical = 12.dp)
        ) {
            when (selectedTab) {
                0 -> SecretaryOverviewTab(
                    viewModel = viewModel,
                    onNavigateToUsers = { selectedTab = 1 },
                    onNavigateToCertificates = { selectedTab = 3 }
                )
                1 -> SecretaryUsersTab(
                    viewModel = viewModel,
                    onAddUser = { showAddUserDialog = true },
                    onEditUser = { userToEdit = it },
                    onDeleteUser = { userToDelete = it }
                )
                2 -> SecretaryScrutinyTab(viewModel)
                3 -> SecretaryCertificatesTab(
                    viewModel = viewModel,
                    onGenerateCertificate = { showGenerateCertDialog = true },
                    onDownloadPdf = { cert ->
                        downloadedCertProtocol = cert.protocolNo
                    }
                )
                4 -> SecretaryAuditTab(
                    viewModel = viewModel,
                    onRefresh = {
                        if (!token.isNullOrBlank()) {
                            coroutineScope.launch { viewModel.loadFromDatabase(token) }
                        }
                    }
                )
                5 -> SecretaryProfileTab(
                    token = token ?: "",
                    viewModel = viewModel
                )
            }
        }
    }

    // Modal: Aggiungi Utente
    if (showAddUserDialog && !token.isNullOrBlank()) {
        AddUserDialog(
            onDismiss = { showAddUserDialog = false },
            onConfirm = { payload ->
                coroutineScope.launch {
                    val ok = viewModel.createUserOnline(token, payload)
                    if (ok) showAddUserDialog = false
                }
            }
        )
    }

    // Modal: Modifica Utente
    userToEdit?.let { user ->
        if (!token.isNullOrBlank()) {
            EditUserDialog(
                user = user,
                onDismiss = { userToEdit = null },
                onConfirm = { payload ->
                    coroutineScope.launch {
                        val ok = viewModel.updateUserOnline(token, user.id, payload)
                        if (ok) userToEdit = null
                    }
                }
            )
        }
    }

    // Modal: Conferma Eliminazione Utente
    userToDelete?.let { user ->
        if (!token.isNullOrBlank()) {
            AlertDialog(
                onDismissRequest = { userToDelete = null },
                title = { Text(stringResource(R.string.confirm_delete_title), fontWeight = FontWeight.Bold) },
                text = {
                    Text(
                        "${stringResource(R.string.confirm_delete_user_msg)}\n\nUtente: ${user.firstName} ${user.lastName} (${user.email})"
                    )
                },
                confirmButton = {
                    Button(
                        onClick = {
                            coroutineScope.launch {
                                viewModel.deleteUserOnline(token, user.id)
                                userToDelete = null
                            }
                        },
                        colors = ButtonDefaults.buttonColors(containerColor = MaterialTheme.colorScheme.error)
                    ) {
                        Text(stringResource(R.string.btn_confirm))
                    }
                },
                dismissButton = {
                    TextButton(onClick = { userToDelete = null }) {
                        Text(stringResource(R.string.btn_cancel))
                    }
                }
            )
        }
    }

    // Modal: Genera Certificato
    if (showGenerateCertDialog && !token.isNullOrBlank()) {
        GenerateCertificateDialog(
            students = viewModel.usersList.filter { it.role.equals("student", ignoreCase = true) },
            onDismiss = { showGenerateCertDialog = false },
            onConfirm = { studentId, type, year, notes ->
                coroutineScope.launch {
                    val ok = viewModel.generateCertificateOnline(token, studentId, type, year, notes)
                    if (ok) showGenerateCertDialog = false
                }
            }
        )
    }

    // Dialog: Feedback Download PDF
    downloadedCertProtocol?.let { prot ->
        AlertDialog(
            onDismissRequest = { downloadedCertProtocol = null },
            title = { Text("Download Certificato", fontWeight = FontWeight.Bold) },
            text = { Text("Il certificato protocollato $prot è pronto per il download e la stampa ufficiale.") },
            confirmButton = {
                Button(onClick = { downloadedCertProtocol = null }) {
                    Text("OK")
                }
            }
        )
    }
}

// -------------------------------------------------------------
// TAB 0: PANORAMICA
// -------------------------------------------------------------
@Composable
fun SecretaryOverviewTab(
    viewModel: SecretaryViewModel,
    onNavigateToUsers: () -> Unit,
    onNavigateToCertificates: () -> Unit
) {
    val stats = viewModel.stats
    val studentsCount = if (stats.totalStudents > 0) stats.totalStudents else viewModel.getUsersByRole("student").size
    val teachersCount = if (stats.totalTeachers > 0) stats.totalTeachers else viewModel.getUsersByRole("teacher").size
    val usersCount = if (stats.totalUsers > 0) stats.totalUsers else viewModel.usersList.size
    val certsCount = viewModel.certificateRequests.size

    LazyColumn(verticalArrangement = Arrangement.spacedBy(16.dp)) {
        item {
            Row(horizontalArrangement = Arrangement.spacedBy(12.dp)) {
                MetricCard(
                    modifier = Modifier.weight(1f),
                    title = stringResource(R.string.total_students),
                    value = studentsCount.toString(),
                    icon = Icons.Default.School,
                    color = Color(0xFF2563EB)
                )
                MetricCard(
                    modifier = Modifier.weight(1f),
                    title = stringResource(R.string.total_teachers),
                    value = teachersCount.toString(),
                    icon = Icons.Default.People,
                    color = Color(0xFF0D9488)
                )
            }
        }

        item {
            Row(horizontalArrangement = Arrangement.spacedBy(12.dp)) {
                MetricCard(
                    modifier = Modifier.weight(1f),
                    title = "Utenti Totali",
                    value = usersCount.toString(),
                    icon = Icons.Default.Badge,
                    color = Color(0xFF581C87)
                )
                MetricCard(
                    modifier = Modifier.weight(1f),
                    title = stringResource(R.string.certificate_requests),
                    value = certsCount.toString(),
                    icon = Icons.Default.Description,
                    color = Color(0xFFEA580C)
                )
            }
        }

        item {
            Card(
                modifier = Modifier.fillMaxWidth(),
                shape = RoundedCornerShape(16.dp),
                colors = CardDefaults.cardColors(containerColor = MaterialTheme.colorScheme.surfaceVariant)
            ) {
                Column(modifier = Modifier.padding(18.dp)) {
                    Text(
                        stringResource(R.string.secretary_dashboard_title),
                        fontWeight = FontWeight.Bold,
                        fontSize = 16.sp
                    )
                    Spacer(modifier = Modifier.height(6.dp))
                    Text(
                        "Pannello amministrativo scolastico sincronizzato con il server centrale.",
                        fontSize = 13.sp,
                        color = Color.Gray
                    )
                    Spacer(modifier = Modifier.height(14.dp))
                    Row(
                        modifier = Modifier.fillMaxWidth(),
                        horizontalArrangement = Arrangement.spacedBy(10.dp)
                    ) {
                        Button(
                            onClick = onNavigateToUsers,
                            modifier = Modifier.weight(1f),
                            shape = RoundedCornerShape(10.dp),
                            colors = ButtonDefaults.buttonColors(containerColor = Color(0xFF581C87))
                        ) {
                            Icon(Icons.Default.PersonAdd, contentDescription = null, modifier = Modifier.size(16.dp))
                            Spacer(modifier = Modifier.width(6.dp))
                            Text("Gestione Utenti", fontSize = 12.sp)
                        }
                        OutlinedButton(
                            onClick = onNavigateToCertificates,
                            modifier = Modifier.weight(1f),
                            shape = RoundedCornerShape(10.dp)
                        ) {
                            Icon(Icons.Default.PostAdd, contentDescription = null, modifier = Modifier.size(16.dp))
                            Spacer(modifier = Modifier.width(6.dp))
                            Text("Certificati", fontSize = 12.sp)
                        }
                    }
                }
            }
        }
    }
}

// -------------------------------------------------------------
// TAB 1: GESTIONE UTENTI
// -------------------------------------------------------------
@Composable
fun SecretaryUsersTab(
    viewModel: SecretaryViewModel,
    onAddUser: () -> Unit,
    onEditUser: (ManagedUser) -> Unit,
    onDeleteUser: (ManagedUser) -> Unit
) {
    var searchQuery by remember { mutableStateOf("") }
    var selectedRoleFilter by remember { mutableStateOf("all") }

    val filteredUsers = viewModel.usersList.filter { user ->
        val matchesQuery = searchQuery.isBlank() ||
                user.firstName.contains(searchQuery, ignoreCase = true) ||
                user.lastName.contains(searchQuery, ignoreCase = true) ||
                user.email.contains(searchQuery, ignoreCase = true)

        val matchesRole = when (selectedRoleFilter) {
            "all" -> true
            "student" -> user.role.equals("student", ignoreCase = true)
            "teacher" -> user.role.equals("teacher", ignoreCase = true)
            "parent" -> user.role.equals("parent", ignoreCase = true)
            "secretary" -> user.role.equals("secretary", ignoreCase = true) || user.role.equals("admin", ignoreCase = true)
            else -> true
        }

        matchesQuery && matchesRole
    }

    LazyColumn(verticalArrangement = Arrangement.spacedBy(12.dp)) {
        item {
            Row(
                modifier = Modifier.fillMaxWidth(),
                horizontalArrangement = Arrangement.SpaceBetween,
                verticalAlignment = Alignment.CenterVertically
            ) {
                Text(
                    stringResource(R.string.user_management),
                    fontWeight = FontWeight.Bold,
                    fontSize = 18.sp
                )
                Button(
                    onClick = onAddUser,
                    shape = RoundedCornerShape(10.dp),
                    colors = ButtonDefaults.buttonColors(containerColor = Color(0xFF581C87))
                ) {
                    Icon(Icons.Default.PersonAdd, contentDescription = null, modifier = Modifier.size(16.dp))
                    Spacer(modifier = Modifier.width(6.dp))
                    Text(stringResource(R.string.add_user_button), fontSize = 13.sp)
                }
            }
        }

        // Search Bar
        item {
            OutlinedTextField(
                value = searchQuery,
                onValueChange = { searchQuery = it },
                modifier = Modifier.fillMaxWidth(),
                placeholder = { Text(stringResource(R.string.search_users_hint), fontSize = 13.sp) },
                leadingIcon = { Icon(Icons.Default.Search, contentDescription = null, tint = Color.Gray) },
                trailingIcon = {
                    if (searchQuery.isNotEmpty()) {
                        IconButton(onClick = { searchQuery = "" }) {
                            Icon(Icons.Default.Close, contentDescription = null, modifier = Modifier.size(16.dp))
                        }
                    }
                },
                singleLine = true,
                shape = RoundedCornerShape(12.dp)
            )
        }

        // Filter Chips
        item {
            LazyRow(horizontalArrangement = Arrangement.spacedBy(8.dp)) {
                val filters = listOf(
                    "all" to R.string.filter_all,
                    "student" to R.string.filter_students,
                    "teacher" to R.string.filter_teachers,
                    "parent" to R.string.filter_parents,
                    "secretary" to R.string.filter_secretary
                )
                items(filters) { (roleKey, strRes) ->
                    FilterChip(
                        selected = selectedRoleFilter == roleKey,
                        onClick = { selectedRoleFilter = roleKey },
                        label = { Text(stringResource(strRes), fontSize = 12.sp) }
                    )
                }
            }
        }

        if (filteredUsers.isEmpty()) {
            item {
                Box(
                    modifier = Modifier
                        .fillMaxWidth()
                        .padding(vertical = 32.dp),
                    contentAlignment = Alignment.Center
                ) {
                    Column(horizontalAlignment = Alignment.CenterHorizontally) {
                        Icon(Icons.Default.PersonOff, contentDescription = null, modifier = Modifier.size(48.dp), tint = Color.LightGray)
                        Spacer(modifier = Modifier.height(8.dp))
                        Text(stringResource(R.string.no_users_found), color = Color.Gray, fontSize = 14.sp)
                    }
                }
            }
        } else {
            items(filteredUsers, key = { it.id }) { user ->
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
                        Row(verticalAlignment = Alignment.CenterVertically, modifier = Modifier.weight(1f)) {
                            // Avatar
                            Box(
                                modifier = Modifier
                                    .size(40.dp)
                                    .background(getRoleColor(user.role).copy(alpha = 0.15f), CircleShape),
                                contentAlignment = Alignment.Center
                            ) {
                                Text(
                                    text = "${user.firstName.take(1)}${user.lastName.take(1)}".uppercase(),
                                    fontWeight = FontWeight.Bold,
                                    color = getRoleColor(user.role),
                                    fontSize = 14.sp
                                )
                            }
                            Spacer(modifier = Modifier.width(12.dp))
                            Column {
                                Text("${user.firstName} ${user.lastName}", fontWeight = FontWeight.Bold, fontSize = 15.sp)
                                Spacer(modifier = Modifier.height(2.dp))
                                Row(verticalAlignment = Alignment.CenterVertically) {
                                    Surface(
                                        color = getRoleColor(user.role),
                                        shape = RoundedCornerShape(4.dp)
                                    ) {
                                        Text(
                                            formatRole(user.role),
                                            color = Color.White,
                                            fontSize = 10.sp,
                                            fontWeight = FontWeight.SemiBold,
                                            modifier = Modifier.padding(horizontal = 6.dp, vertical = 2.dp)
                                        )
                                    }
                                    Spacer(modifier = Modifier.width(6.dp))
                                    Text(user.email, fontSize = 12.sp, color = Color.Gray)
                                }
                            }
                        }

                        // Actions
                        Row {
                            IconButton(onClick = { onEditUser(user) }) {
                                Icon(Icons.Default.Edit, contentDescription = stringResource(R.string.edit_user), tint = Color(0xFF581C87))
                            }
                            IconButton(onClick = { onDeleteUser(user) }) {
                                Icon(Icons.Default.Delete, contentDescription = stringResource(R.string.delete_user), tint = MaterialTheme.colorScheme.error)
                            }
                        }
                    }
                }
            }
        }
    }
}

// -------------------------------------------------------------
// TAB 2: SUPERVISIONE SCRUTINI
// -------------------------------------------------------------
@Composable
fun SecretaryScrutinyTab(viewModel: SecretaryViewModel) {
    val classes = viewModel.scrutinyClasses

    LazyColumn(verticalArrangement = Arrangement.spacedBy(12.dp)) {
        item {
            Text(stringResource(R.string.scrutiny_supervision), fontWeight = FontWeight.Bold, fontSize = 18.sp)
        }
        if (classes.isEmpty()) {
            item {
                Text("Nessuna classe disponibile per gli scrutini", color = Color.Gray, fontSize = 13.sp)
            }
        } else {
            items(classes) { cls ->
                Card(modifier = Modifier.fillMaxWidth(), shape = RoundedCornerShape(12.dp)) {
                    Row(
                        modifier = Modifier
                            .fillMaxWidth()
                            .padding(14.dp),
                        horizontalArrangement = Arrangement.SpaceBetween,
                        verticalAlignment = Alignment.CenterVertically
                    ) {
                        Column {
                            Text(cls.className, fontWeight = FontWeight.Bold, fontSize = 15.sp)
                            Text("Periodo: ${cls.period}° Quadrimestre", fontSize = 12.sp, color = Color.Gray)
                        }
                        Button(
                            onClick = { viewModel.toggleClassScrutinyLock(cls.classId) },
                            colors = ButtonDefaults.buttonColors(
                                containerColor = if (cls.isLocked) Color(0xFF10B981) else Color(0xFFEA580C)
                            ),
                            shape = RoundedCornerShape(8.dp)
                        ) {
                            Text(if (cls.isLocked) "Chiuso / Bloccato" else "Aperto / In Corso", fontSize = 12.sp)
                        }
                    }
                }
            }
        }
    }
}

// -------------------------------------------------------------
// TAB 3: CERTIFICATI & ATTI
// -------------------------------------------------------------
@Composable
fun SecretaryCertificatesTab(
    viewModel: SecretaryViewModel,
    onGenerateCertificate: () -> Unit,
    onDownloadPdf: (it.scuola.registro.secretary.data.CertificateRequest) -> Unit
) {
    val certs = viewModel.certificateRequests

    LazyColumn(verticalArrangement = Arrangement.spacedBy(12.dp)) {
        item {
            Row(
                modifier = Modifier.fillMaxWidth(),
                horizontalArrangement = Arrangement.SpaceBetween,
                verticalAlignment = Alignment.CenterVertically
            ) {
                Text(stringResource(R.string.generate_certificates), fontWeight = FontWeight.Bold, fontSize = 18.sp)
                Button(
                    onClick = onGenerateCertificate,
                    shape = RoundedCornerShape(10.dp),
                    colors = ButtonDefaults.buttonColors(containerColor = Color(0xFF581C87))
                ) {
                    Icon(Icons.Default.Add, contentDescription = null, modifier = Modifier.size(16.dp))
                    Spacer(modifier = Modifier.width(4.dp))
                    Text(stringResource(R.string.new_certificate_btn), fontSize = 13.sp)
                }
            }
        }

        if (certs.isEmpty()) {
            item {
                Box(
                    modifier = Modifier
                        .fillMaxWidth()
                        .padding(vertical = 40.dp),
                    contentAlignment = Alignment.Center
                ) {
                    Column(horizontalAlignment = Alignment.CenterHorizontally) {
                        Icon(Icons.Default.FolderOpen, contentDescription = null, modifier = Modifier.size(48.dp), tint = Color.LightGray)
                        Spacer(modifier = Modifier.height(8.dp))
                        Text(stringResource(R.string.no_certificates_found), color = Color.Gray, fontSize = 14.sp)
                    }
                }
            }
        } else {
            items(certs) { cert ->
                Card(modifier = Modifier.fillMaxWidth(), shape = RoundedCornerShape(12.dp)) {
                    Column(modifier = Modifier.padding(14.dp), verticalArrangement = Arrangement.spacedBy(6.dp)) {
                        Row(
                            modifier = Modifier.fillMaxWidth(),
                            horizontalArrangement = Arrangement.SpaceBetween,
                            verticalAlignment = Alignment.CenterVertically
                        ) {
                            Text(cert.certificateType, fontWeight = FontWeight.Bold, fontSize = 15.sp)
                            Surface(color = Color(0xFF10B981), shape = RoundedCornerShape(6.dp)) {
                                Text(
                                    "Sigillo Digitale Attivo",
                                    color = Color.White,
                                    modifier = Modifier.padding(horizontal = 6.dp, vertical = 2.dp),
                                    fontSize = 10.sp
                                )
                            }
                        }
                        if (cert.studentName.isNotBlank()) {
                            Text(
                                "Studente: ${cert.studentName} ${if (cert.className.isNotBlank()) "(${cert.className})" else ""}",
                                fontSize = 13.sp,
                                fontWeight = FontWeight.Medium,
                                color = MaterialTheme.colorScheme.primary
                            )
                        }
                        Row(
                            modifier = Modifier.fillMaxWidth(),
                            horizontalArrangement = Arrangement.SpaceBetween,
                            verticalAlignment = Alignment.CenterVertically
                        ) {
                            Text(
                                "${stringResource(R.string.protocol_prefix)} ${cert.protocolNo} • ${cert.issuedAt}",
                                fontSize = 12.sp,
                                color = Color.Gray
                            )
                            OutlinedButton(
                                onClick = { onDownloadPdf(cert) },
                                shape = RoundedCornerShape(8.dp),
                                contentPadding = PaddingValues(horizontal = 12.dp, vertical = 4.dp)
                            ) {
                                Icon(Icons.Default.Download, contentDescription = null, modifier = Modifier.size(14.dp))
                                Spacer(modifier = Modifier.width(4.dp))
                                Text("PDF", fontSize = 12.sp)
                            }
                        }
                    }
                }
            }
        }
    }
}

// -------------------------------------------------------------
// TAB 4: AUDIT LOGS
// -------------------------------------------------------------
@Composable
fun SecretaryAuditTab(
    viewModel: SecretaryViewModel,
    onRefresh: () -> Unit
) {
    val logs = viewModel.auditLogs

    LazyColumn(verticalArrangement = Arrangement.spacedBy(10.dp)) {
        item {
            Row(
                modifier = Modifier.fillMaxWidth(),
                horizontalArrangement = Arrangement.SpaceBetween,
                verticalAlignment = Alignment.CenterVertically
            ) {
                Text(stringResource(R.string.audit_logs), fontWeight = FontWeight.Bold, fontSize = 18.sp)
                IconButton(onClick = onRefresh) {
                    Icon(Icons.Default.Refresh, contentDescription = stringResource(R.string.btn_refresh), tint = Color(0xFF581C87))
                }
            }
        }

        if (logs.isEmpty()) {
            item {
                Box(
                    modifier = Modifier
                        .fillMaxWidth()
                        .padding(vertical = 40.dp),
                    contentAlignment = Alignment.Center
                ) {
                    Column(horizontalAlignment = Alignment.CenterHorizontally) {
                        Icon(Icons.Default.Security, contentDescription = null, modifier = Modifier.size(48.dp), tint = Color.LightGray)
                        Spacer(modifier = Modifier.height(8.dp))
                        Text(stringResource(R.string.no_audit_logs), color = Color.Gray, fontSize = 14.sp)
                    }
                }
            }
        } else {
            items(logs) { log ->
                Card(modifier = Modifier.fillMaxWidth(), shape = RoundedCornerShape(10.dp)) {
                    Column(modifier = Modifier.padding(12.dp), verticalArrangement = Arrangement.spacedBy(4.dp)) {
                        Row(
                            modifier = Modifier.fillMaxWidth(),
                            horizontalArrangement = Arrangement.SpaceBetween,
                            verticalAlignment = Alignment.CenterVertically
                        ) {
                            Row(verticalAlignment = Alignment.CenterVertically) {
                                Surface(
                                    color = Color(0xFF581C87),
                                    shape = RoundedCornerShape(4.dp)
                                ) {
                                    Text(
                                        log.action,
                                        color = Color.White,
                                        fontSize = 11.sp,
                                        fontWeight = FontWeight.Bold,
                                        modifier = Modifier.padding(horizontal = 6.dp, vertical = 2.dp)
                                    )
                                }
                                Spacer(modifier = Modifier.width(8.dp))
                                Text(
                                    "${log.actorName} (${log.actorRole})",
                                    fontWeight = FontWeight.SemiBold,
                                    fontSize = 13.sp
                                )
                            }
                            Text(log.createdAt, fontSize = 11.sp, color = Color.Gray)
                        }
                        if (log.details.isNotBlank()) {
                            Text(log.details, fontSize = 12.sp, color = MaterialTheme.colorScheme.onSurfaceVariant)
                        }
                        Text("IP: ${log.ipAddress}", fontSize = 11.sp, color = Color.LightGray)
                    }
                }
            }
        }
    }
}

// -------------------------------------------------------------
// TAB 5: PROFILO & CAMBIO PASSWORD
// -------------------------------------------------------------
@Composable
fun SecretaryProfileTab(
    token: String,
    viewModel: SecretaryViewModel
) {
    val coroutineScope = rememberCoroutineScope()
    val profile = viewModel.currentUserProfile

    var currentPassword by remember { mutableStateOf("") }
    var newPassword by remember { mutableStateOf("") }
    var confirmPassword by remember { mutableStateOf("") }
    var showPassword by remember { mutableStateOf(false) }
    var isChangingPassword by remember { mutableStateOf(false) }
    var localError by remember { mutableStateOf<String?>(null) }
    var localSuccess by remember { mutableStateOf<String?>(null) }

    LazyColumn(verticalArrangement = Arrangement.spacedBy(16.dp)) {
        // Dati Utente
        item {
            Card(modifier = Modifier.fillMaxWidth(), shape = RoundedCornerShape(14.dp)) {
                Column(modifier = Modifier.padding(16.dp), verticalArrangement = Arrangement.spacedBy(8.dp)) {
                    Text(stringResource(R.string.my_account_title), fontWeight = FontWeight.Bold, fontSize = 16.sp)
                    HorizontalDivider()
                    Text("Nome: ${profile?.firstName ?: "Personale"} ${profile?.lastName ?: "Segreteria"}", fontWeight = FontWeight.Medium)
                    Text("Email: ${profile?.email ?: "segreteria@scuola.it"}", color = Color.Gray, fontSize = 13.sp)
                    Text("Ruolo: ${profile?.role ?: "Segreteria"}", color = Color.Gray, fontSize = 13.sp)
                }
            }
        }

        // Cambio Password
        item {
            Card(modifier = Modifier.fillMaxWidth(), shape = RoundedCornerShape(14.dp)) {
                Column(modifier = Modifier.padding(16.dp), verticalArrangement = Arrangement.spacedBy(12.dp)) {
                    Text(stringResource(R.string.security_settings_title), fontWeight = FontWeight.Bold, fontSize = 16.sp)
                    Text("Modifica la password del tuo account personale.", fontSize = 12.sp, color = Color.Gray)

                    OutlinedTextField(
                        value = currentPassword,
                        onValueChange = { currentPassword = it; localError = null },
                        label = { Text(stringResource(R.string.current_password_label)) },
                        visualTransformation = if (showPassword) VisualTransformation.None else PasswordVisualTransformation(),
                        singleLine = true,
                        modifier = Modifier.fillMaxWidth(),
                        shape = RoundedCornerShape(10.dp)
                    )

                    OutlinedTextField(
                        value = newPassword,
                        onValueChange = { newPassword = it; localError = null },
                        label = { Text(stringResource(R.string.new_password_label)) },
                        visualTransformation = if (showPassword) VisualTransformation.None else PasswordVisualTransformation(),
                        singleLine = true,
                        modifier = Modifier.fillMaxWidth(),
                        shape = RoundedCornerShape(10.dp)
                    )

                    OutlinedTextField(
                        value = confirmPassword,
                        onValueChange = { confirmPassword = it; localError = null },
                        label = { Text(stringResource(R.string.confirm_password_label)) },
                        visualTransformation = if (showPassword) VisualTransformation.None else PasswordVisualTransformation(),
                        singleLine = true,
                        modifier = Modifier.fillMaxWidth(),
                        shape = RoundedCornerShape(10.dp)
                    )

                    Row(
                        modifier = Modifier.fillMaxWidth(),
                        horizontalArrangement = Arrangement.SpaceBetween,
                        verticalAlignment = Alignment.CenterVertically
                    ) {
                        TextButton(onClick = { showPassword = !showPassword }) {
                            Text(if (showPassword) "Nascondi password" else "Mostra password", fontSize = 12.sp)
                        }
                    }

                    if (localError != null) {
                        Text(localError!!, color = MaterialTheme.colorScheme.error, fontSize = 12.sp)
                    }
                    if (localSuccess != null) {
                        Text(localSuccess!!, color = Color(0xFF10B981), fontSize = 12.sp, fontWeight = FontWeight.SemiBold)
                    }

                    Button(
                        onClick = {
                            if (currentPassword.isBlank() || newPassword.isBlank() || confirmPassword.isBlank()) {
                                localError = "Compilare tutti i campi"
                                return@Button
                            }
                            if (newPassword.length < 10) {
                                localError = "La nuova password deve contenere almeno 10 caratteri"
                                return@Button
                            }
                            if (newPassword != confirmPassword) {
                                localError = "Le password non coincidono"
                                return@Button
                            }

                            val userId = profile?.id ?: "me"
                            isChangingPassword = true
                            localError = null
                            localSuccess = null

                            coroutineScope.launch {
                                val res = viewModel.changePasswordOnline(token, userId, currentPassword, newPassword)
                                isChangingPassword = false
                                if (res.isSuccess) {
                                    localSuccess = "Password aggiornata con successo"
                                    currentPassword = ""
                                    newPassword = ""
                                    confirmPassword = ""
                                } else {
                                    localError = res.exceptionOrNull()?.message ?: "Errore aggiornamento password"
                                }
                            }
                        },
                        modifier = Modifier.fillMaxWidth(),
                        shape = RoundedCornerShape(10.dp),
                        colors = ButtonDefaults.buttonColors(containerColor = Color(0xFF581C87)),
                        enabled = !isChangingPassword
                    ) {
                        if (isChangingPassword) {
                            CircularProgressIndicator(color = Color.White, modifier = Modifier.size(20.dp))
                        } else {
                            Text(stringResource(R.string.btn_change_password))
                        }
                    }
                }
            }
        }

        // Info App
        item {
            Card(modifier = Modifier.fillMaxWidth(), shape = RoundedCornerShape(14.dp)) {
                Column(modifier = Modifier.padding(16.dp)) {
                    Text(stringResource(R.string.school_info_title), fontWeight = FontWeight.Bold, fontSize = 15.sp)
                    Spacer(modifier = Modifier.height(4.dp))
                    Text(stringResource(R.string.app_version_label), fontSize = 12.sp, color = Color.Gray)
                }
            }
        }
    }
}

// -------------------------------------------------------------
// DIALOG: AGGIUNGI UTENTE
// -------------------------------------------------------------
@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun AddUserDialog(
    onDismiss: () -> Unit,
    onConfirm: (CreateUserPayload) -> Unit
) {
    var firstName by remember { mutableStateOf("") }
    var lastName by remember { mutableStateOf("") }
    var email by remember { mutableStateOf("") }
    var password by remember { mutableStateOf("") }
    var role by remember { mutableStateOf("teacher") }
    var fiscalCode by remember { mutableStateOf("") }
    var phoneNumber by remember { mutableStateOf("") }
    var roleExpanded by remember { mutableStateOf(false) }
    var errorMsg by remember { mutableStateOf<String?>(null) }

    val roles = listOf(
        "teacher" to "Docente",
        "student" to "Studente",
        "parent" to "Genitore",
        "secretary" to "Segreteria"
    )

    AlertDialog(
        onDismissRequest = onDismiss,
        title = { Text(stringResource(R.string.add_user_button), fontWeight = FontWeight.Bold) },
        text = {
            Column(
                modifier = Modifier.fillMaxWidth(),
                verticalArrangement = Arrangement.spacedBy(10.dp)
            ) {
                OutlinedTextField(
                    value = firstName,
                    onValueChange = { firstName = it; errorMsg = null },
                    label = { Text("Nome *") },
                    singleLine = true,
                    modifier = Modifier.fillMaxWidth()
                )
                OutlinedTextField(
                    value = lastName,
                    onValueChange = { lastName = it; errorMsg = null },
                    label = { Text("Cognome *") },
                    singleLine = true,
                    modifier = Modifier.fillMaxWidth()
                )
                OutlinedTextField(
                    value = email,
                    onValueChange = { email = it; errorMsg = null },
                    label = { Text("Email *") },
                    singleLine = true,
                    modifier = Modifier.fillMaxWidth()
                )
                OutlinedTextField(
                    value = password,
                    onValueChange = { password = it; errorMsg = null },
                    label = { Text("Password temporanea (min. 10 car.) *") },
                    singleLine = true,
                    modifier = Modifier.fillMaxWidth()
                )
                // Role Picker
                ExposedDropdownMenuBox(
                    expanded = roleExpanded,
                    onExpandedChange = { roleExpanded = !roleExpanded }
                ) {
                    OutlinedTextField(
                        value = roles.find { it.first == role }?.second ?: "Docente",
                        onValueChange = {},
                        readOnly = true,
                        label = { Text(stringResource(R.string.role_label)) },
                        trailingIcon = { ExposedDropdownMenuDefaults.TrailingIcon(expanded = roleExpanded) },
                        modifier = Modifier.menuAnchor().fillMaxWidth()
                    )
                    ExposedDropdownMenu(
                        expanded = roleExpanded,
                        onDismissRequest = { roleExpanded = false }
                    ) {
                        roles.forEach { (key, label) ->
                            DropdownMenuItem(
                                text = { Text(label) },
                                onClick = {
                                    role = key
                                    roleExpanded = false
                                }
                            )
                        }
                    }
                }
                OutlinedTextField(
                    value = fiscalCode,
                    onValueChange = { fiscalCode = it },
                    label = { Text(stringResource(R.string.fiscal_code_label)) },
                    singleLine = true,
                    modifier = Modifier.fillMaxWidth()
                )

                if (errorMsg != null) {
                    Text(errorMsg!!, color = MaterialTheme.colorScheme.error, fontSize = 12.sp)
                }
            }
        },
        confirmButton = {
            Button(onClick = {
                if (firstName.isBlank() || lastName.isBlank() || email.isBlank() || password.isBlank()) {
                    errorMsg = "Compilare tutti i campi obbligatori (*)"
                    return@Button
                }
                if (password.length < 10) {
                    errorMsg = "La password deve avere almeno 10 caratteri"
                    return@Button
                }
                onConfirm(
                    CreateUserPayload(
                        firstName = firstName.trim(),
                        lastName = lastName.trim(),
                        email = email.trim(),
                        password = password,
                        role = role,
                        fiscalCode = fiscalCode.trim().uppercase(),
                        phoneNumber = phoneNumber.trim()
                    )
                )
            }) {
                Text(stringResource(R.string.btn_save))
            }
        },
        dismissButton = {
            TextButton(onClick = onDismiss) {
                Text(stringResource(R.string.btn_cancel))
            }
        }
    )
}

// -------------------------------------------------------------
// DIALOG: MODIFICA UTENTE
// -------------------------------------------------------------
@Composable
fun EditUserDialog(
    user: ManagedUser,
    onDismiss: () -> Unit,
    onConfirm: (UpdateUserPayload) -> Unit
) {
    var firstName by remember { mutableStateOf(user.firstName) }
    var lastName by remember { mutableStateOf(user.lastName) }
    var fiscalCode by remember { mutableStateOf(user.fiscalCode) }
    var phoneNumber by remember { mutableStateOf(user.phoneNumber) }

    AlertDialog(
        onDismissRequest = onDismiss,
        title = { Text(stringResource(R.string.edit_user), fontWeight = FontWeight.Bold) },
        text = {
            Column(
                modifier = Modifier.fillMaxWidth(),
                verticalArrangement = Arrangement.spacedBy(10.dp)
            ) {
                OutlinedTextField(
                    value = firstName,
                    onValueChange = { firstName = it },
                    label = { Text("Nome") },
                    singleLine = true,
                    modifier = Modifier.fillMaxWidth()
                )
                OutlinedTextField(
                    value = lastName,
                    onValueChange = { lastName = it },
                    label = { Text("Cognome") },
                    singleLine = true,
                    modifier = Modifier.fillMaxWidth()
                )
                OutlinedTextField(
                    value = fiscalCode,
                    onValueChange = { fiscalCode = it },
                    label = { Text(stringResource(R.string.fiscal_code_label)) },
                    singleLine = true,
                    modifier = Modifier.fillMaxWidth()
                )
                OutlinedTextField(
                    value = phoneNumber,
                    onValueChange = { phoneNumber = it },
                    label = { Text(stringResource(R.string.phone_label)) },
                    singleLine = true,
                    modifier = Modifier.fillMaxWidth()
                )
            }
        },
        confirmButton = {
            Button(onClick = {
                onConfirm(
                    UpdateUserPayload(
                        firstName = firstName.trim(),
                        lastName = lastName.trim(),
                        fiscalCode = fiscalCode.trim().uppercase(),
                        phoneNumber = phoneNumber.trim()
                    )
                )
            }) {
                Text(stringResource(R.string.btn_save))
            }
        },
        dismissButton = {
            TextButton(onClick = onDismiss) {
                Text(stringResource(R.string.btn_cancel))
            }
        }
    )
}

// -------------------------------------------------------------
// DIALOG: GENERA CERTIFICATO
// -------------------------------------------------------------
@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun GenerateCertificateDialog(
    students: List<ManagedUser>,
    onDismiss: () -> Unit,
    onConfirm: (studentId: String, type: String, academicYear: String, notes: String) -> Unit
) {
    var selectedStudentId by remember { mutableStateOf(students.firstOrNull()?.id ?: "") }
    var studentNameInput by remember { mutableStateOf("") }
    var selectedType by remember { mutableStateOf("iscrizione") }
    var academicYear by remember { mutableStateOf("2025/2026") }
    var notes by remember { mutableStateOf("") }
    var typeExpanded by remember { mutableStateOf(false) }
    var studentExpanded by remember { mutableStateOf(false) }
    var errorMsg by remember { mutableStateOf<String?>(null) }

    val certTypes = listOf(
        "iscrizione" to "Certificato di Iscrizione",
        "frequenza" to "Certificato di Frequenza",
        "promozione" to "Certificato di Promozione",
        "condotta" to "Certificato di Buona Condotta"
    )

    AlertDialog(
        onDismissRequest = onDismiss,
        title = { Text(stringResource(R.string.new_certificate_dialog_title), fontWeight = FontWeight.Bold) },
        text = {
            Column(
                modifier = Modifier.fillMaxWidth(),
                verticalArrangement = Arrangement.spacedBy(10.dp)
            ) {
                if (students.isNotEmpty()) {
                    ExposedDropdownMenuBox(
                        expanded = studentExpanded,
                        onExpandedChange = { studentExpanded = !studentExpanded }
                    ) {
                        val currentStudent = students.find { it.id == selectedStudentId }
                        OutlinedTextField(
                            value = currentStudent?.let { "${it.firstName} ${it.lastName}" } ?: "Seleziona studente",
                            onValueChange = {},
                            readOnly = true,
                            label = { Text(stringResource(R.string.select_student_label)) },
                            trailingIcon = { ExposedDropdownMenuDefaults.TrailingIcon(expanded = studentExpanded) },
                            modifier = Modifier.menuAnchor().fillMaxWidth()
                        )
                        ExposedDropdownMenu(
                            expanded = studentExpanded,
                            onDismissRequest = { studentExpanded = false }
                        ) {
                            students.forEach { st ->
                                DropdownMenuItem(
                                    text = { Text("${st.firstName} ${st.lastName} (${st.email})") },
                                    onClick = {
                                        selectedStudentId = st.id
                                        studentExpanded = false
                                    }
                                )
                            }
                        }
                    }
                } else {
                    OutlinedTextField(
                        value = studentNameInput,
                        onValueChange = { studentNameInput = it },
                        label = { Text("ID Studente o Nome") },
                        singleLine = true,
                        modifier = Modifier.fillMaxWidth()
                    )
                }

                // Certificate Type
                ExposedDropdownMenuBox(
                    expanded = typeExpanded,
                    onExpandedChange = { typeExpanded = !typeExpanded }
                ) {
                    OutlinedTextField(
                        value = certTypes.find { it.first == selectedType }?.second ?: "",
                        onValueChange = {},
                        readOnly = true,
                        label = { Text(stringResource(R.string.certificate_type_label)) },
                        trailingIcon = { ExposedDropdownMenuDefaults.TrailingIcon(expanded = typeExpanded) },
                        modifier = Modifier.menuAnchor().fillMaxWidth()
                    )
                    ExposedDropdownMenu(
                        expanded = typeExpanded,
                        onDismissRequest = { typeExpanded = false }
                    ) {
                        certTypes.forEach { (key, label) ->
                            DropdownMenuItem(
                                text = { Text(label) },
                                onClick = {
                                    selectedType = key
                                    typeExpanded = false
                                }
                            )
                        }
                    }
                }

                OutlinedTextField(
                    value = academicYear,
                    onValueChange = { academicYear = it },
                    label = { Text(stringResource(R.string.academic_year_label)) },
                    singleLine = true,
                    modifier = Modifier.fillMaxWidth()
                )

                OutlinedTextField(
                    value = notes,
                    onValueChange = { notes = it },
                    label = { Text(stringResource(R.string.notes_label)) },
                    singleLine = true,
                    modifier = Modifier.fillMaxWidth()
                )

                if (errorMsg != null) {
                    Text(errorMsg!!, color = MaterialTheme.colorScheme.error, fontSize = 12.sp)
                }
            }
        },
        confirmButton = {
            Button(onClick = {
                val studentId = if (students.isNotEmpty()) selectedStudentId else studentNameInput.trim()
                if (studentId.isBlank()) {
                    errorMsg = "Seleziona o inserisci uno studente"
                    return@Button
                }
                onConfirm(studentId, selectedType, academicYear, notes)
            }) {
                Text(stringResource(R.string.btn_generate))
            }
        },
        dismissButton = {
            TextButton(onClick = onDismiss) {
                Text(stringResource(R.string.btn_cancel))
            }
        }
    )
}

// -------------------------------------------------------------
// UTILS & COMPONENTS
// -------------------------------------------------------------
@Composable
fun MetricCard(
    modifier: Modifier = Modifier,
    title: String,
    value: String,
    icon: ImageVector,
    color: Color
) {
    Card(
        modifier = modifier,
        shape = RoundedCornerShape(14.dp),
        colors = CardDefaults.cardColors(containerColor = MaterialTheme.colorScheme.surfaceVariant)
    ) {
        Column(modifier = Modifier.padding(16.dp)) {
            Icon(icon, contentDescription = null, tint = color, modifier = Modifier.size(28.dp))
            Spacer(modifier = Modifier.height(8.dp))
            Text(value, fontSize = 24.sp, fontWeight = FontWeight.Bold)
            Text(title, fontSize = 12.sp, color = Color.Gray, maxLines = 1)
        }
    }
}

fun getRoleColor(role: String): Color = when (role.lowercase()) {
    "student" -> Color(0xFF2563EB)
    "teacher", "coordinator" -> Color(0xFF0D9488)
    "parent" -> Color(0xFFEA580C)
    "secretary", "admin", "superadmin" -> Color(0xFF581C87)
    else -> Color(0xFF6B7280)
}

fun formatRole(role: String): String = when (role.lowercase()) {
    "student" -> "Studente"
    "teacher" -> "Docente"
    "coordinator" -> "Coordinatore"
    "parent" -> "Genitore"
    "secretary" -> "Segreteria"
    "admin", "superadmin" -> "Admin"
    else -> role.replaceFirstChar { it.uppercase() }
}

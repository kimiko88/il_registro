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

data class DocumentSignatureItem(
    val id: String,
    val title: String,
    val deadline: String,
    val signatoriesProgress: String,
    var isSigned: Boolean = false
)

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun TeacherDigitalSignatureScreen(
    documents: List<DocumentSignatureItem> = emptyList(),
    onBack: () -> Unit = {}
) {
    var showOtpDialog by remember { mutableStateOf(false) }
    var otpCode by remember { mutableStateOf("") }
    var selectedDocId by remember { mutableStateOf("") }
    val signedSet = remember { mutableStateListOf<String>() }

    Scaffold(
        topBar = {
            TopAppBar(
                title = { Text(stringResource(R.string.sign_hour), fontWeight = FontWeight.Bold) },
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
                        Text(stringResource(R.string.scrutiny_board), fontWeight = FontWeight.Bold, fontSize = 16.sp)
                        Text(stringResource(R.string.select_class), fontSize = 12.sp, color = Color.DarkGray)
                    }
                }
            }

            if (documents.isNotEmpty()) {
                items(documents) { doc ->
                    val isSigned = doc.isSigned || signedSet.contains(doc.id)
                    Card(modifier = Modifier.fillMaxWidth(), shape = RoundedCornerShape(12.dp)) {
                        Column(modifier = Modifier.padding(14.dp)) {
                            Row(
                                modifier = Modifier.fillMaxWidth(),
                                horizontalArrangement = Arrangement.SpaceBetween,
                                verticalAlignment = Alignment.CenterVertically
                            ) {
                                Text(doc.title, fontWeight = FontWeight.Bold)
                                Surface(
                                    color = if (isSigned) Color(0xFF10B981) else Color(0xFFF59E0B),
                                    shape = RoundedCornerShape(6.dp)
                                ) {
                                    Text(
                                        if (isSigned) "Firmato (PAdES)" else "In Attesa",
                                        color = Color.White,
                                        modifier = Modifier.padding(horizontal = 6.dp, vertical = 2.dp),
                                        fontSize = 10.sp
                                    )
                                }
                            }
                            Text("${doc.deadline} • ${doc.signatoriesProgress}", fontSize = 12.sp, color = Color.Gray, modifier = Modifier.padding(vertical = 4.dp))

                            if (!isSigned) {
                                Button(
                                    onClick = {
                                        selectedDocId = doc.id
                                        showOtpDialog = true
                                    },
                                    shape = RoundedCornerShape(8.dp),
                                    modifier = Modifier.fillMaxWidth().padding(top = 6.dp)
                                ) {
                                    Icon(Icons.Default.Fingerprint, contentDescription = null, modifier = Modifier.size(16.dp))
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

    if (showOtpDialog) {
        AlertDialog(
            onDismissRequest = { showOtpDialog = false },
            title = { Text(stringResource(R.string.sign_hour)) },
            text = {
                Column(verticalArrangement = Arrangement.spacedBy(8.dp)) {
                    Text("OTP:")
                    OutlinedTextField(
                        value = otpCode,
                        onValueChange = { if (it.length <= 6) otpCode = it },
                        label = { Text("OTP (6 cifre)") },
                        modifier = Modifier.fillMaxWidth()
                    )
                }
            },
            confirmButton = {
                Button(
                    onClick = {
                        if (otpCode.length == 6) {
                            if (selectedDocId.isNotEmpty()) {
                                signedSet.add(selectedDocId)
                            }
                            showOtpDialog = false
                        }
                    },
                    enabled = otpCode.length == 6
                ) {
                    Text("OK")
                }
            },
            dismissButton = {
                TextButton(onClick = { showOtpDialog = false }) {
                    Text("Annulla")
                }
            }
        )
    }
}

package it.scuola.registro.teacher.ui

import androidx.compose.foundation.layout.*
import androidx.compose.foundation.lazy.LazyColumn
import androidx.compose.foundation.shape.RoundedCornerShape
import androidx.compose.material.icons.Icons
import androidx.compose.material.icons.filled.*
import androidx.compose.material3.*
import androidx.compose.runtime.*
import androidx.compose.ui.Alignment
import androidx.compose.ui.Modifier
import androidx.compose.ui.graphics.Color
import androidx.compose.ui.text.font.FontWeight
import androidx.compose.ui.unit.dp
import androidx.compose.ui.unit.sp

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun TeacherDigitalSignatureScreen(onBack: () -> Unit = {}) {
    var showOtpDialog by remember { mutableStateOf(false) }
    var otpCode by remember { mutableStateOf("") }
    var signatureApplied by remember { mutableStateOf(false) }

    Scaffold(
        topBar = {
            TopAppBar(
                title = { Text("Firma Digitale Verbali (FEA)", fontWeight = FontWeight.Bold) },
                navigationIcon = {
                    IconButton(onClick = onBack) {
                        Icon(Icons.Default.ArrowBack, contentDescription = "Indietro")
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
                        Text("Firma Elettronica Avanzata (FEA - CAD)", fontWeight = FontWeight.Bold, fontSize = 16.sp)
                        Text("Apposizione firma con valore legale sui verbali di scrutinio della classe tramite credenziali protette da OTP o Biometria.", fontSize = 12.sp, color = Color.DarkGray)
                    }
                }
            }

            item {
                Text("Documenti in Attesa di Firma", fontWeight = FontWeight.Bold, fontSize = 16.sp)
            }

            item {
                Card(modifier = Modifier.fillMaxWidth(), shape = RoundedCornerShape(12.dp)) {
                    Column(modifier = Modifier.padding(14.dp)) {
                        Row(
                            modifier = Modifier.fillMaxWidth(),
                            horizontalArrangement = Arrangement.SpaceBetween,
                            verticalAlignment = Alignment.CenterVertically
                        ) {
                            Text("Verbale Scrutinio Finale - Classe 3ª A", fontWeight = FontWeight.Bold)
                            Surface(
                                color = if (signatureApplied) Color(0xFF10B981) else Color(0xFFF59E0B),
                                shape = RoundedCornerShape(6.dp)
                            ) {
                                Text(
                                    if (signatureApplied) "Firmato (PAdES)" else "In Attesa di Firma",
                                    color = Color.White,
                                    modifier = Modifier.padding(horizontal = 6.dp, vertical = 2.dp),
                                    fontSize = 10.sp
                                )
                            }
                        }
                        Text("Scadenza: 15 Giugno 2026 • 12 Docenti Firmatari (11/12 firmati)", fontSize = 12.sp, color = Color.Gray, modifier = Modifier.padding(vertical = 4.dp))

                        if (!signatureApplied) {
                            Button(
                                onClick = { showOtpDialog = true },
                                shape = RoundedCornerShape(8.dp),
                                modifier = Modifier.fillMaxWidth().padding(top = 6.dp)
                            ) {
                                Icon(Icons.Default.Fingerprint, contentDescription = null, modifier = Modifier.size(16.dp))
                                Spacer(modifier = Modifier.width(6.dp))
                                Text("Firma Documento con FEA / OTP")
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
            title = { Text("Conferma Firma FEA") },
            text = {
                Column(verticalArrangement = Arrangement.spacedBy(8.dp)) {
                    Text("Inserisci il codice OTP generato dall'app o ricevuto via SMS/Push:")
                    OutlinedTextField(
                        value = otpCode,
                        onValueChange = { if (it.length <= 6) otpCode = it },
                        label = { Text("Codice OTP (6 cifre)") },
                        modifier = Modifier.fillMaxWidth()
                    )
                }
            },
            confirmButton = {
                Button(
                    onClick = {
                        if (otpCode.length == 6) {
                            signatureApplied = true
                            showOtpDialog = false
                        }
                    },
                    enabled = otpCode.length == 6
                ) {
                    Text("Firma e Sigilla")
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

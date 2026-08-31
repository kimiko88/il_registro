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

data class ParentSignatureDocumentItem(
    val id: String,
    val title: String,
    val description: String,
    val isSigned: Boolean = false
)

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun ParentSignaturesScreen(
    documents: List<ParentSignatureDocumentItem> = emptyList(),
    onSignDocument: (String) -> Unit = {},
    onBack: () -> Unit = {}
) {
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
                        Text(stringResource(R.string.child_grades), fontSize = 12.sp, color = Color.DarkGray)
                    }
                }
            }

            if (documents.isNotEmpty()) {
                items(documents) { doc ->
                    Card(modifier = Modifier.fillMaxWidth(), shape = RoundedCornerShape(12.dp)) {
                        Column(modifier = Modifier.padding(14.dp)) {
                            Row(
                                modifier = Modifier.fillMaxWidth(),
                                horizontalArrangement = Arrangement.SpaceBetween,
                                verticalAlignment = Alignment.CenterVertically
                            ) {
                                Text(doc.title, fontWeight = FontWeight.Bold)
                                Surface(
                                    color = if (doc.isSigned) Color(0xFF10B981) else Color(0xFFEF4444),
                                    shape = RoundedCornerShape(6.dp)
                                ) {
                                    Text(
                                        if (doc.isSigned) "Firmato" else "Richiesto",
                                        color = Color.White,
                                        modifier = Modifier.padding(horizontal = 6.dp, vertical = 2.dp),
                                        fontSize = 10.sp
                                    )
                                }
                            }
                            Text(doc.description, fontSize = 12.sp, color = Color.Gray, modifier = Modifier.padding(vertical = 4.dp))

                            if (!doc.isSigned) {
                                Button(
                                    onClick = { onSignDocument(doc.id) },
                                    shape = RoundedCornerShape(8.dp),
                                    modifier = Modifier.fillMaxWidth().padding(top = 4.dp)
                                ) {
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

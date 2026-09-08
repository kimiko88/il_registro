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

data class PagoPaNoticeItem(
    val id: String,
    val title: String,
    val deadlineAndAmount: String,
    val isPaid: Boolean = false
)

@OptIn(ExperimentalMaterial3Api::class)
@Composable
fun ParentPaymentsScreen(
    notices: List<PagoPaNoticeItem> = emptyList(),
    onPay: (String) -> Unit = {},
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
            if (notices.isNotEmpty()) {
                items(notices) { notice ->
                    Card(modifier = Modifier.fillMaxWidth(), shape = RoundedCornerShape(12.dp)) {
                        Row(
                            modifier = Modifier.fillMaxWidth().padding(14.dp),
                            horizontalArrangement = Arrangement.SpaceBetween,
                            verticalAlignment = Alignment.CenterVertically
                        ) {
                            Column {
                                Text(notice.title, fontWeight = FontWeight.Bold)
                                Text(notice.deadlineAndAmount, fontSize = 12.sp, color = Color.Gray)
                            }
                            if (notice.isPaid) {
                                Surface(color = Color(0xFF10B981), shape = RoundedCornerShape(6.dp)) {
                                    Text("Pagato", color = Color.White, modifier = Modifier.padding(horizontal = 8.dp, vertical = 4.dp), fontSize = 11.sp)
                                }
                            } else {
                                Button(
                                    onClick = { onPay(notice.id) },
                                    colors = ButtonDefaults.buttonColors(containerColor = Color(0xFF0066CC)),
                                    shape = RoundedCornerShape(8.dp)
                                ) {
                                    Text("PagoPA", fontSize = 11.sp)
                                }
                            }
                        }
                    }
                }
            }
        }
    }
}

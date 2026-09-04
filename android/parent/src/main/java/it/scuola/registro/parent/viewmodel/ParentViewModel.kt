package it.scuola.registro.parent.viewmodel

import androidx.compose.runtime.getValue
import androidx.compose.runtime.mutableStateListOf
import androidx.compose.runtime.mutableStateOf
import androidx.compose.runtime.setValue
import it.scuola.registro.parent.data.ColloquioBooking
import it.scuola.registro.parent.data.ParentChild
import it.scuola.registro.parent.data.PendingAbsence
import it.scuola.registro.parent.network.HttpParentApiService
import it.scuola.registro.parent.network.ParentApiService
import it.scuola.registro.parent.ui.ParentGradeDisplayItem

class ParentViewModel(
    private val apiService: ParentApiService = HttpParentApiService()
) {
    var children by mutableStateOf<List<ParentChild>>(emptyList())
        private set

    var selectedChildId by mutableStateOf("")

    var pendingAbsences = mutableStateListOf<PendingAbsence>()
        private set

    var availableColloqui = mutableStateListOf<ColloquioBooking>()
        private set

    var gradesList by mutableStateOf<List<ParentGradeDisplayItem>>(emptyList())
        private set

    var circularsList by mutableStateOf<List<String>>(emptyList())
        private set

    var isLoading by mutableStateOf(false)
        private set

    var errorMessage by mutableStateOf<String?>(null)
        private set

    init {
        loadSampleData()
    }

    suspend fun loadFromDatabase(token: String): Boolean {
        isLoading = true
        errorMessage = null
        try {
            val childrenResult = apiService.getChildren(token)
            if (childrenResult.isSuccess) {
                val apiChildren = childrenResult.getOrNull()
                if (!apiChildren.isNullOrEmpty()) {
                    children = apiChildren
                    selectedChildId = apiChildren.first().id
                    val absencesResult = apiService.getAbsences(token, selectedChildId)
                    if (absencesResult.isSuccess) {
                        val apiAbsences = absencesResult.getOrNull()
                        if (!apiAbsences.isNullOrEmpty()) {
                            pendingAbsences.clear()
                            pendingAbsences.addAll(apiAbsences)
                        }
                    }
                }
            }
            isLoading = false
            return true
        } catch (e: Exception) {
            errorMessage = e.localizedMessage
            isLoading = false
            return false
        }
    }

    fun loadSampleData() {
        children = listOf(
            ParentChild("c1", "Marco", "Rossi", "Classe 2A"),
            ParentChild("c2", "Giulia", "Rossi", "Classe 4B")
        )
        selectedChildId = "c1"

        pendingAbsences.clear()
        pendingAbsences.addAll(listOf(
            PendingAbsence("a1", "c1", "2026-08-26", "Assenza", "Motivi di salute", false),
            PendingAbsence("a2", "c1", "2026-08-18", "Ritardo", "Ingresso 2a ora", false),
            PendingAbsence("a3", "c2", "2026-08-20", "Assenza", "Visita medica", false)
        ))

        availableColloqui.clear()
        availableColloqui.addAll(listOf(
            ColloquioBooking("col1", "Prof. Bianchi", "Matematica", "2026-09-03 15:30", false),
            ColloquioBooking("col2", "Prof.ssa Rossi", "Italiano", "2026-09-04 16:00", true)
        ))

        gradesList = listOf(
            ParentGradeDisplayItem("Matematica", "8½ (Scritto - 28/08)"),
            ParentGradeDisplayItem("Italiano", "7 (Orale - 25/08)"),
            ParentGradeDisplayItem("Inglese", "8 (Pratico - 22/08)"),
            ParentGradeDisplayItem("Scienze", "7½ (Scritto - 19/08)")
        )

        circularsList = listOf(
            "Circolare n. 12: Inizio anno scolastico e orario provvisorio",
            "Circolare n. 11: Modalità di prenotazione colloqui quadrimestrali",
            "Circolare n. 10: Uscite didattiche e autorizzazioni genitori"
        )
    }

    fun selectChild(childId: String) {
        selectedChildId = childId
    }

    fun getAbsencesForSelectedChild(): List<PendingAbsence> {
        return pendingAbsences.filter { it.childId == selectedChildId }
    }

    fun justifyAbsence(absenceId: String, note: String): Boolean {
        val index = pendingAbsences.indexOfFirst { it.id == absenceId }
        if (index != -1) {
            val item = pendingAbsences[index]
            pendingAbsences[index] = item.copy(isJustified = true, justificationNote = note)
            return true
        }
        return false
    }

    fun bookColloquio(colloquioId: String): Boolean {
        val index = availableColloqui.indexOfFirst { it.id == colloquioId }
        if (index != -1) {
            val item = availableColloqui[index]
            availableColloqui[index] = item.copy(isConfirmed = true)
            return true
        }
        return false
    }
}

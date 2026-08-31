package it.scuola.registro.parent.viewmodel

import it.scuola.registro.parent.data.ColloquioBooking
import it.scuola.registro.parent.data.ParentChild
import it.scuola.registro.parent.data.PendingAbsence

class ParentViewModel {
    var children = mutableListOf<ParentChild>()
        private set

    var selectedChildId: String = ""

    var pendingAbsences = mutableListOf<PendingAbsence>()
        private set

    var availableColloqui = mutableListOf<ColloquioBooking>()
        private set

    fun loadSampleData() {
        children = mutableListOf(
            ParentChild("c1", "Marco", "Rossi", "Classe 2A"),
            ParentChild("c2", "Giulia", "Rossi", "Classe 4B")
        )
        selectedChildId = "c1"

        pendingAbsences = mutableListOf(
            PendingAbsence("a1", "c1", "2026-08-26", "Assenza", "Motivi di salute", false),
            PendingAbsence("a2", "c1", "2026-08-18", "Ritardo", "Ingresso 2a ora", false),
            PendingAbsence("a3", "c2", "2026-08-20", "Assenza", "Visita medica", false)
        )

        availableColloqui = mutableListOf(
            ColloquioBooking("col1", "Prof. Bianchi", "Matematica", "2026-09-03 15:30", false),
            ColloquioBooking("col2", "Prof.ssa Rossi", "Italiano", "2026-09-04 16:00", true)
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

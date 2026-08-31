import SwiftUI

struct StudentEPortfolioView: View {
    @State private var showingAddSheet = false
    @State private var title = ""
    @State private var description = ""
    @State private var reflection = ""

    var body: some View {
        NavigationView {
            List {
                Section(header: Text("Linee Guida Nazionali per l'Orientamento")) {
                    VStack(alignment: .leading, spacing: 6) {
                        Text("Tracciamento Ore di Orientamento")
                            .font(.headline)
                        ProgressView(value: 30, total: 30)
                            .tint(.teal)
                        Text("30 / 30 ore annuali svolte (Target raggiunto)")
                            .font(.caption)
                            .foregroundColor(.secondary)
                    }
                    .padding(.vertical, 4)
                }

                Section(header: Text("I Tuoi Capolavori (E-Portfolio)")) {
                    VStack(alignment: .leading, spacing: 6) {
                        HStack {
                            Text("Registro Elettronico Open Source")
                                .fontWeight(.bold)
                            Spacer()
                            Text("A.S. 2025/2026")
                                .font(.caption2)
                                .padding(.horizontal, 6)
                                .padding(.vertical, 2)
                                .background(Color.purple.opacity(0.15))
                                .foregroundColor(.purple)
                                .cornerRadius(4)
                        }
                        Text("Sviluppo di architettura concorrente in Go e client mobile.")
                            .font(.subheadline)
                            .foregroundColor(.secondary)
                        Text("Autovalutazione: Incrementate competenze di coding e lavoro cooperativo.")
                            .font(.caption)
                            .foregroundColor(.primary)
                            .padding(6)
                            .background(Color.gray.opacity(0.1))
                            .cornerRadius(6)
                    }
                }
            }
            .navigationTitle("E-Portfolio & Capolavori")
            .toolbar {
                Button(action: { showingAddSheet = true }) {
                    Image(systemName: "plus")
                }
            }
            .sheet(isPresented: $showingAddSheet) {
                NavigationView {
                    Form {
                        Section(header: Text("Dettagli Capolavoro")) {
                            TextField("Titolo del Capolavoro", text: $title)
                            TextField("Descrizione", text: $description)
                            TextField("Riflessione critica", text: $reflection)
                        }
                    }
                    .navigationTitle("Nuovo Capolavoro")
                    .toolbar {
                        ToolbarItem(placement: .cancellationAction) {
                            Button("Annulla") { showingAddSheet = false }
                        }
                        ToolbarItem(placement: .confirmationAction) {
                            Button("Salva") { showingAddSheet = false }
                        }
                    }
                }
            }
        }
    }
}

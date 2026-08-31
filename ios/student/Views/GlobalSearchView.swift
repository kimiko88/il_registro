import SwiftUI

struct GlobalSearchView: View {
    @State private var searchText = ""

    var body: some View {
        NavigationView {
            List {
                Section(header: Text("Risultati Suggeriti")) {
                    VStack(alignment: .leading, spacing: 4) {
                        Text("Circolare N. 214: Calendario Scrutini")
                            .fontWeight(.bold)
                        Text("Circolari • 26 Maggio 2026")
                            .font(.caption)
                            .foregroundColor(.blue)
                    }
                    .padding(.vertical, 4)

                    VStack(alignment: .leading, spacing: 4) {
                        Text("Formulario Derivate ed Integrali.pdf")
                            .fontWeight(.bold)
                        Text("Materiale Didattico • Matematica")
                            .font(.caption)
                            .foregroundColor(.blue)
                    }
                    .padding(.vertical, 4)
                }
            }
            .searchable(text: $searchText, prompt: "Cerca nel registro...")
            .navigationTitle("Ricerca Globale")
        }
    }
}

import SwiftUI

struct TeacherGlobalReportsView: View {
    var body: some View {
        NavigationView {
            List {
                Section(header: Text("Profilo & Consiglio Orientativo - Mario Rossi")) {
                    VStack(alignment: .leading, spacing: 6) {
                        HStack {
                            Text("Giudizio Globale Finale")
                                .fontWeight(.bold)
                            Spacer()
                            Text("Deliberato")
                                .font(.caption2)
                                .padding(.horizontal, 6)
                                .padding(.vertical, 2)
                                .background(Color.green.opacity(0.2))
                                .foregroundColor(.green)
                                .cornerRadius(4)
                        }
                        Text("Impegno costante e maturo, competenze eccellenti.")
                            .font(.subheadline)
                            .foregroundColor(.blue)
                        Text("Consiglio Orientativo: Percorso Scientifico / Ingegneria")
                            .font(.caption)
                            .foregroundColor(.secondary)
                    }
                    .padding(.vertical, 4)
                }
            }
            .navigationTitle("Giudizi Globali")
        }
    }
}

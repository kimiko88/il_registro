import SwiftUI

struct TeacherSupportView: View {
    var body: some View {
        NavigationView {
            List {
                Section(header: Text("Studenti con Sostegno & PEI")) {
                    VStack(alignment: .leading, spacing: 6) {
                        HStack {
                            Text("Studente L. M. • 3ª A")
                                .fontWeight(.bold)
                            Spacer()
                            Text("18 h/sett")
                                .font(.caption2)
                                .padding(.horizontal, 6)
                                .padding(.vertical, 2)
                                .background(Color.blue.opacity(0.2))
                                .foregroundColor(.blue)
                                .cornerRadius(4)
                        }
                        Text("PEI Approvato • Programmazione Equipollente")
                            .font(.subheadline)
                            .foregroundColor(.blue)
                        Text("Diario: Aggiornato al 28 Maggio 2026")
                            .font(.caption)
                            .foregroundColor(.secondary)
                    }
                    .padding(.vertical, 4)
                }
            }
            .navigationTitle("Registro Sostegno")
        }
    }
}

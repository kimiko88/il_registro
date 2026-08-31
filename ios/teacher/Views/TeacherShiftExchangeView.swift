import SwiftUI

struct TeacherShiftExchangeView: View {
    var body: some View {
        NavigationView {
            List {
                Section(header: Text("Scambi Orario Approvati")) {
                    VStack(alignment: .leading, spacing: 6) {
                        HStack {
                            Text("Prof. Rossi ↔ Prof.ssa Neri")
                                .fontWeight(.bold)
                            Spacer()
                            Text("Approvato")
                                .font(.caption2)
                                .padding(.horizontal, 6)
                                .padding(.vertical, 2)
                                .background(Color.green.opacity(0.2))
                                .foregroundColor(.green)
                                .cornerRadius(4)
                        }
                        Text("Matematica • Giovedì 4 Giu ↔ Venerdì 5 Giu")
                            .font(.subheadline)
                            .foregroundColor(.blue)
                        Text("2ª Ora (3ª A) con 4ª Ora (4ª B)")
                            .font(.caption)
                            .foregroundColor(.secondary)
                    }
                    .padding(.vertical, 4)
                }
            }
            .navigationTitle("Scambi Turno")
        }
    }
}

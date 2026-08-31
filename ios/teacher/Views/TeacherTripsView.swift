import SwiftUI

struct TeacherTripsView: View {
    var body: some View {
        NavigationView {
            List {
                Section(header: Text("Viaggi d'Istruzione Coordinati")) {
                    VStack(alignment: .leading, spacing: 6) {
                        HStack {
                            Text("Laboratori Gran Sasso")
                                .fontWeight(.bold)
                            Spacer()
                            Text("Approvata")
                                .font(.caption2)
                                .padding(.horizontal, 6)
                                .padding(.vertical, 2)
                                .background(Color.green.opacity(0.2))
                                .foregroundColor(.green)
                                .cornerRadius(4)
                        }
                        Text("Classi 4ª A, 4ª B • 12 Giugno 2026")
                            .font(.subheadline)
                            .foregroundColor(.blue)
                        Text("Accompagnatori: Prof. Rossi, Prof.ssa Bianchi • Aut: 42/45")
                            .font(.caption)
                            .foregroundColor(.secondary)
                    }
                    .padding(.vertical, 4)
                }
            }
            .navigationTitle("Viaggi & Uscite")
        }
    }
}

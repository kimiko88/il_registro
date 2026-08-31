import SwiftUI

struct TeacherDailySubstitutionsView: View {
    var body: some View {
        NavigationView {
            List {
                Section(header: Text("Coperture Assegnate Oggi")) {
                    VStack(alignment: .leading, spacing: 6) {
                        HStack {
                            Text("3ª Ora (10:00 - 11:00)")
                                .fontWeight(.bold)
                            Spacer()
                            Text("Classe 2ª C")
                                .font(.caption2)
                                .padding(.horizontal, 6)
                                .padding(.vertical, 2)
                                .background(Color.blue.opacity(0.2))
                                .foregroundColor(.blue)
                                .cornerRadius(4)
                        }
                        Text("Docente: Prof. Gialli • Lab Lingue 1")
                            .font(.subheadline)
                            .foregroundColor(.blue)
                        Text("Attività: Sorveglianza e compiti Classroom")
                            .font(.caption)
                            .foregroundColor(.secondary)
                    }
                    .padding(.vertical, 4)
                }
            }
            .navigationTitle("Supplenze del Giorno")
        }
    }
}

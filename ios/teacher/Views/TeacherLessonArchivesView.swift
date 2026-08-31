import SwiftUI

struct TeacherLessonArchivesView: View {
    var body: some View {
        NavigationView {
            List {
                Section(header: Text("Lezioni Svolte - Classe 3ª A")) {
                    VStack(alignment: .leading, spacing: 6) {
                        HStack {
                            Text("Lezione N. 142: De L'Hôpital")
                                .fontWeight(.bold)
                            Spacer()
                            Text("Firmata")
                                .font(.caption2)
                                .padding(.horizontal, 6)
                                .padding(.vertical, 2)
                                .background(Color.green.opacity(0.2))
                                .foregroundColor(.green)
                                .cornerRadius(4)
                        }
                        Text("20 Maggio 2026 • 2ª Ora (09:00 - 10:00)")
                            .font(.subheadline)
                            .foregroundColor(.blue)
                        Text("Compiti: Esercizi pag. 320 nn. 45-52")
                            .font(.caption)
                            .foregroundColor(.secondary)
                    }
                    .padding(.vertical, 4)
                }
            }
            .navigationTitle("Archivio Lezioni")
        }
    }
}

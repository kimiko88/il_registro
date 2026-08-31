import SwiftUI

struct TeacherRecoveryCoursesView: View {
    var body: some View {
        NavigationView {
            List {
                Section(header: Text("Corsi di Recupero e PAI Attivi")) {
                    VStack(alignment: .leading, spacing: 6) {
                        HStack {
                            Text("Recupero Matematica 3ª")
                                .fontWeight(.bold)
                            Spacer()
                            Text("In Corso (8/10 h)")
                                .font(.caption2)
                                .padding(.horizontal, 6)
                                .padding(.vertical, 2)
                                .background(Color.green.opacity(0.2))
                                .foregroundColor(.green)
                                .cornerRadius(4)
                        }
                        Text("8 studenti con sospensione del giudizio")
                            .font(.subheadline)
                            .foregroundColor(.secondary)
                        Text("Verifica finale: 28 Agosto 2026")
                            .font(.caption)
                            .foregroundColor(.blue)
                    }
                    .padding(.vertical, 4)
                }
            }
            .navigationTitle("Corsi di Recupero")
        }
    }
}

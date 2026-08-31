import SwiftUI

struct TeacherTestBookingView: View {
    var body: some View {
        NavigationView {
            List {
                Section(header: Text("Verifiche Programmate - Classe 3ª A")) {
                    VStack(alignment: .leading, spacing: 6) {
                        HStack {
                            Text("Compito di Matematica")
                                .fontWeight(.bold)
                            Spacer()
                            Text("Nessun Conflitto")
                                .font(.caption2)
                                .padding(.horizontal, 6)
                                .padding(.vertical, 2)
                                .background(Color.green.opacity(0.2))
                                .foregroundColor(.green)
                                .cornerRadius(4)
                        }
                        Text("Giovedì 4 Giugno • 09:00 - 11:00 (2h)")
                            .font(.subheadline)
                            .foregroundColor(.blue)
                        Text("Funzioni esponenziali e logaritmiche")
                            .font(.caption)
                            .foregroundColor(.secondary)
                    }
                    .padding(.vertical, 4)
                }
            }
            .navigationTitle("Prenotazione Verifiche")
        }
    }
}

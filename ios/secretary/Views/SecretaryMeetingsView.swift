import SwiftUI

struct SecretaryMeetingsView: View {
    var body: some View {
        NavigationView {
            List {
                Section(header: Text("Riunioni & Consigli Convocati")) {
                    VStack(alignment: .leading, spacing: 6) {
                        HStack {
                            Text("Consiglio di Classe Straordinario - 3ª A")
                                .fontWeight(.bold)
                            Spacer()
                            Text("Convocato")
                                .font(.caption2)
                                .padding(.horizontal, 6)
                                .padding(.vertical, 2)
                                .background(Color.blue.opacity(0.2))
                                .foregroundColor(.blue)
                                .cornerRadius(4)
                        }
                        Text("28 Maggio • Ore 17:00 • Aula Magna")
                            .font(.subheadline)
                            .foregroundColor(.blue)
                        Text("O.d.G.: Andamento didattico e PDP")
                            .font(.caption)
                            .foregroundColor(.secondary)
                    }
                    .padding(.vertical, 4)
                }
            }
            .navigationTitle("Riunioni & Consigli")
        }
    }
}

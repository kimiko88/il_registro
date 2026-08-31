import SwiftUI

struct TeacherSeatingPlanView: View {
    var body: some View {
        NavigationView {
            List {
                Section(header: Text("Layout Aula 12 - Classe 3ª A")) {
                    VStack(alignment: .leading, spacing: 6) {
                        HStack {
                            Text("Isole di Lavoro (Cooperative)")
                                .fontWeight(.bold)
                            Spacer()
                            Text("24 Postazioni")
                                .font(.caption2)
                                .padding(.horizontal, 6)
                                .padding(.vertical, 2)
                                .background(Color.green.opacity(0.2))
                                .foregroundColor(.green)
                                .cornerRadius(4)
                        }
                        Text("Isola 1: Mario Rossi, Luca Bianchi, Giulia Verdi, Sara Neri")
                            .font(.subheadline)
                            .foregroundColor(.blue)
                        Text("Visibilità ottimale LIM e sintesi vocale PDP")
                            .font(.caption)
                            .foregroundColor(.secondary)
                    }
                    .padding(.vertical, 4)
                }
            }
            .navigationTitle("Disposizione Banchi")
        }
    }
}

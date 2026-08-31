import SwiftUI

struct TeacherColloquiView: View {
    var body: some View {
        NavigationView {
            List {
                Section(header: Text("Genitori Prenotati")) {
                    VStack(alignment: .leading, spacing: 6) {
                        HStack {
                            Text("Sig. Rossi (Studente: Mario)")
                                .fontWeight(.bold)
                            Spacer()
                            Text("10:00 - 10:15")
                                .font(.caption2)
                                .padding(.horizontal, 6)
                                .padding(.vertical, 2)
                                .background(Color.blue.opacity(0.2))
                                .foregroundColor(.blue)
                                .cornerRadius(4)
                        }
                        Text("In Presenza • Aula 12")
                            .font(.subheadline)
                            .foregroundColor(.blue)
                        Text("Motivo: Andamento didattico e recupero")
                            .font(.caption)
                            .foregroundColor(.secondary)
                    }
                    .padding(.vertical, 4)
                }
            }
            .navigationTitle("Ricevimento Famiglie")
        }
    }
}

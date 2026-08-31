import SwiftUI

struct ParentOnlineEnrollmentView: View {
    var body: some View {
        NavigationView {
            List {
                Section(header: Text("Domanda Iscrizione Classi Prime")) {
                    VStack(alignment: .leading, spacing: 6) {
                        HStack {
                            Text("Domanda ISC-2026-4019")
                                .fontWeight(.bold)
                            Spacer()
                            Text("Accettata")
                                .font(.caption2)
                                .padding(.horizontal, 6)
                                .padding(.vertical, 2)
                                .background(Color.green.opacity(0.2))
                                .foregroundColor(.green)
                                .cornerRadius(4)
                        }
                        Text("Liceo Scientifico Scienze Applicate")
                            .font(.subheadline)
                            .foregroundColor(.blue)
                        Text("Studente: Giovanni Rossi • Sede Centrale • Lingua: Spagnolo")
                            .font(.caption)
                            .foregroundColor(.secondary)
                    }
                    .padding(.vertical, 4)
                }
            }
            .navigationTitle("Iscrizioni Online")
        }
    }
}

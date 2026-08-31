import SwiftUI

struct TeacherExamKioskView: View {
    var body: some View {
        NavigationView {
            List {
                Section(header: Text("Sessione Kiosk Aula 12 - Classe 3ª A")) {
                    VStack(alignment: .leading, spacing: 6) {
                        HStack {
                            Text("Verifica di Informatica")
                                .fontWeight(.bold)
                            Spacer()
                            Text("24 Bloccati")
                                .font(.caption2)
                                .padding(.horizontal, 6)
                                .padding(.vertical, 2)
                                .background(Color.green.opacity(0.2))
                                .foregroundColor(.green)
                                .cornerRadius(4)
                        }
                        Text("Autonomous Single App Mode (ASAM)")
                            .font(.subheadline)
                            .foregroundColor(.blue)
                        Text("Tentativi di uscita: 0 • Batteria media: 94%")
                            .font(.caption)
                            .foregroundColor(.secondary)
                    }
                    .padding(.vertical, 4)
                }
            }
            .navigationTitle("Modalità Esame Kiosk")
        }
    }
}

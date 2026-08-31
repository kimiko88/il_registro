import SwiftUI

struct StudentTripsView: View {
    var body: some View {
        NavigationView {
            List {
                Section(header: Text("Viaggi d'Istruzione e Visite")) {
                    VStack(alignment: .leading, spacing: 6) {
                        HStack {
                            Text("Viaggio d'Istruzione a Vienna & Praga")
                                .fontWeight(.bold)
                            Spacer()
                            Text("Approvato")
                                .font(.caption2)
                                .padding(.horizontal, 6)
                                .padding(.vertical, 2)
                                .background(Color.green.opacity(0.2))
                                .foregroundColor(.green)
                                .cornerRadius(4)
                        }
                        Text("12 - 17 Aprile 2026 • 5 Giorni")
                            .font(.subheadline)
                            .foregroundColor(.blue)
                        Text("Accompagnatori: Prof. Rossi, Prof.ssa Verdi")
                            .font(.caption)
                            .foregroundColor(.secondary)
                        Text("Quota PagoPA: Saldata")
                            .font(.caption2)
                            .foregroundColor(.green)
                    }
                    .padding(.vertical, 4)
                }
            }
            .navigationTitle("Viaggi d'Istruzione")
        }
    }
}
